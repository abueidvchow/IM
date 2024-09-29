package service

import (
	"IM/common"
	"IM/common/web/request"
	"IM/model"
	"IM/pkg/jwt"
	"IM/pkg/snowflake"
	"crypto/md5"
	"fmt"
	"strconv"
)

func UserRegisterService(p *request.RegisterParam) (common.HttpStatusCode, error) {
	//检验username是否存在
	exist, err := model.CheckUserExist(p.Username)
	if err != nil {
		return common.ERROR_MYSQL, err
	}
	if exist {
		return common.ERROR_USER_EXIST, nil

	}
	// 写入数据库
	hash := md5.Sum([]byte(p.Password))
	md5String := fmt.Sprintf("%x", hash)
	user := &model.User{
		UserID:   snowflake.GenID(),
		Username: p.Username,
		Nickname: p.Nickname,
		Password: md5String,
	}
	n, err := model.RegisterUser(user)
	if err != nil {
		return common.ERROR_MYSQL, err

	} else if n == 0 {
		return common.ERROR_INVALID_PARAMS, nil

	}
	return common.SUCCESS_REGISTER, nil
}

func UserLoginService(p *request.LoginParam) (code common.HttpStatusCode, token, user_id string, err error) {
	//检验username是否存在
	exist, err := model.CheckUserExist(p.Username)
	if err != nil {
		return common.ERROR_MYSQL, "", "", err
	}
	if !exist {
		return common.ERROR_USER_NOT_EXIST, "", "", nil
	}
	user, err := model.LoginUser(p)
	if err != nil {
		return common.ERROR_MYSQL, "", "", err
	} else if user == nil {
		return common.ERROR_INVALID_PARAMS, "", "", nil
	}

	//发放token
	aToken, _, err := jwt.GentToken(user.UserID, user.Username)
	if err != nil {
		return common.ERROR_GENERATE_JWT, "", "", err
	}

	return common.SUCCESS_LOGIN, aToken, strconv.FormatInt(user.UserID, 10), nil

}
