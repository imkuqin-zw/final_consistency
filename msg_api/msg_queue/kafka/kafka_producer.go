package kafka

import (
	"context"
	"encoding/json"
	"github.com/Shopify/sarama"
	"github.com/imkuqin-zw/final_consistency/msg_api/msg_queue"
	"github.com/imkuqin-zw/final_consistency/plugins/kafka/producer"
	z "shop/plugins/zap"
	"sync"
)

var (
	once sync.Once
	log  *z.Logger
	p    producer.Producer
)

const (
	TransactionMsgTopic = "transaction_msg_topic"
)

func Init() {
	once.Do(func() {
		log = z.GetLogger()
		p = producer.GetProducer()
	})
}

type KafkaQueue struct {
}

func (kq *KafkaQueue) SendMsg(ctx context.Context, msg interface{}) error {
	val, err := json.Marshal(msg)
	if err != nil {
		return err
	}
	m := &sarama.ProducerMessage{
		Topic:    TransactionMsgTopic,
		Metadata: ctx,
		Value:    sarama.ByteEncoder(val),
	}
	return p.Send(m)
}

func NewMsgQueue() msg_queue.MsgQue {
	return &KafkaQueue{}
}
