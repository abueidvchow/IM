package model

import (
	"IM/pkg/db"
	"go.uber.org/zap"
	"strconv"
)

// 聊天室对应的用户
type UserRoom struct {
	ID        int64 `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	RoomID    int64 `json:"room_id" gorm:"room_id"`
	UserID    int64 `json:"user_id" gorm:"user_id"`
	MessageID int64 `json:"message_id" gorm:"message_id"`
	//CreateTime time.Time `json:"create_time" gorm:"create_time"`
	//UpdateTime time.Time `json:"update_time" gorm:"update_time"`
}

func (UserRoom) TableName() string {
	return "user_room"
}

func CheckUserFromUserRoom(uid int64, rid string) bool {
	roomID, _ := strconv.ParseInt(rid, 10, 64)
	userRoom := UserRoom{}
	result := db.DB.Where("user_id = ? AND room_id = ?", uid, roomID).Find(&userRoom)
	if result.Error != nil {
		zap.L().Error(result.Error.Error())
		return false
	} else if result.RowsAffected == 0 {
		return false
	} else {
		return true
	}
}

func GetUserFromUserRoom(rid int64) ([]*UserRoom, error) {
	userRoom := []*UserRoom{}
	result := db.DB.Where("room_id = ?", rid).Find(&userRoom)
	if result.Error != nil {
		return nil, result.Error
	} else if result.RowsAffected == 0 {
		return nil, nil
	} else {
		return userRoom, nil
	}
}
