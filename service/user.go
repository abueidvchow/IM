package service

import (
	"IM/common"
	"IM/common/web/request"
	"IM/config"
	"IM/model"
	"IM/pkg/db"
	"IM/pkg/jwt"
	"IM/pkg/snowflake"
	"IM/service/ws"
	"crypto/md5"
	"errors"
	"fmt"
	"strconv"
)

func UserRegisterService(p *request.RegisterParam) (common.HttpStatusCode, error) {
	//检验username是否存在
	exist, err := model.CheckUserNameExist(p.Username)
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
	exist, err := model.CheckUserNameExist(p.Username)
	if err != nil {
		return common.ERROR_MYSQL, "", "", err
	}
	if !exist {
		return common.ERROR_USER_NOT_EXIST, "", "", nil
	}

	// 检验用户账号密码并返回该用户信息
	user, err := model.LoginUser(p)
	if err != nil {
		return common.ERROR_MYSQL, "", "", err
	} else if user == nil {
		return common.ERROR_INVALID_PARAMS, "", "", nil
	}

	// 检查用户是否已经在其他连接登录
	onlineAddr, err := db.GetUserOnline(user.UserID)
	if err != nil {
		return common.ERROR_REDIS, "", "", err
	}
	if onlineAddr != "" {
		fmt.Println("[用户登录] 用户已经在其他连接登录")
		return common.ERROR_USER_LOGINED, "", "", errors.New(common.HttpMsg[common.ERROR_USER_LOGINED])
	}

	//发放token
	aToken, _, err := jwt.GentToken(user.UserID, user.Username)
	if err != nil {
		return common.ERROR_GENERATE_JWT, "", "", err
	}

	// 加入Redis在线列表
	rpcAddr := config.Conf.IP + ":" + strconv.Itoa(config.Conf.RPCPort)
	err = db.SetUserOnline(user.UserID, rpcAddr)
	if err != nil {
		return common.ERROR_REDIS, "", "", err
	}
	// 加入用户在线列表
	ws.OnlineUser = append(ws.OnlineUser, user.UserID)

	return common.SUCCESS_LOGIN, aToken, strconv.FormatInt(user.UserID, 10), nil

}
