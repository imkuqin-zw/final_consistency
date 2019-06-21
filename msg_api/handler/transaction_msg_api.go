package handler

import (
	"context"
	"github.com/imkuqin-zw/final_consistency/msg_api/models"
	pb "github.com/imkuqin-zw/final_consistency/msg_api/proto/transaction"
	"github.com/micro/go-micro/errors"
	"go.uber.org/zap"
)

type TransactionMsgApi struct {
}

func (c TransactionMsgApi) PingPong(ctx context.Context, stream pb.Transaction_PingPongStream) error {
	for {
		req, err := stream.Recv()
		if err != nil {
			return err
		}
		log.Debug("Got ping ", zap.Int64("stroke", req.Stroke))
		if err := stream.Send(&pb.Pong{Stroke: req.Stroke}); err != nil {
			return err
		}
	}
}

func (c TransactionMsgApi) CreateTransactionMsg(ctx context.Context, in *pb.ReqCreateTransMsg, out *pb.MsgId) error {
	if in.MsgBody == "" {
		return errors.BadRequest("", "%s", "msg body is empty")
	}
	if in.ConsumerQue == "" {
		return errors.BadRequest("", "%s", "consumer queue is empty")
	}
	msg := &models.TransactionMsg{
		Editor:      in.Editor,
		Creator:     in.Creator,
		MsgDataType: in.MsgDataType,
		MsgBody:     in.MsgBody,
		ConsumerQue: in.ConsumerQue,
		Remark:      in.Remark,
		Extension:   in.Extension,
	}
	if err := s.StoreMsgWaitingConfirm(msg); err != nil {
		log.Error("StoreMsgWaitingConfirm", zap.Error(err))
		return err
	}
	return nil
}

func (c TransactionMsgApi) ConfirmAndSendMessage(ctx context.Context, in *pb.MsgId, out *pb.NoContent) error {
	if in.Value == "" {
		return errors.BadRequest("", "%s", "msg_id is empty")
	}
	if err := s.ConfirmAndSendMessage(ctx, in.Value); err != nil {
		log.Error("ConfirmAndSendMessage", zap.Error(err))
		return err
	}
	return nil
}

func (c TransactionMsgApi) DeleteMessageByMsgId(ctx context.Context, in *pb.MsgId, out *pb.NoContent) error {
	if in.Value == "" {
		return errors.BadRequest("", "%s", "msg_id is empty")
	}
	if err := s.DeleteMessageByMsgId(in.Value); err != nil {
		log.Error("DeleteMessageByMsgId", zap.Error(err))
		return err
	}
	return nil
}

func (c TransactionMsgApi) ResendMessage(ctx context.Context, in *pb.MsgId, out *pb.NoContent) error {
	if in.Value == "" {
		return errors.BadRequest("", "%s", "msg_id is empty")
	}
	if err := s.ResendMessage(ctx, in.Value); err != nil {
		log.Error("ResendMessage", zap.Error(err))
		return err
	}
	return nil
}

func (c TransactionMsgApi) SetMessageToAlreadyDead(ctx context.Context, in *pb.MsgId, out *pb.NoContent) error {
	if in.Value == "" {
		return errors.BadRequest("", "%s", "msg_id is empty")
	}
	if err := s.SetMessageToAlreadyDead(in.Value); err != nil {
		log.Error("SetMessageToAlreadyDead", zap.Error(err))
		return err
	}
	return nil
}

func (c TransactionMsgApi) ConfirmMessageSuccess(ctx context.Context, in *pb.MsgId, out *pb.NoContent) error {
	if in.Value == "" {
		return errors.BadRequest("", "%s", "msg_id is empty")
	}
	if err := s.ConfirmMessageSuccess(in.Value); err != nil {
		log.Error("SetMessageToAlreadyDead", zap.Error(err))
		return err
	}
	return nil
}
