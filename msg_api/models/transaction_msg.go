package models

import (
	"time"
)

const (
	TransStatusWaitingConfirm = iota
	TransStatusSending
	TransStatusFanish
)

type TransactionMsg struct {
	Id          uint64      `json:"id"`            //主键
	MsgId       uint64      `json:"msg_id"`        //消息id
	Version     uint32      `json:"version"`       //版本号
	Status      uint32      `json:"status"`        //状态
	Editor      uint64      `json:"editor"`        //修改者
	Creator     uint64      `json:"creator"`       //创建者
	MsgDataType uint64      `json:"msg_data_type"` //消息数据类型
	MsgBody     string      `json:"msg_body"`      //消息内容
	UpdateAt    time.Time   `json:"update_at"`     //最后修改时间
	CreateAt    time.Time   `json:"create_at"`     //创建时间
	ConsumerQue string      `json:"consumer_que"`  //消息队列
	SendTimes   uint8       `json:"send_times"`    //消息重发次数
	AlreadyDead bool        `json:"already_dead"`  //是否死亡
	Remark      string      `json:"remark"`        //备注
	Extension   interface{} `json:"extension"`     //扩展
}
