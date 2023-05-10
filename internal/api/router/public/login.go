package public

import (
	"github.com/gin-gonic/gin"
	controllers "github.com/ncuhome/GeniusAuthoritarian/internal/api/controllers/public/login"
)

func routerLogin(G *gin.RouterGroup) {
	G.POST("/", controllers.Login) // 个人页面登录
	G.POST("/verify", controllers.VerifyToken)

	thirdParty := G.GET(":app")
	thirdParty.GET("/link", controllers.GetLoginLink)
	thirdParty.POST("/", controllers.ThirdPartyLogin)
}
