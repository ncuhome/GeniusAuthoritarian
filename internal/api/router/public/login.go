package public

import (
	"github.com/gin-gonic/gin"
	controllers "github.com/ncuhome/GeniusAuthoritarian/internal/api/controllers/login"
)

func routerLogin(G *gin.RouterGroup) {
	G.POST("/verify", controllers.VerifyToken)

	feishu := G.Group("feishu")
	feishu.GET("link", controllers.FeishuLoginLink)
	feishu.POST("/", controllers.FeishuLogin)
}
