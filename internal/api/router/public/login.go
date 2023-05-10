package public

import (
	"github.com/gin-gonic/gin"
	controllers "github.com/ncuhome/GeniusAuthoritarian/internal/api/controllers/public/login"
)

func routerLogin(G *gin.RouterGroup) {
	G.POST("/", controllers.Login) // 个人页面登录
	G.POST("/verify", controllers.VerifyToken)

	feishu := G.Group("feishu")
	feishu.GET("link", controllers.FeishuLoginLink)
	feishu.POST("/", controllers.FeishuLogin)

	dingTalk := G.Group("dingTalk")
	dingTalk.GET("link", controllers.DingTalkLoginLink)
	dingTalk.POST("/", controllers.DingTalkLogin)
}
