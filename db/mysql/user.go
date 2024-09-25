package mysql

import (
	"IM/common/web/request"
	"IM/model"
	"crypto/md5"
	"fmt"
)

func CheckUserExist(username string) (bool, error) {

	var count int64
	DB.Model(&model.User{}).Count(&count)
	if count > 0 {
		var user model.User
		db := DB.Where("user_name = ?", username).Find(&user)
		if db.Error != nil {
			return false, db.Error
		}
		//不存在
		if db.RowsAffected == 0 {
			return false, nil
		} else { //存在
			return true, nil
		}
	} else {
		return false, nil
	}

}

func RegisterUser(user *model.User) (n int, err error) {
	var count int64
	DB.Model(&model.User{}).Count(&count)
	if count == 0 { //如果数据表为空
		db := DB.Create(&user)
		if db.Error != nil {
			return 0, db.Error
		}
		return int(db.RowsAffected), nil
	} else {
		db := DB.Create(&user)
		if db.Error != nil {
			return 0, db.Error
		}
		if db.RowsAffected == 0 {
			return 0, nil
		}
		return int(db.RowsAffected), nil
	}

}

func LoginUser(p *request.LoginParam) (user *model.User, err error) {
	//检验密码

	//加密密码
	data := []byte(p.Password)
	hash := md5.Sum(data)
	md5String := fmt.Sprintf("%x", hash)

	//到数据库里查询该用户名的密码
	user = &model.User{}
	db := DB.Where("user_name = ?", p.Username).Find(&user)
	if db.Error != nil {
		return nil, db.Error
	}
	//密码不对
	if md5String != user.Password {
		return nil, nil
	}

	//密码正确
	return user, nil
}
