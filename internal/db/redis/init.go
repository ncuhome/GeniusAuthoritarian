package redis

import (
	"github.com/go-redis/redis/v8"
	"github.com/ncuhome/GeniusAuthoritarian/internal/global"
	driver "github.com/ncuhome/GeniusAuthoritarian/pkg/drivers/redis"
	log "github.com/sirupsen/logrus"
)

var Client *redis.Client

func init() {
	var err error
	Client, err = driver.New(&global.Config.Redis)
	if err != nil {
		log.Fatalln("初始化 redis 失败:", err)
	}
}
