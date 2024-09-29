package model

import (
	"IM/pkg/db"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"time"
)

type Group struct {
	ID         int64     `gorm:"primary_key;auto_increment;comment:'自增主键'" json:"id"`
	Name       string    `gorm:"not null;comment:'群组名称'" json:"name"`
	OwnerID    int64     `gorm:"not null;comment:'群主id'" json:"owner_id"`
	CreateTime time.Time `gorm:"not null;default:CURRENT_TIMESTAMP;comment:'创建时间'" json:"create_time"`
	UpdateTime time.Time `gorm:"not null;default:CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP;comment:'更新时间'" json:"update_time"`
}

func (Group) TableName() string {
	return "group"
}

func CreateGroup(group *Group, ids []int64) (err error) {
	//批量插入数据可以开启事务
	return db.DB.Transaction(func(tx *gorm.DB) error {
		result := tx.Create(&group)
		if result.Error != nil {
			zap.L().Error("model.CreateGroup Error:", zap.Error(result.Error))
			return result.Error
		}
		var groupUsers []GroupUser
		for i := 0; i < len(ids); i++ {
			groupUsers = append(groupUsers, GroupUser{
				GroupID:    group.ID,
				UserID:     ids[i],
				CreateTime: time.Now(),
				UpdateTime: time.Now(),
			})
		}
		return tx.Create(&groupUsers).Error

	})

}
