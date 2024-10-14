package controller

import (
	"IM/service"
	"github.com/gin-gonic/gin"
)

func Chat(c *gin.Context) {
	// 检验是否有接收者ID
	//receiverIDStr := c.Query("receiver_id") // 接收此消息的ID，可以是用户ID也可以是群
	//receiverID, _ := strconv.ParseInt(receiverIDStr, 10, 64)
	//exist, err := model.CheckUserIDExist(receiverID)
	//if err != nil {
	//	ResponseError(c, common.ERROR_MYSQL, err)
	//	return
	//} else if !exist {
	//	ResponseError(c, common.ERROR_USER_NOT_EXIST, nil)
	//	return
	//}
	uid, _ := c.Get(CtxUserIdKey)
	//发送消息
	service.ChatService(uid.(int64), c.Writer, c.Request)
}
