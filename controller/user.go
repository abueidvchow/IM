package controller

import (
	"IM/common"
	"IM/common/web/request"
	"IM/service"
	"github.com/gin-gonic/gin"
)

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
	code, data, err := service.UserLoginService(p)
	if err != nil {
		ResponseError(c, code, err.Error())
		return
	}
	if code < 2000 {
		ResponseSuccess(c, code, data)
	} else {
		ResponseError(c, code)
	}
}
