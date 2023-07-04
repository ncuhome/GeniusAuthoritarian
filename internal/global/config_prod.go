//go:build !dev

package global

import (
	"github.com/Mmx233/EnvConfig"
)

func initConfig() {
	EnvConfig.Load("", &Config)
}
