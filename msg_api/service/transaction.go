package service

import (
	"context"
	"github.com/imkuqin-zw/final_consistency/msg_api/models"
	"time"
)

//存储预发送消息
func (s *Service) StoreMsgWaitingConfirm(msg *models.TransactionMsg) error {
	msg.MsgId = uuid.GetUUID()
	msg.Version = s.version
	msg.Status = models.TransStatusWaitingConfirm
	msg.CreateAt = time.Now()
	msg.UpdateAt = msg.CreateAt
	msg.MsgSendTimes = 0
	msg.AlreadyDead = false

	return repo.InsertTransMsg(msg)
}

//确认预发送消息并发送消息
func (s *Service) ConfirmAndSendMessage(ctx context.Context, msgId string) error {
	msg, err := repo.GetTransMsgByMsgId(msgId)
	if err != nil {
		return err
	}
	var affect int64
	affect, err = repo.UpdateTransMsgStatusByMsgId(msg, models.TransStatusSending)
	if err != nil {
		return err
	}
	if affect == 0 {
		return nil
	}
	return msgQue.SendMsg(ctx, msg)
}

//查询并处理超时的预发送消息
func (s *Service) DealTimeOutMsgWaitingConfirm() {

}

//func (s *Service) DealUnSend
