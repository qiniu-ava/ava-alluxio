package utils

type mongoConfig struct {
	Host            string `json:"host"`
	DBName          string `json:"dbname"`
	SessionPoolSize int    `json:"sessionPoolSize"`
}

type kafkaConfig struct {
	Brokers   []string `json:"brokers"`
	CertFile  string   `json:"certFile"`
	KeyFile   string   `json:"keyFile"`
	CaFile    string   `json:"caFile"`
	VerifySsl bool     `json:"verifySsl"`
	Topic     string   `json:"topic"`
}

type Config struct {
	DB                   mongoConfig `json:"db"`
	KAFKA                kafkaConfig `json:"kafka"`
	CollSessionPoolLimit int         `json:"coll_session_pool_limit"`
	ReadTimeout          int         `json:"readTimeout"`
	WriteTimeout         int         `json:"writeTimeout"`
}
