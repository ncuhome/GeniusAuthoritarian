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

	var err error
	DB, err = mysql.New(&global.Config.Mysql, conf)
	if err != nil {
		log.Fatalln("连接 Mysql 失败:", err)
	}

	if err = DB.AutoMigrate(
		&User{},
		&SiteWhiteList{},
		&BaseGroup{},
		&App{},
		&LoginRecord{},
		&UserGroups{},
		&FeishuGroups{},
		&AppGroup{},
		&UserSsh{},
	); err != nil {
		log.Fatalln("AutoMigration failed:", err)
	}
}
