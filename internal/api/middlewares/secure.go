package middlewares

import (
	"github.com/Mmx233/secure/v2"
	"github.com/Mmx233/secure/v2/drivers"
	"github.com/gin-gonic/gin"
	"github.com/ncuhome/GeniusAuthoritarian/internal/api/callback"
	"github.com/ncuhome/GeniusAuthoritarian/internal/db/redis"
	log "github.com/sirupsen/logrus"
)

func Secure() gin.HandlerFunc {
	middleware, e := secure.New(&secure.Config{
		Driver: &drivers.RedisDriver{
			Client: redis.Client,
		},
		HandleReachLimit: func(c *gin.Context) {
			callback.Error(c, callback.ErrRequestFrequency)
		},
		RateLimit: 60, // API 每分钟最大请求数
	})
	if e != nil {
		log.Fatalln(e)
	}
	return middleware.Handler
}
