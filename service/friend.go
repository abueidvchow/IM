package service

import (
	"IM/model"
	"errors"
)

func AddFriend(uid, fid int64) (err error) {
	//检查是否是好友
	flag, err := model.IsFriend(uid, fid)
	if err != nil {
		return err
	}
	if flag {
		return errors.New("不能重复添加为好友")
	}

	//添加好友
	err = model.UserAddFriend(uid, fid)
	if err != nil {
		return err
	}
	return
}
