package producer

import (
	"context"
	"github.com/Shopify/sarama"
	"time"
)

type AsyncProducer struct {
	sarama.AsyncProducer
}

func newAsyncProducer(cfg *Config) Producer {
	c := sarama.NewConfig()
	c.Producer.RequiredAcks = sarama.WaitForLocal     // Only wait for the leader to ack
	c.Producer.Compression = sarama.CompressionSnappy // Compress messages
	c.Producer.Return.Successes = true
	c.Producer.Return.Errors = true
	c.Producer.Flush.Frequency = 500 * time.Millisecond // Flush batches every 500ms
	for i := uint32(0); i < cfg.Retries; i++ {
		asyncProducer, err := sarama.NewAsyncProducer(cfg.Brokers, c)
		if err == nil {
			return &AsyncProducer{asyncProducer}
		}
		sarama.Logger.Printf("new async producer fault times(%d) error(%v)", i, err)
		time.Sleep(time.Second)
	}
	panic("[kafka] init async producer fault")
}

func (ap *AsyncProducer) Send(c context.Context, msg *sarama.ProducerMessage) error {
	msg.Metadata = c
	ap.Input() <- msg
	return nil
}

func (ap *AsyncProducer) isSync() bool {
	return false
}

func (ap *AsyncProducer) errProcess(deal errHandle) {
	err := ap.Errors()
	for {
		e, ok := <-err
		if !ok {
			return
		}
		sarama.Logger.Printf("kafka producer send message(%v) failed error(%v)", e.Msg, e.Err)
		if deal != nil {
			deal(e)
		}
	}
}

func (ap *AsyncProducer) successProcess(deal sucHandle) {
	suc := ap.Successes()
	for {
		msg, ok := <-suc
		if !ok {
			return
		}
		sarama.Logger.Printf("kafka producer send message(%v) sucsess", msg)
		if deal != nil {
			deal(msg)
		}
	}
}
