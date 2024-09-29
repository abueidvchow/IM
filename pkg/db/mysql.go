package db

import (
	"IM/config"
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB
var err error

func InitMySQL(cfg *config.MySQLConfig) error {
	dsl := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local", cfg.UserName,
		cfg.Password, cfg.Host, cfg.Port, cfg.DB)

	DB, err = gorm.Open(mysql.Open(dsl), &gorm.Config{
		QueryFields: true,
	})
	if err != nil {
		return err
	}
	db, err := DB.DB()
	if err != nil {
		return err
	}

	db.SetMaxIdleConns(cfg.MaxIdleConns) //设置最大闲置连接
	db.SetMaxOpenConns(cfg.MaxOpenConns) //设置最大连接数

	return nil
}
