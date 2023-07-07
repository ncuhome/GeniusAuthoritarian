package main

import (
	"github.com/ncuhome/GeniusAuthoritarian/internal/pkg/GroupOperator"
	"github.com/ncuhome/GeniusAuthoritarian/internal/pkg/agent"
	"github.com/ncuhome/GeniusAuthoritarian/internal/pkg/feishu"
	"github.com/ncuhome/GeniusAuthoritarian/internal/pkg/views"
	"github.com/ncuhome/GeniusAuthoritarian/internal/router"
	"github.com/ncuhome/GeniusAuthoritarian/internal/tools"
	log "github.com/sirupsen/logrus"
)

func init() {
	agent.Init()
	feishu.InitSync()
	views.InitRenewAgent()
	GroupOperator.InitGroupRelation()
}

// 主程序，包含所有路由，不可多实例运行
func main() {
	log.Infoln("Sys Boost")
	if e := tools.SoftHttpSrv(router.Engine()); e != nil {
		log.Fatalln("启动监听失败:", e)
	}
}
