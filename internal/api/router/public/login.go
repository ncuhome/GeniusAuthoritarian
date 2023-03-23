package public

import "github.com/gin-gonic/gin"

func routerLogin(G *gin.RouterGroup) {
	feishu := G.Group("feishu")
	feishu.POST("/")
}
