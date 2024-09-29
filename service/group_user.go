package service

import "IM/model"

func GroupUserListService(groupId int64) (groupUsers []model.GroupUser, err error) {
	groupUsers, err = model.GetGroupUserList(groupId)
	if err != nil {
		return
	}
	return groupUsers, nil
}
