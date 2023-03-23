package public

import (
	"github.com/gin-gonic/gin"
	controllers "github.com/ncuhome/GeniusAuthoritarian/internal/api/controllers/login"
)

func routerLogin(G *gin.RouterGroup) {
	feishu := G.Group("feishu")
	feishu.GET("/", controllers.GoFeishuLogin)
	feishu.POST("/")
}
