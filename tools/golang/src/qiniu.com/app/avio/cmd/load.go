package cmd

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/spf13/cobra"
	util "qiniu.com/app/avio/util"
)

var (
	zoneRP = regexp.MustCompile("//(z[0-9]){0, 1}/")
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
	path            string
	log             string
	scannerFinished bool
	readWg          sync.WaitGroup
	readCh          chan string
	messageWg       sync.WaitGroup
	messageCh       chan FetchMessage
	metrics         *WorkerMetrics
	metricsMu       sync.Mutex
}

func (c *CMD) init(flags *loadFlags) {
	c.input = flags.Source
	c.prefix = ""
	c.pool = flags.PoolSize

	c.isJSONList = false
	if flags.IsJSONList == "true" {
		c.isJSONList = true
	}

	c.path = ""
	if len(flags.Args) == 1 {
		c.path = flags.Args[0]
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

	p := util.Url2LocalPath(u, c.prefix)
	_, e := ioutil.ReadFile(p)

	if e != nil {
		msg.ok = false
		msg.message = fmt.Sprintf("read file failed, read fail message: %v", e)
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

type loadFlags struct {
	Source     string
	IsJSONList string
	LogPath    string
	PoolSize   uint
	Depth      uint
	Args       []string
}

func (f *loadFlags) Validate() error {
	if f.PoolSize > 100 || f.PoolSize < 1 {
		return &util.AvioError{
			Msg: "请设置合法的参数，poolSize 参数最大值为 100，最小值为 1",
		}
	}

	if f.Depth > 10 || f.Depth < 1 {
		return &util.AvioError{
			Msg: "请设置合法的参数，depth 参数最大值为 10，最小值为 1",
		}
	}

	if len(f.Args) > 1 {
		return &util.AvioError{
			Msg: "当前只支持一次加载一个指定目录或者文件，一次加载多个目录或者文件的功能将在后续支持",
		}
	}

	if len(f.Args) == 1 && f.Source != "" {
		return &util.AvioError{
			Msg: "-i/--input-file 参数和 path 只能同时设置一个",
		}
	}

	if f.IsJSONList != "false" && f.IsJSONList != "true" {
		return &util.AvioError{
			Msg: "--is-jsonlist 参数的可选项为 false/true",
		}
	}

	if f.IsJSONList == "true" && f.Source != "" {
		return &util.AvioError{
			Msg: "--is-jsonlist 参数只有在指定了 -i/--input-file 参数时有效",
		}
	}

	if f.Source != "" && f.Depth != 4 {
		fmt.Printf("[WARN] -d/--depth 被设置为 %d，在使用 avio preload -i/--input-file=<list_file> 的模式下将被忽略", f.Depth)
	}
	return nil
}

var mLoadFlags = &loadFlags{}

func init() {
	preloadCmd.Flags().UintVarP(&mLoadFlags.Depth, "depth", "d", 4, "递归加载的最大深度，默认值为 4，最大值为 10")
	preloadCmd.Flags().StringVarP(&mLoadFlags.Source, "input-source", "i", "", "所有需要 preload 的文件的列表文件")
	preloadCmd.Flags().StringVar(&mLoadFlags.IsJSONList, "is-jsonlist", "false", "只在 -i/--input-file 被设置时有效，可选项：true/false，默认值为 false")
	preloadCmd.Flags().UintVarP(&mLoadFlags.PoolSize, "pool", "p", 50, "并行 preload 的并发数，默认值为 50，最大值为 200")
	preloadCmd.Flags().StringVarP(&mLoadFlags.LogPath, "log", "l", "/var/log/avio/preload-<unix_time>.log", "当前任务的日志文件，默认值为 /var/log/avio/preload-<unix_time>.log")
	preloadCmd.SetUsageTemplate(`Usage:{{if .Runnable}}
  {{.UseLine}} [path] {{end}}{{if .HasAvailableSubCommands}}
  {{.CommandPath}} [command]{{end}}{{if gt (len .Aliases) 0}}

Aliases:
  {{.NameAndAliases}}{{end}}{{if .HasExample}}

Examples:
{{.Example}}{{end}}{{if .HasAvailableSubCommands}}

Available Commands:{{range .Commands}}{{if (or .IsAvailableCommand (eq .Name "help"))}}
  {{rpad .Name .NamePadding }} {{.Short}}{{end}}{{end}}{{end}}{{if .HasAvailableLocalFlags}}

Flags:
{{.LocalFlags.FlagUsages | trimTrailingWhitespaces}}{{end}}{{if .HasAvailableInheritedFlags}}

Global Flags:
{{.InheritedFlags.FlagUsages | trimTrailingWhitespaces}}{{end}}{{if .HasHelpSubCommands}}

Additional help topics:{{range .Commands}}{{if .IsAdditionalHelpTopicCommand}}
  {{rpad .CommandPath .CommandPathPadding}} {{.Short}}{{end}}{{end}}{{end}}{{if .HasAvailableSubCommands}}

Use "{{.CommandPath}} [command] --help" for more information about a command.{{end}}
`)
	rootCmd.AddCommand(preloadCmd)
}

var preloadAliases = []string{"ld", "pl", "load"}
var preloadCmd = &cobra.Command{
	Use:                "preload",
	Short:              "将指定目录或文件加载到 alluxio 缓存系统中",
	Long:               `将指定目录或文件加载到 alluxio 缓存系统中，此操作可能导致其他最近未使用的数据被踢出 alluxio 缓存系统或从缓存系统中更高级别的缓存降级到更低级别的缓存中。`,
	Aliases:            preloadAliases,
	DisableSuggestions: false,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		mLoadFlags.Args = args
		return mLoadFlags.Validate()
	},
	Run: func(cmd *cobra.Command, args []string) {
		mCMD := &CMD{}
		mCMD.init(mLoadFlags)
		mCMD.start()
	},
}
