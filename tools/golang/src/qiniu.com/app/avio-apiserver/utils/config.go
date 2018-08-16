package utils

import (
	"qiniu.com/app/common"
)

type kubeConfig struct {
	KubeConfigPath string `json:"kube_config_path"`
	Namespace      string `json:"namespace"`
	test           bool   // just for ci
}

type httpServerConfig struct {
	Port int `json:"port"`
}

type Config struct {
	Server               httpServerConfig   `json:"server"`
	DB                   common.MongoConfig `json:"db"`
	CollSessionPoolLimit int                `json:"coll_session_pool_limit"`
	KubeConfig           kubeConfig         `json:"kube_config"`
	Log                  common.Log         `json:"log"`
}
