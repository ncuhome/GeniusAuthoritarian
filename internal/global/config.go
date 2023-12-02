package global

import (
	"github.com/ncuhome/GeniusAuthoritarian/internal/global/models"
	log "github.com/sirupsen/logrus"
)

var Config models.Config

const ThisAppName = "统一授权系统"

func checkConfig() {
	if Config.SshDev.Token == "" {
		log.Fatalln("请配置 Token")
	}
}
