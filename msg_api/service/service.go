package service

import (
	"github.com/imkuqin-zw/final_consistency/msg_api/msg_queue"
	"github.com/imkuqin-zw/final_consistency/msg_api/repository"
	"github.com/imkuqin-zw/final_consistency/msg_api/utils/snowflake"
	z "shop/plugins/zap"
	"sync"
)

var (
	once   sync.Once
	log    *z.Logger
	s      *Service
	uuid   *snowflake.Snowflake
	repo   repository.Repository
	msgQue msg_queue.MsgQue
)

type Service struct {
	version uint32
}

func Init(version uint32) {
	once.Do(func() {
		log = z.GetLogger()
		uuid = snowflake.GetSnowkflak()
		repo = repository.GetRepo()
		s = &Service{version: version}
	})
}

func GetService() *Service {
	return s
}
