package public

import (
	"github.com/gin-gonic/gin"
	controllers "github.com/ncuhome/GeniusAuthoritarian/internal/api/controllers/login"
	"github.com/ncuhome/GeniusAuthoritarian/internal/api/middlewares"
)

func routerLogin(G *gin.RouterGroup) {
	G.Use(middlewares.SiteFilter)

	feishu := G.Group("feishu")
	feishu.GET("link", controllers.FeishuLoginLink)
	feishu.POST("/", controllers.FeishuLogin)
}
