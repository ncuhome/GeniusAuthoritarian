package router

import (
	"github.com/gin-gonic/gin"
	"github.com/ncuhome/GeniusAuthoritarian/internal/api/middlewares"
	"github.com/ncuhome/GeniusAuthoritarian/internal/api/router/public"
)

func ApiV1(G *gin.RouterGroup) {
	G.Use(middlewares.SiteFilter)

	public.Router(G.Group("public"))
}
