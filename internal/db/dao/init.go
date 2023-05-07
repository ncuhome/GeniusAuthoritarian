package dao

import (
	"github.com/ncuhome/GeniusAuthoritarian/internal/global"
	"github.com/ncuhome/GeniusAuthoritarian/pkg/drivers/mysql"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func init() {
	conf := &gorm.Config{
		SkipDefaultTransaction: true,
	}
	if global.Config.TraceMode {
		conf.Logger = logger.Default.LogMode(logger.Info)
	}

	var e error
	DB, e = mysql.New(&global.Config.Mysql, conf)
	if e != nil {
		log.Fatalln("连接 Mysql 失败:", e)
	}

	if e = DB.AutoMigrate(
		&User{},
		&SiteWhiteList{},
		&Group{},
		&LoginRecordWithForeignKey{},
		&UserGroupsWithForeignKey{},
		&FeishuGroupsWithForeignKey{},
		&App{},
	); e != nil {
		log.Fatalln("AutoMigration failed:", e)
	}

	if e = DB.AutoMigrate(
		&LoginRecord{},
		&UserGroups{},
		&FeishuGroups{},
	); e != nil {
		log.Fatalln("AutoMigration failed:", e)
	}
}
