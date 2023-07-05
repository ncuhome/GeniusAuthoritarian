package main

import (
	"github.com/ncuhome/GeniusAuthoritarian/internal/router"
	log "github.com/sirupsen/logrus"
)

// 登录逻辑子程序，请注意该程序允许多实例并行
func main() {
	log.Infoln("Sys Boost")
	if e := router.Engine().Run(":80"); e != nil {
		log.Fatalln("启动监听失败:", e)
	}
}
