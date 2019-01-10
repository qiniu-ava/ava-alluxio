package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/pkg/errors"
	yaml "gopkg.in/yaml.v2"
	"qiniu.com/app/alluxio-exporter/collectors"
)

type AlluxioConfig struct {
	MasterHost string `yaml:"master_host"`
	Group      string `yaml:"group"`
}
type WorkerConfig struct {
	Address workerAddress `json:"address"`
}
type workerAddress struct {
	Host    string `json:"host"`
	WebPort int    `json:"webPort"`
	Role    string `json:"role"`
}

// Config is the top-level configuration for Metastord.
type Config struct {
	Alluxio []*AlluxioConfig
}

// fileExists returns true if the path exists and is a file.
func fileExists(path string) bool {
	stat, err := os.Stat(path)
	return !os.IsNotExist(err) && !stat.IsDir()
}

func ParseConfig(p string) (*Config, error) {
	if !fileExists(p) {
		return nil, errors.Errorf("Config file does not exist.")
	}

	cfgData, err := ioutil.ReadFile(p)
	if err != nil {
		return nil, errors.Errorf("ReadFile: %v", err)
	}

	var cfg Config
	err = yaml.Unmarshal(cfgData, &cfg)
	if err != nil {
		return nil, fmt.Errorf("yaml parse: %v", err)
	}

	return &cfg, nil
}

func GetWorkerInfoList(masterHost string) ([]WorkerConfig, error) {
	url := "http://" + masterHost + "/api/v1/master/worker_info_list"
	method := "GET"
	res, e := collectors.HTTPRequest(url, method, nil, nil)
	if e != nil {
		return nil, errors.Errorf("HTTPRequest for worker : %v ", e)
	}
	var workerList []WorkerConfig
	e = json.Unmarshal(res, &workerList)
	if e != nil {
		return nil, fmt.Errorf("json parse: %v", e)
	}
	return workerList, nil
}
