package mysql

import (
	"IM/model"
	"errors"
	"go.uber.org/zap"
)

func SaveMessage(msg *model.Message) error {

	db := DB.Create(&msg)
	if db.Error != nil {
		zap.L().Error("消息保存出错：", zap.Error(db.Error))
		return db.Error
	}
	if db.RowsAffected == 0 {
		err := errors.New("消息未保存成功")
		zap.L().Error(err.Error())
		return err
	}
	return nil
}
