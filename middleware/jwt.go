package middleware

import (
	"IM/common"
	"IM/controller"
	"IM/pkg/jwt"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"strings"
)

// 基于JWT的认证中间件
func JWTAuthMiddleware(c *gin.Context) {
	// 客户端携带Token有三种方式 1.放在请求头 2.放在请求体 3.放在URI
	// 这里假设Token放在Header的Authorization中，并使用Bearer开头
	// 这里的具体实现方式要依据你的实际业务情况决定
	authHeader := c.Request.Header.Get("Authorization")
	// 1.先判断请求头是否为空
	if authHeader == "" {
		controller.ResponseError(c, common.ERROR_NEED_LOGIN)
		zap.L().Error("需要登录")
		c.Abort()
		return
	}
	// 2.不为空的话，判断请求头是否符合格式
	parts := strings.SplitN(authHeader, " ", 2)
	if !(len(parts) == 2 && parts[0] == "Bearer") {
		controller.ResponseError(c, common.ERROR_INVALID_TOKEN)
		//中间件执行失败后，必须要执行Abort()和return
		c.Abort()
		return
	}
	// 3.符合JWT格式的话就执行解析
	mc, err := jwt.ParseToken(parts[1])
	if err != nil {
		controller.ResponseError(c, common.ERROR_INVALID_TOKEN)
		c.Abort()
		return
	}
	c.Set(controller.CtxUserIdKey, mc.UserId)

	c.Next()
}
