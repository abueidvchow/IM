package controller

import (
	"IM/common"
	"IM/service"
	"github.com/gin-gonic/gin"
	"strconv"
)

func CreateGroup(c *gin.Context) {
	//参数校验
	name := c.PostForm("name")
	idsStr := c.PostFormArray("ids") // 群成员 id，不包括群创建者
	if name == "" || len(idsStr) == 0 {
		ResponseError(c, common.ERROR_INVALID_PARAMS)
		return
	}
	uid, _ := c.Get(CtxUserIdKey)
	err := service.CreateGroupService(uid.(int64), name, idsStr)
	if err != nil {
		ResponseError(c, common.ERROR_CREATE_GROUP, err)
		return
	}

	ResponseSuccess(c, common.SUCCESS_CREATE_GROUP)
}

func GroupUserList(c *gin.Context) {
	// 参数校验
	groupIdStr := c.Query("group_id")
	groupId, _ := strconv.ParseInt(groupIdStr, 10, 64)
	if groupId == 0 {
		ResponseError(c, common.ERROR_INVALID_PARAMS)
		return
	}
	groupUsers, err := service.GroupUserListService(groupId)
	if err != nil {
		ResponseError(c, common.ERROR_GET_GROUP_USERS, err)
		return
	}
	ResponseSuccess(c, common.SUCCESS_GET_GROUP_USERS, groupUsers)
}
