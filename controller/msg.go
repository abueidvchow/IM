package controller

import (
	"IM/service"
	"github.com/gin-gonic/gin"
)

func SendMsg(c *gin.Context) {
	uid, _ := c.Get(CtxUserIdKey)

	var userID int64 = uid.(int64)
	service.SendMessageService(c, userID)
}

func ChatList(c *gin.Context) {
	uid, _ := c.Get(CtxUserIdKey)

	var userID int64 = uid.(int64)
	service.ChatListService(c, userID)
}
