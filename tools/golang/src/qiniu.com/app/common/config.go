package common

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

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

type MongoAuthConfig struct {
	UserName string `json:"user"`
	Password string `json:"password"`
}

func LoadJsonObjectFromFile(path string, target interface{}) error {
	f, e := os.Open(path)
	if e != nil {
		return e
	}

	defer f.Close()
	b, e := ioutil.ReadAll(f)
	if e != nil {
		return e
	}

	return json.Unmarshal(b, target)
}
