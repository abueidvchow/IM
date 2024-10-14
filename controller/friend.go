package controller

import (
	"IM/common"
	"IM/model"
	"IM/service"
	"github.com/gin-gonic/gin"
	"strconv"
)

func AddFriend(c *gin.Context) {
	//输入用户id添加
	/*POSTMAN 添加好友格式
	post body form-data
	{"friend_id":2658626632155136}
	*/
	uid, _ := c.Get(CtxUserIdKey)

	var userID int64 = uid.(int64)
	friendIDStr := c.PostForm("friend_id")
	friendID, _ := strconv.ParseInt(friendIDStr, 10, 64)
	if friendID == userID {
		ResponseError(c, common.ERROR_INVALID_PARAMS, "不能添加自己为好友")
		return
	}
	exist, err := model.CheckUserIDExist(friendID)
	if err != nil {
		ResponseError(c, common.ERROR_MYSQL, "查询用户ID时出错")
		return
	}
	if !exist {
		ResponseError(c, common.ERROR_USER_NOT_EXIST)
		return
	}
	err = service.AddFriend(userID, friendID)
	if err != nil {
		ResponseError(c, common.ERROR_ADD_FRIEND, err)
		return
	}

	ResponseSuccess(c, common.SUCCESS_ADD_FRIEND)

}
