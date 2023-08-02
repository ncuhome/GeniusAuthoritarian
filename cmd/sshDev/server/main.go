package main

import (
	"github.com/ncuhome/GeniusAuthoritarian/internal/global"
	"github.com/ncuhome/GeniusAuthoritarian/internal/pkg/agent"
	"github.com/ncuhome/GeniusAuthoritarian/internal/pkg/sshDev"
	"github.com/ncuhome/GeniusAuthoritarian/internal/pkg/sshDev/rpc"
	log "github.com/sirupsen/logrus"
)

func init() {
	agent.Init()
	// 建议放在用户同步的时间之后
	sshDev.AddCron("0 6 * * *")
}

func main() {
	log.Infoln("Sys Boost")

	if err := rpc.Run(global.Config.SshDev.Token); err != nil {
		log.Fatalln("启动 rpc 服务失败:", err)
	}
}
