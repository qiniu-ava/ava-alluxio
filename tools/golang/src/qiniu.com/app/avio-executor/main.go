package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/imdario/mergo"
	"github.com/pkg/errors"
	"qiniupkg.com/x/config.v7"
)

const defaultConfig = `{
	"kafka": {
    "address": ["0.0.0.0:9092"],
    "topic": "avio_job_item_prd"
  }
}`

const appName = "avio-executor"

func initConf(conf *Config) error {
	var defaultConf Config
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
	var conf Config
	if e := initConf(&conf); e != nil {
		panic(e)
	}

	consumer, e := NewConsumer(&conf)
	if e != nil {
		// TODO fix me
		panic(e)
	}

	consumer.HandleMessage()
	ch := make(chan os.Signal)
	signal.Notify(ch, os.Interrupt, syscall.SIGTERM)
	go (func(c *Consumer) {
		<-ch
		c.GracefullyShutdown()
		// os.Exit(1)
	})(consumer)
}
