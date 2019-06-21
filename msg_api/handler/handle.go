package handler

import (
	"github.com/imkuqin-zw/final_consistency/msg_api/service"
	z "github.com/imkuqin-zw/final_consistency/plugins/zap"
)

var log *z.Logger
var s *service.Service

func Init() {
	log = z.GetLogger()
	s = service.GetService()
}
