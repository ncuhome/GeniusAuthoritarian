//go:build dev

package global

import "github.com/Mmx233/config"

func initConfig() {
	// 调试模式从 yaml 载入配置
	c := config.NewConfig(&config.Options{
		Config:    &Config,
		Default:   &Config,
		Overwrite: true,
	})
	if e := c.Load(); e != nil {
		panic(e)
	}
}
