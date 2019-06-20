package consumer

import (
	"final_consistency/basic/config"
	"github.com/Shopify/sarama"
	"github.com/bsm/sarama-cluster"
	"time"
)

type errHandle func(error)
type ntfHandle func(notification *cluster.Notification)

type Config struct {
	Brokers  []string `json:"brokers" yaml:"brokers"`
	Topics   []string `json:"topics" yaml:"topics"`
	GroupId  string   `json:"group_id" yaml:"group_id"`
	Offset   bool     `json:"offset" yaml:"offset"`
	Retries  uint32   `json:"retries" yaml:"retries"`
	LogDebug bool     `json:"log_debug" yaml:"log_debug"`
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

func newClusterConsumer() *cluster.Consumer {
	cfg := getConf()
	clusterCfg := cluster.NewConfig()
	clusterCfg.Consumer.Return.Errors = true
	clusterCfg.Group.Return.Notifications = true
	if cfg.Offset {
		clusterCfg.Consumer.Offsets.Initial = sarama.OffsetNewest
	} else {
		clusterCfg.Consumer.Offsets.Initial = sarama.OffsetOldest
	}
	for i := uint32(0); i < cfg.Retries; i++ {
		s, err := cluster.NewConsumer(cfg.Brokers, cfg.GroupId, cfg.Topics, clusterCfg)
		if err == nil {
			return s
		}
		sarama.Logger.Printf("[kafka] new cluster consumer fault times(%d) error(%v)", i, err)
		time.Sleep(time.Second)
	}
	panic("[kafka] init cluster consumer fault")
}

func notificationProcess(deal ntfHandle) {
	if s == nil {
		sarama.Logger.Printf("[kafka] cluster consumer not initialized yet")
		return
	}
	notify := s.Notifications()
	for {
		ntf, ok := <-notify
		if !ok {
			return
		}
		sarama.Logger.Printf("[kafka] cluster consumer notification(%v)", ntf)
		if deal != nil {
			deal(ntf)
		}
	}
}

func errProcess(deal errHandle) {
	if s == nil {
		return
	}
	cfg := getConf()
	err := s.Errors()
	for {
		e, ok := <-err
		if !ok {
			break
		}
		sarama.Logger.Printf("[kafka] cluster consumer group_id(%d) error(%v)", cfg.GroupId, e)
		if deal != nil {
			deal(e)
		}
	}
}
