package service

import (
	"IM/model"
	"fmt"
	"go.uber.org/zap"
	"strconv"
	"time"
)

func CreateGroupService(uid int64, name string, idsStr []string) (err error) {
	ids := make([]int64, len(idsStr)+1)
	ids[0] = uid
	fmt.Println("idsStr:", idsStr)
	for i, idStr := range idsStr {
		ids[i+1], err = strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			zap.L().Error("CreateGroupServiceError:", zap.Error(err))
			return err
		}
	}

	//创建群
	group := &model.Group{
		Name:       name,
		OwnerID:    uid,
		CreateTime: time.Now(),
		UpdateTime: time.Now(),
	}

	return model.CreateGroup(group, ids)

}
