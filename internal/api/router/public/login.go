package public

import (
	"github.com/gin-gonic/gin"
	controllers "github.com/ncuhome/GeniusAuthoritarian/internal/api/controllers/public/login"
)

func routerLogin(G *gin.RouterGroup) {
	G.POST("/", controllers.Login) // 个人页面登录
	G.POST("/verify", controllers.VerifyToken)

	thirdParty := G.Group(":app")
	thirdParty.POST("/")
	thirdParty.POST("/:appCode", controllers.ThirdPartyLogin)

	thirdPartyLink := thirdParty.Group("link")
	thirdPartyLink.GET("/", controllers.GetSelfLoginLink)
	thirdPartyLink.GET("/:appCode", controllers.GetLoginLink)
}
