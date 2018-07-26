package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/imdario/mergo"
	"github.com/pkg/errors"
	"qiniu.com/app/avio-walker/services"
	"qiniu.com/app/avio-walker/utils"
	"qiniu.com/app/avio-walker/walker"
	"qiniupkg.com/x/config.v7"
)

const defaultConfig = `{
	"db": {
    "host": "localhost:27017",
		"dbname": "avio",
		"sessionPoolSize": 200
	},
	"readTimeout": 10,
	"writeTimeout": 600,
  "log": {
    "debug": false
  }
}`

const appName = "avio-walker"

func initConf(conf *utils.Config) error {
	var defaultConf utils.Config
	if e := config.LoadString(&defaultConf, defaultConfig); e != nil {
		return errors.Errorf("load default configuration failed")
	}

	config.Init("f", appName, appName+".conf")

	if e := config.Load(conf); e != nil {
		return errors.Errorf("load configuration failed")
	}

	if e := mergo.MergeWithOverwrite(conf, defaultConf); e != nil {
		return errors.Errorf("merge options failed")
	}

	return nil
}

func main() {
	var conf utils.Config
	if e := initConf(&conf); e != nil {
		panic(e)
	}

	w, e := walker.NewWalker(&conf)
	if e != nil {
		panic(e)
	}

	w.SetServerRouter(services.InitRouter(w))

	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go (func(w walker.Walker) {
		<-c
		w.GracefullyShutdown()
		// os.Exit(1)
	})(w)

	w.Run()
}
