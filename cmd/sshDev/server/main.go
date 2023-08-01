package main

import (
	"github.com/ncuhome/GeniusAuthoritarian/internal/global"
	"github.com/ncuhome/GeniusAuthoritarian/internal/pkg/sshDev/rpc"
	log "github.com/sirupsen/logrus"
)

func main() {
	log.Infoln("Sys Boost")

	if err := rpc.Run(global.Config.SshDev.Token); err != nil {
		log.Fatalln("启动 rpc 服务失败:", err)
	}
}
