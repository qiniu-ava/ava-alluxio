package main

import (
	"encoding/json"
	"os"
	"os/signal"
	"syscall"

	"github.com/imdario/mergo"
	"github.com/pkg/errors"
	"qiniupkg.com/x/config.v7"
	log "qiniupkg.com/x/log.v7"
)

const defaultConfig = `{
	"kafka": {
    "address": ["0.0.0.0:9092"],
    "topic": "avio_job_item_prd"
  },
  "alluxio": {
    "masterHost": "0.0.0.0:19999",
    "proxyHost": "0.0.0.0:39999"
  }
}`

const appName = "avio-executor"

func initConf(conf *Config) error {
	var defaultConf, fileConf Config
	if e := config.LoadString(&defaultConf, defaultConfig); e != nil {
		return errors.Errorf("load default configuration failed, error: %v", e)
	}

	config.Init("f", appName, appName+".conf")

	if e := config.Load(&fileConf); e != nil {
		return errors.Errorf("load configuration failed, error: %v", e)
	}

	if e := mergo.MergeWithOverwrite(conf, defaultConf); e != nil {
		return errors.Errorf("merge options failed, erorr: %v", e)
	}

	if e := mergo.MergeWithOverwrite(conf, fileConf); e != nil {
		return errors.Errorf("merge options failed, error: %v", e)
	}

	return nil
}

func Boot() {
	var conf Config
	if e := initConf(&conf); e != nil {
		panic(e)
	}

	if conf.Log.Debug {
		log.SetOutputLevel(log.Ldebug)
	} else {
		log.SetOutputLevel(log.Linfo)
	}

	b, e := json.MarshalIndent(conf, "", "  ")
	log.Infof("running avio-executor with following configure:\n%v", string(b))

	consumer, e := NewConsumer(&conf)
	if e != nil {
		// TODO fix me
		panic(e)
	}

	ch := make(chan os.Signal)
	signal.Notify(ch, os.Interrupt, syscall.SIGTERM)
	go (func(c *Consumer) {
		for {
			<-ch
			c.GracefullyShutdown()
		}
	})(consumer)
	consumer.HandleMessage()
}
