package kafka

import (
	"context"
	"github.com/Shopify/sarama"
	"github.com/imkuqin-zw/final_consistency/plugins/kafka/producer"
	z "shop/plugins/zap"
	"sync"
)

var (
	once sync.Once
	log  *z.Logger
	p    producer.Producer
)

func Init() {
	once.Do(func() {
		log = z.GetLogger()
		p = producer.GetProducer()
	})
}

type Queue struct {
}

func (kq *Queue) SendMsg(ctx context.Context, topic string, msg string) error {
	m := &sarama.ProducerMessage{
		Topic:    topic,
		Metadata: ctx,
		Value:    sarama.StringEncoder(msg),
	}
	return p.Send(m)
}

func NewMsgQueue() *Queue {
	return &Queue{}
}
