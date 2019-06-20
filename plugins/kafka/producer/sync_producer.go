package producer

import (
	"github.com/Shopify/sarama"
	"time"
)

type SyncProducer struct {
	sarama.SyncProducer
}

func newSyncProducer(cfg *Config) Producer {
	c := sarama.NewConfig()
	c.Producer.RequiredAcks = sarama.WaitForAll
	c.Producer.Retry.Max = 10
	for i := uint32(0); i < cfg.Retries; i++ {
		syncProducer, err := sarama.NewSyncProducer(cfg.Brokers, c)
		if err == nil {
			return &SyncProducer{syncProducer}
		}
		sarama.Logger.Printf("new sync producer fault times(%d) error(%v)", i, err)
		time.Sleep(time.Second)
	}
	panic("[kafka] init sync producer fault")
}

func (sp *SyncProducer) Send(msg *sarama.ProducerMessage) error {
	_, _, err := sp.SendMessage(msg)
	return err
}
