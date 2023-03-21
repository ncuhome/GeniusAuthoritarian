package global

import (
	"github.com/Mmx233/EnvConfig"
	"github.com/ncuhome/GeniusAuthoritarian/internal/global/models"
	"os"
)

var Config models.Config

func initConfig() {
	EnvConfig.Load("", &Config)
}

var DevMode = os.Getenv("DEV_MODE") == "TRUE"
