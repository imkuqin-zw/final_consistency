package producer

import (
	"github.com/Shopify/sarama"
	"github.com/imkuqin-zw/final_consistency/basic/config"
	"sync"
)

var (
	initOnce    sync.Once
	errorOnce   sync.Once
	successOnce sync.Once
	p           Producer
)

type Config struct {
	Brokers []string `json:"brokers" yaml:"brokers"`
	Sync    bool     `json:"sync" yaml:"sync"`
	Retries uint32   `json:"retries" yaml:"retries"`
}

type Producer interface {
	Close() error
	Send(*sarama.ProducerMessage) error
}

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

func init() {
	initOnce.Do(func() {
		cfg := getConf()
		if cfg.Sync {
			p = newSyncProducer(cfg)
		} else {
			p = newAsyncProducer(cfg)
		}
	})
}

func GetProducer() Producer {
	return p
}

func SetErrHandle(handle errHandle) {
	errorOnce.Do(func() {
		if c, ok := p.(*AsyncProducer); ok {
			go c.errProcess(handle)
		}
	})
}

func SetSuccessHandle(handle sucHandle) {
	successOnce.Do(func() {
		if c, ok := p.(*AsyncProducer); ok {
			go c.successProcess(handle)
		}
	})
}
