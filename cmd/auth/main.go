package main

import (
	"github.com/ncuhome/GeniusAuthoritarian/internal/pkg/GroupOperator"
	"github.com/ncuhome/GeniusAuthoritarian/internal/router"
	"github.com/ncuhome/GeniusAuthoritarian/internal/tools"
	log "github.com/sirupsen/logrus"
)

func init() {
	GroupOperator.InitGroupRelation()
}

// 登录逻辑子程序，请注意该程序允许多实例并行
func main() {
	log.Infoln("Sys Boost")

	if e := tools.SoftHttpSrv(router.Engine()); e != nil {
		log.Fatalln("启动监听失败:", e)
	}
}
