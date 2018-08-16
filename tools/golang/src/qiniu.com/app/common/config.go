package common

type MongoConfig struct {
	Host            string `json:"host"`
	DBName          string `json:"dbname"`
	SessionPoolSize int    `json:"sessionPoolSize"`
}

type Log struct {
	Debug bool `json:"debug"`
}

type AlluxioConfig struct {
	MasterHost string `json:"masterHost"`
	ProxyHost  string `json:"proxyHost"`
}
