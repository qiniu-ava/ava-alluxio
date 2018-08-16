package main

import (
	"encoding/json"
	"os"
	"sync"
	"time"

	"qiniu.com/app/common/typo"

	kafka "gopkg.in/Shopify/sarama.v1"
	log "qiniupkg.com/x/log.v7"
)

type Consumer struct {
	topic              string
	kafkaClient        kafka.Client
	preloadExecutor    Executor
	saveExecutor       Executor
	shutDownWg         sync.WaitGroup
	setToShutDown      bool
	lastTimeToShutDown time.Time
}

var durationToShutDownForce = 5 * time.Second

func NewConsumer(conf *Config) (c *Consumer, e error) {
	// TODO fix me
	config := kafka.NewConfig()
	client, e := kafka.NewClient(conf.Kafka.Address, config)
	if e != nil {
		return nil, e
	}
	aClient, e := NewAlluxioClient(&conf.Alluxio)
	if e != nil {
		return nil, e
	}
	c = &Consumer{kafkaClient: client}
	c.topic = conf.Kafka.Topic
	c.preloadExecutor = NewPreloadExecutor(aClient)
	c.saveExecutor = NewSaveExecutor(aClient)
	c.shutDownWg = sync.WaitGroup{}
	c.setToShutDown = false
	c.lastTimeToShutDown = time.Date(1970, 0, 0, 0, 0, 0, 0, time.Local)
	return
}

func (c *Consumer) HandleMessage() error {
	// TODO implement me
	pList, e := c.kafkaClient.Partitions(c.topic)
	if e != nil {
		return e
	}

	log.Infof("partitions list length: %d", len(pList))

	consumer, e := kafka.NewConsumerFromClient(c.kafkaClient)
	if e != nil {
		return e
	}

	messages := make(chan *kafka.ConsumerMessage, 64)

	for _, partition := range pList {
		offset, e := c.kafkaClient.GetOffset(c.topic, partition, kafka.OffsetNewest)
		if e != nil {
			return e
		}

		log.Infof("partition %d offset %d", partition, offset)

		pConsumer, e := consumer.ConsumePartition(c.topic, partition, offset)
		if e != nil {
			return e
		}

		go func(pc kafka.PartitionConsumer) {
			for msg := range pConsumer.Messages() {
				log.Infof("pc message: %v, about to send to messages", msg)
				c.shutDownWg.Add(1)
				messages <- msg
			}
		}(pConsumer)
	}

	go func() {
		for {
			alluxioMsg := typo.KafkaMessage{}
			select {
			case m := <-messages:
				e := json.Unmarshal(m.Value, &alluxioMsg)
				if e != nil {
					log.Errorf("unmarshal kafka message from kafka failed, error: %v", e)
				} else {
					c.handleKafkaMessage(alluxioMsg)
				}
				c.shutDownWg.Done()
			}
		}
	}()

	for {
		if c.setToShutDown {
			c.shutDownWg.Wait()
			os.Exit(ExitByGracefullyShutDown)
		}

		time.Sleep(durationToShutDownForce)
	}

	return nil
}

func (c *Consumer) handleKafkaMessage(msg typo.KafkaMessage) {
	switch string(msg.MsgType) {
	case string(typo.AvioCMDMsg):
		c.handleAvioCMDMsg(msg.AvioCMDData)
	}
}

func (c *Consumer) handleAvioCMDMsg(data typo.AvioCMDData) {
	log.Infof("about to handle avio cmd message: job type: %d, job path: %s", data.JobType, data.Path)
	switch int(data.JobType) {
	case int(typo.PreloadJobType):
		c.preloadExecutor.Do(data.Path)
	case int(typo.SaveJobType):
		c.saveExecutor.Do(data.Path)
	}
}

func (c *Consumer) GracefullyShutdown() {
	c.close()
}

func (c *Consumer) close() error {
	if c.kafkaClient != nil && !c.kafkaClient.Closed() {
		c.kafkaClient.Close()
	}

	if time.Since(c.lastTimeToShutDown) < durationToShutDownForce {
		log.Warn("force shutdown immediately")
		os.Exit(ExitByForceShutDown)
	} else if c.setToShutDown {
		log.Warnf("you are trying to shut down executor while it's shutting down gracefully,"+
			" you can shut down forcely by sending multi TERM signal within %v", durationToShutDownForce)
	} else {
		c.setToShutDown = true
		log.Warn("shutting down gracefully")
	}

	c.lastTimeToShutDown = time.Now()

	return nil
}
