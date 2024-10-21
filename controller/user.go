package controller

import (
	"IM/common"
	"IM/common/web/request"
	"IM/service"
	"github.com/gin-gonic/gin"
)

const CtxUserIdKey = "user_id"

func UserRegister(c *gin.Context) {
	// 检验参数
	p := &request.RegisterParam{}
	if err := c.ShouldBind(p); err != nil {
		ResponseError(c, common.ERROR_INVALID_PARAMS, err.Error())
		return
	}

	code, err := service.UserRegisterService(p)
	if err != nil {
		ResponseError(c, code, err.Error())
		return
	}
	if code < 2000 {
		ResponseSuccess(c, code, p)
	} else {
		ResponseError(c, code)
	}

}

func UserLogin(c *gin.Context) {
	//检验参数
	p := &request.LoginParam{}
	if err := c.ShouldBind(p); err != nil {
		ResponseError(c, common.ERROR_INVALID_PARAMS, err.Error())
		return
	}
	code, token, user_id, err := service.UserLoginService(p)
	if err != nil {
		ResponseError(c, code, err.Error())
		return
	}
	if code < 2000 {
		c.Set(CtxUserIdKey, user_id)
		type Data struct {
			Token  string `json:"token"`
			UserId string `json:"user_id"`
		}
		data := Data{Token: token, UserId: user_id}
		ResponseSuccess(c, code, data)
	} else {
		ResponseError(c, code)
	}
}
