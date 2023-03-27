package global

import (
	"github.com/Mmx233/EnvConfig"
	"github.com/Mmx233/config"
	"github.com/ncuhome/GeniusAuthoritarian/internal/global/models"
	"os"
)

var Config models.Config

var DevMode = os.Getenv("DEV_MODE") == "TRUE"

func initConfig() {
	if DevMode {
		// 调试模式从 yaml 载入配置
		c := config.NewConfig(&config.Options{
			Config:    &Config,
			Default:   &Config,
			Overwrite: true,
		})
		if e := c.Load(); e != nil {
			panic(e)
		}
	} else {
		EnvConfig.Load("", &Config)
	}
}
