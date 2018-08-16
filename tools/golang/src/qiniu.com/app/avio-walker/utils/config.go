package utils

import (
	"qiniu.com/app/common"
)

type kafkaConfig struct {
	Brokers   []string `json:"brokers"`
	CertFile  string   `json:"certFile"`
	KeyFile   string   `json:"keyFile"`
	CaFile    string   `json:"caFile"`
	VerifySsl bool     `json:"verifySsl"`
	Topic     string   `json:"topic"`
}

type httpServerConfig struct {
	Port int `json:"port"`
}

type Config struct {
	Server  httpServerConfig     `json:"server"`
	DB      common.MongoConfig   `json:"db"`
	KAFKA   kafkaConfig          `json:"kafka"`
	Alluxio common.AlluxioConfig `json:"alluxio"`
	Log     common.Log           `json:"log"`
}
