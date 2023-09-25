package main

import (
	"github.com/ncuhome/GeniusAuthoritarian/internal/global"
	"github.com/ncuhome/GeniusAuthoritarian/internal/pkg/cronAgent"
	"github.com/ncuhome/GeniusAuthoritarian/internal/pkg/department2BaseGroup"
	"github.com/ncuhome/GeniusAuthoritarian/internal/pkg/feishu"
	sshDevServer "github.com/ncuhome/GeniusAuthoritarian/internal/pkg/sshDev/server"
	"github.com/ncuhome/GeniusAuthoritarian/internal/pkg/sshDev/server/rpc"
	"github.com/ncuhome/GeniusAuthoritarian/internal/pkg/views"
	"github.com/ncuhome/GeniusAuthoritarian/internal/router"
	"github.com/ncuhome/GeniusAuthoritarian/internal/tools"
	log "github.com/sirupsen/logrus"
)

func init() {
	department2BaseGroup.Init()
	cron := cronAgent.New()
	feishu.InitSync(cron)
	views.InitRenewAgent(cron)
	// 建议放在用户同步的时间之后
	sshDevServer.AddSshAccountCron(cron, "0 6 * * *")
}

// 主程序，包含所有路由，不可多实例运行
func main() {
	log.Infoln("Sys Boost")

	go func() {
		if global.Config.SshDev.Token == "" {
			log.Fatalln("请配置 Token")
		}

		if err := rpc.Run(global.Config.SshDev.Token, ":81"); err != nil {
			log.Fatalln("启动 sshDev rpc 服务失败:", err)
		}
	}()

	// :80
	if err := tools.SoftHttpSrv(router.Engine()); err != nil {
		log.Fatalln("启动监听失败:", err)
	}
}
