package main

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
)

type AlluxioConfig struct {
	AlluxioType    string `yaml:"type"`
	AlluxioWebHost string `yaml:"host"`
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
