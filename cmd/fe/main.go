package main

import (
	"github.com/ncuhome/GeniusAuthoritarian/internal/router"
	log "github.com/sirupsen/logrus"
)

func main() {
	log.Infoln("Sys Boost")
	if e := router.Engine().Run(":80"); e != nil {
		log.Fatalln("启动监听失败:", e)
	}
}
