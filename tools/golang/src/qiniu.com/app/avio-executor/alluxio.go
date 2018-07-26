package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	atypo "qiniu.com/app/avio-apiserver/typo"
	"qiniu.com/app/common"
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
	}

	e = c.updateMaster()
	return
}

func (a *AlluxioClient) updateMaster() error {
	if a.conf.MasterHost != "" && !a.conf.ZooKeeper.Enabled {
		a.currentMaster = a.conf.MasterHost
	} else {
		master, e := common.ZKElect(a.conf.ZooKeeper.Servers, a.conf.ZooKeeper.Election)
		if e != nil {
			return e
		}
		a.currentMaster = master
	}

	// TODO handle master host not reachable
	return nil
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
	data := make([]byte, 32)
retryOpen:
	openRes, e := a.client.Post(url, "application/json", r)
	if e == nil {
		_, e = openRes.Body.Read(data)
		openRes.Body.Close()
	}

	if e != nil {
		if timesToRetry > 0 {
			timesToRetry--
			goto retryOpen
		} else {
			log.Errorf("failed to open file from %s after 3 retries", path)
			return e
		}
	}

	fd, e := strconv.Atoi(string(data))

	url = fmt.Sprintf("http://%s/api/v1/streams/%d/read", a.currentProxy, fd)
	timesToRetry = 3
retryRead:
	readRes, e := a.client.Post(url, "application/x-www-form-urlencoded", nil)
	if e != nil {
		if timesToRetry > 0 {
			timesToRetry--
			goto retryRead
		} else {
			log.Errorf("failed to read file from %s with stream id %d after 3 retries, error: %v", path, fd, e)
			return e
		}
	}

	bf := make([]byte, 4096)
	n := 4096
	for n == 4096 && e == nil {
		n, e = readRes.Body.Read(bf)
	}
	readRes.Body.Close()

	url = fmt.Sprintf("http://%s/api/v1/streams/%d/close", a.currentProxy, fd)
	timesToRetry = 3
retryClose:
	_, e = a.client.Post(url, "application/x-www-form-urlencoded", nil)
	if e != nil {
		if timesToRetry > 0 {
			timesToRetry--
			goto retryClose
		} else {
			log.Errorf("failed to close stream %d from file %s after 3 retries, error: %v", fd, path, e)
			return e
		}
	}

	return nil
}

func (a *AlluxioClient) save(path string) error {
	masterToRetry := 3
retryMaster:
	url := "http://" + a.currentMaster + "/api/v1/master/file/schedule_async_persist/?path=" + path

	timesToRetry := 3
retry:
	_, e := a.client.Post(url, "application/x-www-form-urlencoded", nil)
	if e != nil {
		if timesToRetry > 0 {
			timesToRetry--
			goto retry
		} else if masterToRetry > 0 {
			masterToRetry--
			a.updateMaster()
			goto retryMaster
		} else {
			log.Errorf("failed to call alluxio master to schedule persist %s after 9 retries", path)
		}
	}

	return nil
}

func (a *AlluxioClient) Call(typ atypo.JobType, path string) error {
	switch typ {
	case atypo.PreloadJobType:
		return a.preload(path)
	case atypo.SaveJobType:
		return a.save(path)
	}
	return nil
}
