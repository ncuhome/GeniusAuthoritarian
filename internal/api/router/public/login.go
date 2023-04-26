package public

import (
	"github.com/gin-gonic/gin"
	controllers "github.com/ncuhome/GeniusAuthoritarian/internal/api/controllers/login"
	dingTalkPkg "github.com/ncuhome/GeniusAuthoritarian/internal/pkg/dingTalk"
	feishuPkg "github.com/ncuhome/GeniusAuthoritarian/internal/pkg/feishu"
)

func routerLogin(G *gin.RouterGroup) {
	G.POST("/verify", controllers.VerifyToken)

	feishu := G.Group("feishu")
	feishu.GET("link", controllers.GetLoginLink(feishuPkg.Api.LoginLink))
	feishu.POST("/", controllers.FeishuLogin)

	dingTalk := G.Group("dingTalk")
	dingTalk.GET("link", controllers.GetLoginLink(dingTalkPkg.Api.LoginLink))
}
