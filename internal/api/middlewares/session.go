package middlewares

import (
	"github.com/Mmx233/tool"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/redis"
	"github.com/gin-gonic/gin"
	"github.com/ncuhome/GeniusAuthoritarian/internal/global"
	"math/rand"
	"time"
	"unsafe"
)

var sessionStore sessions.Store

func init() {
	var err error
	secretKey := tool.RandString(rand.NewSource(time.Now().UnixNano()), 20)
	secretBytes := unsafe.Slice(unsafe.StringData(secretKey), len(secretKey))
	sessionStore, err = redis.NewStoreWithDB(
		10, "tcp",
		global.Config.Redis.Addr, global.Config.Redis.Password,
		"1",
		secretBytes,
	)
	if err != nil {
		panic(err)
	}
}

func EnableSession(name string) gin.HandlerFunc {
	return sessions.Sessions(name, sessionStore)
}
