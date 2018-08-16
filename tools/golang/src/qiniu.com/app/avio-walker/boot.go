package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/imdario/mergo"
	"github.com/pkg/errors"
	"qiniu.com/app/avio-walker/compile"
	"qiniu.com/app/avio-walker/services"
	"qiniu.com/app/avio-walker/utils"
	"qiniu.com/app/avio-walker/walker"
	"qiniu.com/app/common"
	"qiniupkg.com/x/config.v7"
	log "qiniupkg.com/x/log.v7"
)

const defaultConfig = `{
	"db": {
    "host": "localhost:27017",
		"dbname": "avio",
		"sessionPoolSize": 200
	},
	"server": {
		"port": 8080
	},
	"readTimeout": 10,
	"writeTimeout": 600,
  "log": {
    "debug": false
  }
}`

const appName = "avio-walker"

func initConf(conf *utils.Config) error {
	var defaultConf, fileConf utils.Config
	if e := config.LoadString(&defaultConf, defaultConfig); e != nil {
		return errors.Errorf("load default configuration failed")
	}

	config.Init("f", appName, appName+".conf")

	if e := config.Load(&fileConf); e != nil {
		return errors.Errorf("load configuration failed")
	}

	if e := mergo.MergeWithOverwrite(conf, defaultConf); e != nil {
		return errors.Errorf("merge options failed")
	}

	if e := mergo.MergeWithOverwrite(conf, fileConf); e != nil {
		return errors.Errorf("merge options failed")
	}

	return nil
}

func Boot() {
	var conf utils.Config
	if e := initConf(&conf); e != nil {
		panic(e)
	}

	if conf.Log.Debug {
		log.SetOutputLevel(log.Ldebug)
	} else {
		log.SetOutputLevel(log.Linfo)
	}

	var auth *common.MongoAuthConfig
	if e := common.LoadJsonObjectFromFile(compile.MongoAuthConfigPath, auth); e != nil {
		log.Warnf("load mongo auth config failed, error: %v", e)
	}

	w, e := walker.NewWalker(&conf, auth)
	if e != nil {
		panic(e)
	}

	w.SetServerRouter(services.InitRouter(w))

	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go (func(w walker.Walker) {
		<-c
		w.GracefullyShutdown()
	})(w)

	w.Run()
}
