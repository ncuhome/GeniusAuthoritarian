package feishu

import (
	"context"
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
	if err = departmentBackoff.Run(context.Background()); err != nil {
		log.Fatalln(err)
	}

	userSyncSchedule, err := cronAgent.Parser.Parse("30 5 * * *")
	if err != nil {
		log.Fatalf("规划同步飞书用户定时任务失败: %v", err)
	}

	userSyncBackoff := NewUserSyncBackoff(
		redis.NewSyncStat("feishu-user"),
		userSyncSchedule,
	)
	if err = userSyncBackoff.Run(context.Background()); err != nil {
		log.Fatalln(err)
	}

	c.Schedule(deparmentSchedule, cronAgent.FuncJobWithSingleton(departmentBackoff))
	c.Schedule(userSyncSchedule, cronAgent.FuncJobWithSingleton(userSyncBackoff))
}
