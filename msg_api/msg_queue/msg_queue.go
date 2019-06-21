package msg_queue

import (
	"context"
	"github.com/imkuqin-zw/final_consistency/basic/config"
	"github.com/imkuqin-zw/final_consistency/msg_api/msg_queue/kafka"
	z "github.com/imkuqin-zw/final_consistency/plugins/zap"
	"go.uber.org/zap"
	"sync"
)

var (
	once sync.Once
	p    MsgQue
	log  *z.Logger
)

const (
	_ = iota
	KafkaQueue
)

type RepoConf struct {
	DbType int8 `json:"db_type"`
}

type MsgQue interface {
	SendMsg(context.Context, string, string) error
}

func Init() {
	once.Do(func() {
		log = z.GetLogger()
		InitRepository()
	})
}

func InitRepository() {
	c := config.C()
	var cfg RepoConf
	if err := c.Path("repo", &cfg); err != nil {
		log.Panic("get app queue type config fault", zap.Error(err))
	}
	switch cfg.DbType {
	case KafkaQueue:
		kafka.Init()
		p = kafka.NewMsgQueue()
	default:
		log.Panic("unknown msg queue type", zap.Any("queue_cfg", cfg))
	}
}

func GetMsgQue() MsgQue {
	return p
}
