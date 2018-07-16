package main

import (
	"net/http"
	"time"

	"github.com/imdario/mergo"
	db "qiniu.com/app/avio-apiserver/database"
	"qiniu.com/app/avio-apiserver/utils"
	"qiniupkg.com/x/config.v7"
)

const defaultConfig = `{
	"db": {
    "host": "localhost:27017",
    "dbname": "avio"
  },
  "log": {
    "debug": false
  }
}`

const appName = "avio-apiserver"

func main() {
	var defaultConf, conf utils.Config
	config.LoadString(&defaultConf, defaultConfig)
	config.Init("f", appName, appName+".conf")

	if e := config.Load(&conf); e != nil {
		panic("load configuration failed")
	}

	if e := mergo.MergeWithOverwrite(&conf, defaultConf); e != nil {
		panic("merge options failed")
	}

	session, err := db.NewMongoSession(&conf)
	defer session.Close()
	if err != nil {
		panic("connect to mongodb failed")
	}

	db.Init()

	server := &http.Server{
		Addr:         ":8080",
		Handler:      initRouter(),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	server.ListenAndServe()
}
