package service

import (
	"context"
	"github.com/imkuqin-zw/final_consistency/msg_api/models"
	"github.com/jinzhu/gorm"
	"time"
)

//存储预发送消息
func (s *Service) StoreMsgWaitingConfirm(msg *models.TransactionMsg) error {
	msg.MsgId = uuid.GetUUID()
	msg.Version = s.version
	msg.Status = models.TransStatusWaitingConfirm
	msg.CreateAt = time.Now()
	msg.UpdateAt = msg.CreateAt
	msg.SendTimes = 0
	msg.AlreadyDead = false

	return repo.InsertTransMsg(msg)
}

//确认预发送消息并发送消息
func (s *Service) ConfirmAndSendMessage(ctx context.Context, msgId string) error {
	msg, err := repo.GetTransMsgByMsgId(msgId)
	if err != nil {
		return err
	}
	affect, err := repo.UpdateTransMsgStatusByMsgId(msg, models.TransStatusSending)
	if err != nil && err != gorm.ErrRecordNotFound {
		return err
	}
	if affect == 0 {
		return nil
	}
	return msgQue.SendMsg(ctx, msg.ConsumerQue, msg.MsgBody)
}

//删除消息
func (s *Service) DeleteMessageByMsgId(msgId string) error {
	_, err := repo.DeleteTransMsgByMsgId(msgId)
	return err
}

//重新发送消息
func (s *Service) ResendMessage(ctx context.Context, msgId string) error {
	msg, err := repo.GetTransMsgByMsgId(msgId)
	if err != nil {
		return err
	}
	affect, err := repo.UpdateTransMsgSendTimesByMsgId(msgId)
	if err != nil {
		return err
	}
	if affect == 0 {
		return nil
	}
	return msgQue.SendMsg(ctx, msg.ConsumerQue, msg.MsgBody)
}

//设置消息已经死亡
func (s *Service) SetMessageToAlreadyDead(msgId string) error {
	_, err := repo.SetTransMsgAlreadyDeadByMsgId(msgId)
	return err
}

//确认消息已被成功消费
func (s *Service) ConfirmMessageSuccess(msgId string) error {
	msg, err := repo.GetTransMsgByMsgId(msgId)
	if err != nil {
		return err
	}
	_, err = repo.UpdateTransMsgStatusByMsgId(msg, models.TransStatusFanish)
	if err != nil && err != gorm.ErrRecordNotFound {
		return err
	}
	return nil
}
