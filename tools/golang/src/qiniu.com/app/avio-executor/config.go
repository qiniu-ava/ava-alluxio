package main

import (
	"qiniu.com/app/common"
)

type kafkaConfig struct {
	Address []string `json:"address"`
	Topic   string   `json:"topic"`
}

type Config struct {
	Kafka   kafkaConfig          `json:"kafka"`
	Alluxio common.AlluxioConfig `json:"alluxio"`
	Log     common.Log           `json:"log"`
}
