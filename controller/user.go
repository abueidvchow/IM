package controller

import (
	"IM/common/web"
	"IM/common/web/request"
	"github.com/gin-gonic/gin"
)

type UserController struct {
}

func UserRegister(c *gin.Context) {
	//检验参数
	p := &request.LoginParam{}
	if err := c.ShouldBind(p); err != nil {
		ResponseError(c, web.ERROR_INVALID_PARAMS, err.Error())
		return
	}
	ResponseSuccess(c, web.SUCCESS_REGISTER, p)
}

func UserLogin(c *gin.Context) {
	//检验参数
}
