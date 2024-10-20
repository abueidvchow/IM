package model

import (
	"IM/pkg/db"
	"IM/pkg/protocol/pb"
	"fmt"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
	"time"
)

type Message struct {
	ID          int64     `gorm:"primary_key;auto_increment;comment:'自增主键'" json:"id"`
	UserID      int64     `gorm:"not null;comment:'用户id，指接受者用户id'" json:"user_id"`
	SenderID    int64     `gorm:"not null;comment:'发送者用户id'" json:"sender_id"`
	SessionType int8      `gorm:"not null;comment:'聊天类型，群聊/单聊'" json:"session_type"` // 1：单聊 2：群聊
	ReceiverId  int64     `gorm:"not null;comment:'接收者id，群聊id/用户id'" json:"receiver_id"`
	MessageType int8      `gorm:"not null;comment:'消息类型,文字、语音、图片'" json:"message_type"`
	Content     []byte    `gorm:"not null;comment:'消息内容'" json:"content"`
	Seq         int64     `gorm:"not null;comment:'消息序列号'" json:"seq"`
	SendTime    time.Time `gorm:"not null;default:CURRENT_TIMESTAMP;comment:'消息发送时间'" json:"send_time"`
	CreateTime  time.Time `gorm:"not null;default:CURRENT_TIMESTAMP;comment:'创建时间'" json:"create_time"`
	UpdateTime  time.Time `gorm:"not null;default:CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP;comment:'更新时间'" json:"update_time"`
}

func (Message) TableName() string {
	return "message"
}

func ProtoMarshalToMessage(data []byte) []*Message {
	var messages []*Message
	mqMessages := &pb.MQMessages{}
	err := proto.Unmarshal(data, mqMessages)
	if err != nil {
		fmt.Println("json.Unmarshal(mqMessages) 失败,err:", err)
		return nil
	}
	for _, mqMessage := range mqMessages.Messages {
		message := &Message{
			UserID:      mqMessage.UserId,
			SenderID:    mqMessage.SenderId,
			SessionType: int8(mqMessage.SessionType),
			ReceiverId:  mqMessage.ReceiverId,
			MessageType: int8(mqMessage.MessageType),
			Content:     mqMessage.Content,
			Seq:         mqMessage.Seq,
			SendTime:    mqMessage.SendTime.AsTime(),
		}
		messages = append(messages, message)
	}
	return messages
}

func MessageToProtoMarshal(messages ...*Message) []byte {
	if len(messages) == 0 {
		return nil
	}
	var mqMessage []*pb.MQMessage
	for _, message := range messages {
		mqMessage = append(mqMessage, &pb.MQMessage{
			UserId:      message.UserID,
			SenderId:    message.SenderID,
			SessionType: int32(message.SessionType),
			ReceiverId:  message.ReceiverId,
			MessageType: int32(message.MessageType),
			Content:     message.Content,
			Seq:         message.Seq,
			SendTime:    timestamppb.New(message.SendTime),
		})
	}
	bytes, err := proto.Marshal(&pb.MQMessages{Messages: mqMessage})
	if err != nil {
		fmt.Println("json.Marshal(messages) 失败,err:", err)
		return nil
	}
	return bytes
}

func CreateMessages(messages []*Message) (err error) {
	result := db.DB.Create(messages)
	if err = result.Error; err != nil {
		return err
	}
	return nil
}
