package main

import (
	"github.com/ncuhome/GeniusAuthoritarian/internal/router"
	log "github.com/sirupsen/logrus"
)

func main() {
	log.Infoln("Sys Boost")
	if err := router.Engine().Run(":80"); err != nil {
		log.Fatalln("启动监听失败:", err)
	}
}
