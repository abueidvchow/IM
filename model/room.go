package model

import "time"

type Room struct {
	RoomID     int64     `json:"room_id" gorm:"column:room_id"`           //聊天室ID
	RoomName   string    `json:"room_name" gorm:"column:room_name"`       //聊天室名字
	RoomInfo   string    `json:"room_info" gorm:"column:room_info"`       //聊天室简介
	RoomUserID int64     `json:"room_user_id" gorm:"column:room_user_id"` //聊天室创建者
	CreateTime time.Time `json:"create_time" gorm:"column:create_time"`   //创建时间
	UpdateTime time.Time `json:"update_time" gorm:"column:update_time"`   //更新时间
}

func (Room) TableName() string {
	return "room"
}
