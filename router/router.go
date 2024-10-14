package router

import (
	"IM/controller"
	"IM/pkg/logger"
	"IM/pkg/middleware"
	"github.com/gin-gonic/gin"
)

func SetUpRouter() (r *gin.Engine) {
	r = gin.New()
	r.Use(logger.GinLogger(), logger.GinRecovery(true))

	r.GET("/", func(c *gin.Context) {
		c.String(200, "你好,世界")
	})

	r.POST("/register", controller.UserRegister)
	r.POST("/login", controller.UserLogin)

	auth := r.Group("/api", middleware.JWTAuthMiddleware)
	{

		auth.POST("/friend/add", controller.AddFriend)
		auth.POST("/group/create", controller.CreateGroup)
		auth.GET("/group/userList", controller.GroupUserList)

		auth.GET("/ws", controller.Chat)
	}

	return r
}
