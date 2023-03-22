package router

import (
	"github.com/gin-gonic/gin"
	"github.com/ncuhome/GeniusAuthoritarian/internal/api/middlewares"
)

func Api(G *gin.RouterGroup) {
	G.Use(middlewares.SiteFilter)
}
