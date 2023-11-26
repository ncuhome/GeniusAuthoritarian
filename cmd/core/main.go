package main

import (
	"github.com/ncuhome/GeniusAuthoritarian/internal/db/redis"
	"github.com/ncuhome/GeniusAuthoritarian/internal/global"
	"github.com/ncuhome/GeniusAuthoritarian/internal/pkg/cronAgent"
	"github.com/ncuhome/GeniusAuthoritarian/internal/pkg/department2BaseGroup"
	"github.com/ncuhome/GeniusAuthoritarian/internal/pkg/feishu"
	"github.com/ncuhome/GeniusAuthoritarian/internal/pkg/views"
	"github.com/ncuhome/GeniusAuthoritarian/internal/router"
	sshDevServer "github.com/ncuhome/GeniusAuthoritarian/internal/rpc/sshDev"
	"github.com/ncuhome/GeniusAuthoritarian/internal/rpc/sshDev/rpc"
	"github.com/ncuhome/GeniusAuthoritarian/internal/tools"
	log "github.com/sirupsen/logrus"
)

func init() {
	department2BaseGroup.Init(redis.NewSyncStat("init-base-groups"))

	cron := cronAgent.New()
	feishu.InitSync(cron)
	views.InitRenewAgent(cron, redis.NewSyncStat("renew-views"))
	// 建议放在用户同步的时间之后
	sshDevServer.AddSshAccountCron(cron, redis.NewSyncStat("dev-ssh"))

	cron.Start()
}

func sshDevRpc() {
	if global.Config.SshDev.Token == "" {
		log.Fatalln("请配置 Token")
	}

	if err := rpc.Run(global.Config.SshDev.Token, ":81"); err != nil {
		log.Fatalln("启动 sshDev rpc 服务失败:", err)
	}
}

func main() {
	log.Infoln("Sys Boost")

	go sshDevRpc()

	// :80
	if err := tools.SoftHttpSrv(router.Engine()); err != nil {
		log.Fatalln("启动监听失败:", err)
	}
}
