package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/url"
	"os"
	"path"
	"regexp"
	"strings"
	"sync"
	"time"
)

var (
	VERSION = "0.1.2" // declare-version
	zoneRP  = regexp.MustCompile("//(z[0-9]){0, 1}/")
)

type FetchMessage struct {
	path    string
	ok      bool
	message string
}

type WorkerMetrics struct {
	jobs           int64
	doneJobs       int64
	successfulJobs int64
	failedJobs     int64
}

type CMD struct {
	input           string
	prefix          string
	isJSONList      bool
	pool            uint
	scannerFinished bool
	readWg          sync.WaitGroup
	readCh          chan string
	messageWg       sync.WaitGroup
	messageCh       chan FetchMessage
	metrics         *WorkerMetrics
	metricsMu       sync.Mutex
}

func (c *CMD) init() {
	flag.StringVar(&c.input, "input", "", "input file list")
	flag.StringVar(&c.prefix, "prefix", "/", "local path prefix, default is root")
	flag.UintVar(&c.pool, "pool", 20, "concurrent fetch pool size, default 20")
	flag.BoolVar(&c.isJSONList, "is-jsonlist", false, "each line should be a json string if this is set true, default false")

	vFlag := flag.Bool("V", false, "version")
	hFlag := flag.Bool("h", false, "help")

	flag.Parse()
	if *vFlag {
		fmt.Println(VERSION)
		os.Exit(0)
	}
	if *hFlag {
		flag.Usage()
		os.Exit(0)
	}
	if c.input == "" {
		flag.Usage()
		os.Exit(1)
	}
	c.metrics = &WorkerMetrics{
		jobs:           0,
		doneJobs:       0,
		successfulJobs: 0,
		failedJobs:     0,
	}
}

func (c *CMD) start() {
	c.readWg = sync.WaitGroup{}
	c.messageWg = sync.WaitGroup{}
	c.readCh = make(chan string, c.pool)
	c.messageCh = make(chan FetchMessage, c.pool)
	c.metricsMu = sync.Mutex{}
	c.scannerFinished = false
	file, err := os.Open(c.input)
	if err != nil {
		log.Fatalf("open failed: %v", err)
	}
	defer file.Close()

	go func() {
		for {
			select {
			case msg := <-c.messageCh:
				c.messageWorker(msg)
				c.messageWg.Done()
			}
		}
	}()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		c.scannerAdjust()
		c.readCh <- scanner.Text()
		c.metricsMu.Lock()
		c.metrics.jobs++
		c.metricsMu.Unlock()
		c.readWg.Add(1)
		go c.worker()
	}

	if err := scanner.Err(); err != nil {
		fmt.Printf("[%s] scan error: %v\n", time.Now().Format("2006-01-02 15:04:05"), err)
	} else {
		func() {
			c.messageWg.Add(1)
			msg := FetchMessage{
				path:    "",
				ok:      true,
				message: "all done",
			}
			c.messageCh <- msg
		}()
	}

	time.Sleep(time.Second)

	c.readWg.Wait()
	c.messageWg.Wait()
	close(c.readCh)
	close(c.messageCh)
}

func (c *CMD) postMsg(msg FetchMessage) {
	c.messageWg.Add(1)
	c.messageCh <- msg
	c.readWg.Done()
	return
}

func (c *CMD) worker() {
	text := <-c.readCh
	var msg FetchMessage
	var u string

	if text == "" {
		msg.ok = true
		msg.message = "empty line"
		msg.path = text
		c.postMsg(msg)
		return
	} else if c.isJSONList {
		var tree map[string]interface{}
		e := json.Unmarshal([]byte(text), &tree)
		if e != nil {
			msg.ok = false
			msg.message = fmt.Sprintf("parse item `%s` to json failed", text)
			msg.path = text
			c.postMsg(msg)
			return
		}

		u, _ = tree["url"].(string)

		if u == "" {
			msg.ok = false
			msg.message = fmt.Sprintf("failed to get url from item `%s`", text)
			msg.path = text
			c.postMsg(msg)
			return
		}
	} else {
		slices := strings.Split(text, " ")
		pathSlice := 0
		for i := 0; i < len(slices); i++ {
			if i == 0 || strings.HasSuffix(slices[i], "\\") {
				pathSlice++
			} else {
				break
			}
		}
		u = strings.Join(slices[:pathSlice], " ")
	}

	p := url2LocalPath(u, c.prefix)
	_, e := ioutil.ReadFile(p)

	if e != nil {
		msg.ok = false
		msg.message = "read file failed"
	} else {
		msg.ok = true
		msg.message = ""
	}

	msg.path = u
	c.postMsg(msg)
}

func (c *CMD) messageWorker(msg FetchMessage) {
	c.metricsMu.Lock()
	if !msg.ok {
		fmt.Printf("[%s] load file from %s failed, error: %s\n", time.Now().Format("2006-01-02 15:04:05"), msg.path, msg.message)
		c.metrics.failedJobs++
	} else {
		c.metrics.successfulJobs++
	}
	c.metrics.doneJobs++

	if c.metrics.doneJobs%100 == 0 {
		fmt.Printf("[%s] loaded %d, failed: %d, success: %d\n", time.Now().Format("2006-01-02 15:04:05"), c.metrics.doneJobs, c.metrics.failedJobs, c.metrics.successfulJobs)
	}

	if msg.ok && msg.message == "all done" {
		c.scannerFinished = true
		fmt.Printf("[%s] finish scanning, totally %d files to load\n", time.Now().Format("2006-01-02 15:04:05"), c.metrics.jobs)
	}

	if c.scannerFinished && c.metrics.doneJobs == c.metrics.jobs+1 {
		fmt.Printf("[%s] handled all items, total: %d, failed: %d, success: %d\n", time.Now().Format("2006-01-02 15:04:05"), c.metrics.jobs, c.metrics.failedJobs, c.metrics.successfulJobs-1)
	}
	c.metricsMu.Unlock()
}

func (c *CMD) scannerAdjust() {
	var pendingJobs int64
	for {
		c.metricsMu.Lock()
		pendingJobs = c.metrics.jobs - c.metrics.doneJobs
		c.metricsMu.Unlock()

		if uint(pendingJobs) > 3*c.pool {
			time.Sleep(100 * time.Microsecond)
		} else {
			break
		}
	}
}

func url2LocalPath(fpath, prefix string) string {
	if strings.HasPrefix(fpath, "http") {
		u, e := url.Parse(fpath)
		if e != nil {
			return ""
		}
		return path.Join(prefix, u.Path)
	} else if strings.HasPrefix(fpath, "qiniu") {
		zoneRP.ReplaceAllString(fpath, "//")
		u, e := url.Parse(fpath)
		if e != nil {
			return ""
		}

		return path.Join(prefix, u.Path)
	}
	return path.Join(prefix, fpath)
}

func main() {
	var cmd = &CMD{}
	cmd.init()
	cmd.start()
}
