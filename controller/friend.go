package controller

import (
	"IM/common"
	"IM/service"
	"github.com/gin-gonic/gin"
	"strconv"
)

func AddFriend(c *gin.Context) {
	//输入用户id添加
	uid, _ := c.Get(CtxUserIdKey)

	var userID int64 = uid.(int64)
	friendIDStr := c.PostForm("friend_id")
	friendID, _ := strconv.ParseInt(friendIDStr, 10, 64)
	if friendID == userID {
		ResponseError(c, common.ERROR_INVALID_PARAMS, "不能添加自己为好友")
		return
	}
	err := service.AddFriend(userID, friendID)
	if err != nil {
		ResponseError(c, common.ERROR_ADD_FRIEND, err)
		return
	}

	ResponseSuccess(c, common.SUCCESS_ADD_FRIEND)

}
