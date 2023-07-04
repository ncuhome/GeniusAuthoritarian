//go:build gate

package router

import (
	"github.com/gin-gonic/gin"
	"github.com/ncuhome/GeniusAuthoritarian/internal/api/router/public"
)

func ApiV1(G *gin.RouterGroup) {
	public.Router(G.Group("public"))
}
