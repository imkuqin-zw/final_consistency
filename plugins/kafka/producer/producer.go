package producer

import (
	"context"
	"final_consistency/basic/config"
	"github.com/Shopify/sarama"
	"time"
)

type Config struct {
	Brokers []string `json:"brokers" yaml:"brokers"`
	Sync    bool     `json:"sync" yaml:"sync"`
	Retries uint32   `json:"retries" yaml:"retries"`
}

type Producer interface {
	isSync() bool
	Send(context.Context, *sarama.ProducerMessage) (err error)
}

type errHandle func(*sarama.ProducerError)
type sucHandle func(*sarama.ProducerMessage)

func getConf() *Config {
	c := config.C()
	cfg := &Config{}
	err := c.App("kafka", cfg)
	if err != nil {
		sarama.Logger.Printf("[kafka] %v", err)
		panic(err)
	}
	if cfg.Retries == 0 {
		cfg.Retries = 3
	}
	return cfg
}

func Input(c context.Context, msg *sarama.ProducerMessage) (err error) {
	if !p.conf.Sync {
		msg.Metadata = c
		p.AsyncProducer.Input() <- msg
	} else {
		if _, _, err = p.SyncProducer.SendMessage(msg); err != nil {
			sarama.Logger.Print("syncProducer send msg(%v) fault error(%v): ", msg, err)
		}
	}
	return
}

func Close() (err error) {
	if !p.conf.Sync {
		if p.AsyncProducer != nil {
			return p.AsyncProducer.Close()
		}
	}
	if p.SyncProducer != nil {
		return p.SyncProducer.Close()
	}
	return
}
