package model

import (
	"IM/pkg/db"
	"fmt"
	"time"
)

type Message struct {
	ID          int64     `gorm:"primary_key;auto_increment;comment:'自增主键'" json:"id"`
	UserID      int64     `gorm:"not null;comment:'用户id，指接受者用户id'" json:"user_id"`
	SenderID    int64     `gorm:"not null;comment:'发送者用户id'" json:"sender_id"`
	SessionType int8      `gorm:"not null;comment:'聊天类型，群聊/单聊'" json:"session_type"` // 1：单聊 2：群聊
	ReceiverId  int64     `gorm:"not null;comment:'接收者id，群聊id/用户id'" json:"receiver_id"`
	MessageType int8      `gorm:"not null;comment:'消息类型,文字、语音、图片'" json:"message_type"`
	Content     string    `gorm:"not null;comment:'消息内容'" json:"content"`
	Seq         int64     `gorm:"not null;comment:'消息序列号'" json:"seq"`
	SendTime    time.Time `gorm:"not null;default:CURRENT_TIMESTAMP;comment:'消息发送时间'" json:"send_time"`
	CreateTime  time.Time `gorm:"not null;default:CURRENT_TIMESTAMP;comment:'创建时间'" json:"create_time"`
	UpdateTime  time.Time `gorm:"not null;default:CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP;comment:'更新时间'" json:"update_time"`
}

func (Message) TableName() string {
	return "message"
}

func CreateMessages(messages []Message) (err error) {
	result := db.DB.Create(messages)
	if err = result.Error; err != nil {
		return err
	}
	fmt.Println("消息保存成功")
	return nil
}
