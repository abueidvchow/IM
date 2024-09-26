package mysql

import (
	"IM/model"
	"fmt"
	"go.uber.org/zap"
	"strconv"
)

func CheckUserFromUserRoom(uid string, rid string) bool {
	userID, _ := strconv.ParseInt(uid, 10, 64)
	roomID, _ := strconv.ParseInt(rid, 10, 64)
	userRoom := &model.UserRoom{}
	db := DB.Where("user_id = ? AND room_id = ?", userID, roomID).Find(&userRoom)
	if db.Error != nil {
		zap.L().Error(db.Error.Error())
		return false
	} else if db.RowsAffected == 0 {
		return false
	} else {
		return true
	}
}

func GetUserFromUserRoom(rid int64) ([]*model.UserRoom, error) {
	userRoom := []*model.UserRoom{}
	db := DB.Where("room_id = ?", rid).Find(&userRoom)
	fmt.Println("userRoom", userRoom)
	if db.Error != nil {
		return nil, db.Error
	} else if db.RowsAffected == 0 {
		return nil, nil
	} else {
		return userRoom, nil
	}
}
