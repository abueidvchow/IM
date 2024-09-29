package model

import (
	"IM/pkg/db"
	"go.uber.org/zap"
	"time"
)

type Friend struct {
	ID         int64     `gorm:"primary_key;auto_increment;comment:'自增主键'" json:"id"`
	UserID     int64     `gorm:"not null;comment:'用户id'" json:"user_id"`
	FriendID   int64     `gorm:"not null;comment:'好友id'" json:"friend_id"`
	CreateTime time.Time `gorm:"not null;default:CURRENT_TIMESTAMP;comment:'创建时间'" json:"create_time"`
	UpdateTime time.Time `gorm:"not null;default:CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP;comment:'更新时间'" json:"update_time"`
}

func (Friend) TableName() string {
	return "friend"
}

func IsFriend(uid, fid int64) (bool, error) {
	friend := &Friend{}
	result := db.DB.Where("user_id = ? AND friend_id = ?", uid, fid).Find(&friend)
	if result.Error != nil {
		zap.L().Error("IsFriend error:", zap.Error(result.Error))
		return false, result.Error
	}
	return result.RowsAffected > 0, nil
}

func UserAddFriend(uid, fid int64) error {
	friend := &Friend{
		FriendID:   fid,
		UserID:     uid,
		CreateTime: time.Now(),
		UpdateTime: time.Now(),
	}
	result := db.DB.Create(&friend)
	if result.Error != nil {
		zap.L().Error("UserAddFriend error:", zap.Error(result.Error))
		return result.Error
	}
	return nil
}
