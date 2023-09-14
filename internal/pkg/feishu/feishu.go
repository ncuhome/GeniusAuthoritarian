package feishu

import (
	"github.com/ncuhome/GeniusAuthoritarian/internal/db/dao"
	"github.com/ncuhome/GeniusAuthoritarian/internal/global"
	"github.com/ncuhome/GeniusAuthoritarian/internal/tools"
	"github.com/ncuhome/GeniusAuthoritarian/pkg/feishuApi"
	log "github.com/sirupsen/logrus"
)

var Api = feishuApi.New(global.Config.Feishu.ClientID, global.Config.Feishu.Secret, tools.Http.Client)

func InitSync() {
	var count int64
	err := dao.DB.Model(&dao.FeishuGroups{}).Count(&count).Error
	if err != nil {
		log.Fatalln(err)
	}
	if count == 0 {
		if err = DepartmentSync(); err != nil {
			log.Fatalf("同步飞书部门失败: %v", err)
		}
	}
	if err = dao.DB.Model(&dao.User{}).Count(&count).Error; err != nil {
		log.Fatalln(err)
	}
	if count == 0 {
		var sync = UserSyncProcessor{}
		if err = sync.Run(); err != nil {
			log.Fatalf("同步飞书用户失败: %v", err)
		} else {
			log.Infoln("飞书用户列表已同步")
			sync.PrintSyncResult()
		}
	}

	if err = AddDepartmentSyncCron("0 5 * * *"); err != nil {
		log.Fatalf("添加定时同步飞书部门任务失败: %v", err)
	}
	if err = AddUserSyncCron("30 5 * * *"); err != nil {
		log.Fatalf("添加定时同步飞书用户任务失败: %v", err)
	}
}
