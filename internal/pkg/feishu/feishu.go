package feishu

import (
	"github.com/ncuhome/GeniusAuthoritarian/internal/db/redis"
	"github.com/ncuhome/GeniusAuthoritarian/internal/global"
	"github.com/ncuhome/GeniusAuthoritarian/internal/pkg/cronAgent"
	"github.com/ncuhome/GeniusAuthoritarian/internal/tools"
	"github.com/ncuhome/GeniusAuthoritarian/pkg/feishuApi"
	"github.com/robfig/cron/v3"
	log "github.com/sirupsen/logrus"
)

var Api = feishuApi.New(global.Config.Feishu.ClientID, global.Config.Feishu.Secret, tools.Http.Client)

func InitSync(c *cron.Cron) {
	deparmentSchedule, err := cronAgent.Parser.Parse("0 5 * * *")
	if err != nil {
		log.Fatalf("规划同步飞书部门定时任务失败: %v", err)
	}

	departmentBackoff := NewDepartmentSyncBackOff(
		redis.NewSyncStat("feishu-department"),
		deparmentSchedule,
	)
	if err = departmentBackoff.Content(); err != nil {
		log.Fatalln(err)
	}

	userSyncBackoff := NewUserSyncBackoff()
	if err = userSyncBackoff.Content(); err != nil {
		log.Fatalln(err)
	}

	c.Schedule(deparmentSchedule, cron.FuncJob(departmentBackoff.Start))
	if _, err = userSyncBackoff.AddCron(c, "30 5 * * *"); err != nil {
		log.Fatalf("添加定时同步飞书用户任务失败: %v", err)
	}
}
