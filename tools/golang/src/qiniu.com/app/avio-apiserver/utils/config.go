package utils

type mongoConfig struct {
	Host            string `json:"host"`
	DBName          string `json:"dbname"`
	SessionPoolSize int    `json:"sessionPoolSize"`
}

type kubeConfig struct {
	KubeConfigPath string `json:"kube_config_path"`
	Namespace      string `json:"namespace"`
	test           bool   // just for ci
}

type Config struct {
	DB                   mongoConfig `json:"db"`
	CollSessionPoolLimit int         `json:"coll_session_pool_limit"`
	KubeConfig           kubeConfig  `json:"kube_config"`
}
