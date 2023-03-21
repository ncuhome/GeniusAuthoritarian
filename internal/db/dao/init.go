package dao

import (
	"github.com/ncuhome/GeniusAuthoritarian/internal/global"
	"github.com/ncuhome/GeniusAuthoritarian/pkg/drivers/mysql"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

var DB *gorm.DB

func init() {
	var e error
	DB, e = mysql.New(&global.Config.Mysql)
	if e != nil {
		log.Fatalln("连接 Mysql 失败:", e)
	}

	if e = DB.AutoMigrate(
		&LoginRecord{},
	); e != nil {
		log.Fatalln("AutoMigration failed:", e)
	}
}
