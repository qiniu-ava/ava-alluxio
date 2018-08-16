package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"

	"qiniu.com/app/common"
	"qiniu.com/app/common/typo"
	log "qiniupkg.com/x/log.v7"
)

type OpenFileOption struct {
	ReadType                   string `json:"readType"`
	LocationPolicyClass        string `json:"locationPolicyClass,omitempty"`
	CacheLocationPolicyClass   string `json:"cacheLocationPolicyClass,omitempty"`
	UfsReadLocationPolicyClass string `json:"ufsReadLocationPolicyClass,omitempty"`
	MaxUfsReadConcurrency      string `json:"maxUfsReadConcurrency,omitempty"`
}

type AlluxioClient struct {
	conf          *common.AlluxioConfig
	currentMaster string
	currentProxy  string
	client        *http.Client
}

func NewAlluxioClient(conf *common.AlluxioConfig) (c *AlluxioClient, e error) {
	c = &AlluxioClient{
		conf: conf,
		client: &http.Client{
			Timeout: time.Second * 60,
		},
		currentMaster: conf.MasterHost,
		currentProxy:  conf.ProxyHost,
	}

	return
}

func (a *AlluxioClient) preload(path string) error {
	url := "http://" + a.currentProxy + "/api/v1/paths/" + path + "/open-file/"
	timesToRetry := 3
	options := OpenFileOption{
		ReadType: "CACHE_PROMOTE",
	}
	d, e := json.Marshal(options)
	if e != nil {
		log.Warnf("we should never be here, %v", e)
		return e
	}

	r := bytes.NewReader(d)
	var data []byte
retryOpen:
	openRes, e := a.client.Post(url, "application/json", r)
	if e != nil {
		log.Warnf("open file from %s failed, error: %v", path, e)
	}

	if e == nil {
		data = make([]byte, openRes.ContentLength)
		_, e = openRes.Body.Read(data)
		openRes.Body.Close()
	}

	if e != nil && e != io.EOF {
		if timesToRetry > 0 {
			timesToRetry--
			time.Sleep(time.Duration((3-timesToRetry)*2) * time.Second)
			goto retryOpen
		} else {
			log.Errorf("failed to open file from %s after 3 retries, error: %v", path, e)
			return e
		}
	}

	fd, e := strconv.Atoi(string(data))

	if e != nil {
		log.Warnf("failed to get stream id from response body, error: %v", e)
	} else {
		log.Debugf("get stream id %d for path %s", fd, path)
	}

	url = fmt.Sprintf("http://%s/api/v1/streams/%d/read", a.currentProxy, fd)
	timesToRetry = 3
retryRead:
	readRes, e := a.client.Post(url, "application/x-www-form-urlencoded", nil)
	if e != nil {
		if timesToRetry > 0 {
			timesToRetry--
			time.Sleep(time.Duration((3-timesToRetry)*2) * time.Second)
			goto retryRead
		} else {
			log.Errorf("failed to read file from %s with stream id %d after 3 retries, error: %v", path, fd, e)
			return e
		}
	}

	bf := make([]byte, 4096)
	n := 4096
	for n == 4096 && e == nil {
		log.Debugf("read stream from %s again", path)
		n, e = readRes.Body.Read(bf)
	}

	log.Debugf("after read stream, n: %d, error: %v", n, e)

	readRes.Body.Close()

	url = fmt.Sprintf("http://%s/api/v1/streams/%d/close", a.currentProxy, fd)
	timesToRetry = 3
retryClose:
	_, e = a.client.Post(url, "application/x-www-form-urlencoded", nil)
	if e != nil {
		if timesToRetry > 0 {
			timesToRetry--
			time.Sleep(time.Duration((3-timesToRetry)*2) * time.Second)
			goto retryClose
		} else {
			log.Errorf("failed to close stream %d from file %s after 3 retries, error: %v", fd, path, e)
			return e
		}
	}

	return nil
}

func (a *AlluxioClient) save(path string) error {
	url := "http://" + a.currentMaster + "/api/v1/master/file/schedule_async_persist/?path=" + path
	timesToRetry := 3
retry:
	// TODO post with config
	_, e := a.client.Post(url, "application/x-www-form-urlencoded", nil)
	if e != nil {
		if timesToRetry > 0 {
			timesToRetry--
			time.Sleep(time.Duration((3-timesToRetry)*2) * time.Second)
			goto retry
		} else {
			log.Errorf("failed to call alluxio master to schedule persist %s after 9 retries", path)
		}
	}

	return nil
}

func (a *AlluxioClient) Call(typ typo.JobType, path string) error {
	switch typ {
	case typo.PreloadJobType:
		return a.preload(path)
	case typo.SaveJobType:
		return a.save(path)
	}
	return nil
}
