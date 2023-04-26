package mysql

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"time"
)

func New(conf *Config, gormConfig *gorm.Config) (*gorm.DB, error) {
	//数据库初始化
	db, e := gorm.Open(mysql.Open(fmt.Sprintf(
		"%s:%s@tcp(%s:%d)/%s?%s",
		conf.Username,
		conf.Password,
		conf.Host,
		conf.Port,
		conf.Database,
		conf.Arg,
	)), gormConfig)
	if e != nil {
		return nil, e
	}

	//连接池设置
	if sqlDB, err := db.DB(); err == nil {
		sqlDB.SetConnMaxLifetime(time.Hour * 5)
	}

	return db, nil
}
