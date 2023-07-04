package main

import (
	"github.com/ncuhome/GeniusAuthoritarian/internal/pkg/agent"
	"github.com/ncuhome/GeniusAuthoritarian/internal/pkg/feishu"
	"github.com/ncuhome/GeniusAuthoritarian/internal/router"
	log "github.com/sirupsen/logrus"
)

func init() {
	agent.Init()
	feishu.InitSync()
}

func main() {
	log.Infoln("Sys Boost")
	if e := router.CoreEngine().Run(":80"); e != nil {
		log.Fatalln("启动监听失败:", e)
	}
}
