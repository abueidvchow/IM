package model

import (
	"IM/common/web/request"
	"IM/pkg/db"
	"crypto/md5"
	"fmt"
	"time"
)

type User struct {
	ID         uint64    `gorm:"primary_key;auto_increment;comment:'自增主键'" json:"id"`
	UserID     int64     `json:"user_id" gorm:"column:user_id"`
	Username   string    `json:"username" gorm:"column:user_name"`
	Nickname   string    `gorm:"column:nick_name;not null;comment:'昵称'" json:"nickname"`
	Password   string    `json:"password" gorm:"column:password"`
	CreateTime time.Time `gorm:"not null;default:CURRENT_TIMESTAMP;comment:'创建时间'" json:"create_time"`
	UpdateTime time.Time `gorm:"not null;default:CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP;comment:'更新时间'" json:"update_time"`
}

func (User) TableName() string {
	return "user"
}

func CheckUserExist(username string) (bool, error) {
	var count int64
	db.DB.Model(User{}).Count(&count)
	if count > 0 {
		var user User
		result := db.DB.Where("user_name = ?", username).Find(&user)
		if result.Error != nil {
			return false, result.Error
		}
		//不存在
		if result.RowsAffected == 0 {
			return false, nil
		} else { //存在
			return true, nil
		}
	} else {
		return false, nil
	}

}

func RegisterUser(user *User) (n int, err error) {
	var count int64
	db.DB.Model(User{}).Count(&count)
	if count == 0 { //如果数据表为空
		result := db.DB.Create(&user)
		if result.Error != nil {
			return 0, result.Error
		}
		return int(result.RowsAffected), nil
	} else {
		result := db.DB.Create(&user)
		if result.Error != nil {
			return 0, result.Error
		}
		if result.RowsAffected == 0 {
			return 0, nil
		}
		return int(result.RowsAffected), nil
	}

}

func LoginUser(p *request.LoginParam) (user *User, err error) {
	//检验密码

	//加密密码
	data := []byte(p.Password)
	hash := md5.Sum(data)
	md5String := fmt.Sprintf("%x", hash)

	//到数据库里查询该用户名的密码
	user = &User{}
	result := db.DB.Where("user_name = ?", p.Username).Find(&user)
	if result.Error != nil {
		return nil, result.Error
	}
	//密码不对
	if md5String != user.Password {
		return nil, nil
	}

	//密码正确
	return user, nil
}
