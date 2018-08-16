package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/imdario/mergo"
	"qiniu.com/app/avio-apiserver/compile"
	"qiniu.com/app/avio-apiserver/utils"
	"qiniu.com/app/common"
	"qiniu.com/app/common/database"
	"qiniupkg.com/x/config.v7"
	log "qiniupkg.com/x/log.v7"
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

func loadMongoAuthConfig() *common.MongoAuthConfig {
	var conf common.MongoAuthConfig

	if e := common.LoadJsonObjectFromFile(compile.MongoAuthConfigPath, &conf); e != nil {
		log.Warnf("load mongo auth filed, error: %v", e)
		return nil
	}

	return &conf
}

func Boot() {
	var defaultConf, fileConf, conf utils.Config
	config.LoadString(&defaultConf, defaultConfig)
	config.Init("f", appName, appName+".conf")

	if e := mergo.MergeWithOverwrite(&conf, defaultConf); e != nil {
		panic("merge options failed")
	}

	if e := config.Load(&fileConf); e != nil {
		panic("load configuration failed")
	}

	if e := mergo.MergeWithOverwrite(&conf, fileConf); e != nil {
		panic("merge options failed")
	}

	s, _ := json.Marshal(&conf)
	log.Infof("apiserver is about to start with following configuration:\n%s", string(s))

	if conf.Log.Debug {
		log.SetOutputLevel(log.Ldebug)
	} else {
		log.SetOutputLevel(log.Linfo)
	}

	auth := loadMongoAuthConfig()
	session, err := database.NewMongoSession(&conf.DB, auth)
	defer session.Close()
	if err != nil {
		panic("connect to mongodb failed")
	}

	if e := database.Init(); e != nil {
		panic(fmt.Sprintf("init mongodb failed, error: %v", e))
	}

	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", conf.Server.Port),
		Handler:      initRouter(),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	log.Infof("listening %d...", conf.Server.Port)
	server.ListenAndServe()
}
