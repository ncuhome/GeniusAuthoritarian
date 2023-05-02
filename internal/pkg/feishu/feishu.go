package feishu

import (
	"github.com/ncuhome/GeniusAuthoritarian/internal/db/dao"
	"github.com/ncuhome/GeniusAuthoritarian/internal/global"
	"github.com/ncuhome/GeniusAuthoritarian/pkg/feishuApi"
	"github.com/ncuhome/GeniusAuthoritarian/tools"
	log "github.com/sirupsen/logrus"
)

func init() {
	var count int64
	e := dao.DB.Model(&dao.FeishuGroups{}).Count(&count).Error
	if e != nil {
		log.Fatalln(e)
	}
	if count == 0 {
		if e = DepartmentSync(); e != nil {
			log.Fatalf("同步飞书部门失败: %v", e)
		}
	}
	if e = dao.DB.Model(&dao.User{}).Count(&count).Error; e != nil {
		log.Fatalln(e)
	}
	if count == 0 {
		var sync = UserSyncProcessor{}
		if e = sync.Run(); e != nil {
			log.Fatalf("同步飞书用户失败: %v", e)
		} else {
			log.Infoln("飞书用户列表已同步")
			sync.PrintSyncResult()
		}
	}

	if e = AddDepartmentSyncCron("0 5 * * *"); e != nil {
		log.Fatalf("添加定时同步飞书部门任务失败: %v", e)
	}
	if e = AddUserSyncCron("30 5 * * *"); e != nil {
		log.Fatalf("添加定时同步飞书用户任务失败: %v", e)
	}
}

var Api = feishuApi.New(global.Config.Feishu.ClientID, global.Config.Feishu.Secret, tools.Http.Client)
