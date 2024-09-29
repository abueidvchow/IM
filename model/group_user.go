package model

import (
	"IM/pkg/db"
	"go.uber.org/zap"
	"time"
)

type GroupUser struct {
	ID         int64     `gorm:"primary_key;auto_increment;comment:'自增主键'" json:"id"`
	GroupID    int64     `gorm:"not null;comment:'组id'" json:"group_id"`
	UserID     int64     `gorm:"not null;comment:'用户id'" json:"user_id"`
	CreateTime time.Time `gorm:"not null;default:CURRENT_TIMESTAMP;comment:'创建时间'" json:"create_time"`
	UpdateTime time.Time `gorm:"not null;default:CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP;comment:'更新时间'" json:"update_time"`
}

func (GroupUser) TableName() string {
	return "group_user"
}

func GetGroupUserList(groupID int64) (groupUsers []GroupUser, err error) {
	groupUsers = []GroupUser{}
	result := db.DB.Where("group_id=?", groupID).Find(&groupUsers)
	if result.Error != nil {
		zap.L().Error("model.GetGroupUserList Error:", zap.Error(result.Error))
		return nil, result.Error
	}
	return groupUsers, nil
}
