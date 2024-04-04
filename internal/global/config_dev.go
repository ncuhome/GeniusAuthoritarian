//go:build dev

package global

import (
	"github.com/Mmx233/config"
	"github.com/Mmx233/tool"
	log "github.com/sirupsen/logrus"
)

func initConfig() {
	// load config from yaml
	c := config.NewConfig(&config.Options{
		Config:    &Config,
		Default:   &Config,
		Overwrite: true,
	})
	if err := c.Load(); err != nil {
		panic(err)
	}

	// generate keys
	exist, err := tool.File.Exists(ConfigDir)
	if err != nil {
		panic(err)
	} else if !exist {
		initKeypair()
		log.Infoln("keypair generated")
	}
}
