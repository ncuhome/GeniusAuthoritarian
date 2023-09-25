package feishu

import (
	"github.com/ncuhome/GeniusAuthoritarian/internal/global"
	"github.com/ncuhome/GeniusAuthoritarian/internal/tools"
	"github.com/ncuhome/GeniusAuthoritarian/pkg/feishuApi"
	"github.com/robfig/cron/v3"
	log "github.com/sirupsen/logrus"
)

var Api = feishuApi.New(global.Config.Feishu.ClientID, global.Config.Feishu.Secret, tools.Http.Client)

func InitSync(c *cron.Cron) {
	departmentBackoff := NewDepartmentSyncBackOff()
	err := departmentBackoff.Content()
	if err != nil {
		log.Fatalln(err)
	}

	userSyncBackoff := NewUserSyncBackoff()
	if err = userSyncBackoff.Content(); err != nil {
		log.Fatalln(err)
	}

	if _, err = departmentBackoff.AddCron(c, "0 5 * * *"); err != nil {
		log.Fatalf("添加定时同步飞书部门任务失败: %v", err)
	}
	if _, err = userSyncBackoff.AddCron(c, "30 5 * * *"); err != nil {
		log.Fatalf("添加定时同步飞书用户任务失败: %v", err)
	}
}
