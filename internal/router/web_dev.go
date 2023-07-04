//go:build dev && web

package router

import (
	gateway "github.com/Mmx233/Gateway/v2"
	"github.com/gin-gonic/gin"
	"github.com/ncuhome/GeniusAuthoritarian/internal/tools"
	log "github.com/sirupsen/logrus"
	"net/http"
)

func frontendHandler() gin.HandlerFunc {
	return gateway.Proxy(&gateway.ApiConf{
		Addr:      "localhost:5173",
		Transport: tools.Http.Client.Transport,
		ErrorHandler: func(_ http.ResponseWriter, _ *http.Request, err error) {
			log.Warnf("调试页面请求转发失败: %v", err)
		},
		AllowRequest: func(c *gin.Context) bool {
			return frontendRouterCheck(c)
		},
	})
}
