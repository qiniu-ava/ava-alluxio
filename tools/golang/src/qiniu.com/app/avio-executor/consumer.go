package main

import (
	kafka "gopkg.in/Shopify/sarama.v1"
)

type Consumer struct {
	kafkaClient     kafka.Consumer
	preloadExecutor Executor
	saveExecutor    Executor
}

func NewConsumer(conf *Config) (c *Consumer, e error) {
	// TODO fix me
	config := kafka.NewConfig()
	client, e := kafka.NewConsumer(conf.Kafka.Address, config)
	if e != nil {
		return nil, e
	}
	aClient, e := NewAlluxioClient(&conf.Alluxio)
	if e != nil {
		return nil, e
	}
	c = &Consumer{kafkaClient: client}
	c.preloadExecutor = NewPreloadExecutor(aClient)
	c.saveExecutor = NewSaveExecutor(aClient)
	return
}

func (c *Consumer) HandleMessage() error {
	// TODO implement me
	return nil
}

func (c *Consumer) GracefullyShutdown() {
	c.close()
}

func (c *Consumer) close() error {
	if c.kafkaClient != nil {
		return c.kafkaClient.Close()
	}

	return nil
}
