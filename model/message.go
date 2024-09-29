package model

import (
	"IM/pkg/db"
	"errors"
	"go.uber.org/zap"
	"time"
)

// 真正发送的消息结构体
type MessageStruct struct {
	UserID  int64  `json:"user_id" gorm:"column:user_id"` // 发送者ID
	RoomID  int64  `json:"room_id" gorm:"column:room_id"` // 房间ID
	Content string `json:"content" gorm:"column:content"` // 消息内容
}

// 数据库的
type Message struct {
	//MessageID int64 `json:"message_id" gorm:"column:message_id"` // 消息ID
	UserID int64 `json:"user_id" gorm:"column:user_id"` // 发送者ID
	RoomID int64 `json:"room_id" gorm:"column:room_id"` // 房间ID
	//TargetID    int64  `json:"target_id" gorm:"column:target_id"`       // 接收者ID
	//SendType    int    `json:"send_type" gorm:"column:send_type"`       // 发送类型 私发 群发 广播
	//MessageType int    `json:"message_type" gorm:"column:message_type"` // 文字 图片...
	Content    string    `json:"content" gorm:"column:content"` // 消息内容
	CreateTime time.Time `json:"create_time" gorm:"column:create_time"`
	UpdateTime time.Time `json:"update_time" gorm:"column:update_time"`
}

func (Message) TableName() string {
	return "message"
}

func SaveMessage(msg *Message) error {

	result := db.DB.Create(&msg)
	if result.Error != nil {
		zap.L().Error("消息保存出错：", zap.Error(result.Error))
		return result.Error
	}
	if result.RowsAffected == 0 {
		err := errors.New("消息未保存成功")
		zap.L().Error(err.Error())
		return err
	}
	return nil
}

func GetChatList(room_id string, page, size int) ([]Message, error) {
	var msgs []Message = make([]Message, 0)
	skip := (page - 1) * size
	result := db.DB.Offset(skip).Limit(skip+size).Find(&msgs, "room_id=?", room_id)
	if result.Error != nil {
		zap.L().Error(result.Error.Error())
		return nil, result.Error
	}
	return msgs, nil
}
