package service

import (
	"github.com/imkuqin-zw/final_consistency/msg_api/repository"
	"github.com/imkuqin-zw/final_consistency/msg_api/utils/snowflake"
	z "shop/plugins/zap"
	"sync"
)

var (
	once sync.Once
	log  *z.Logger
	s    *Service
	uuid *snowflake.Snowflake
	repo repository.Repository
)

type Service struct {
	version uint32
}

func Init() {
	once.Do(func() {
		log = z.GetLogger()
		uuid = snowflake.GetSnowkflak()
		repo = repository.GetRepo()
		s = &Service{}
	})
}

func GetService() *Service {
	return s
}
