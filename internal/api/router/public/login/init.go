package login

import (
	"github.com/gin-gonic/gin"
	controllers "github.com/ncuhome/GeniusAuthoritarian/internal/api/controllers/public/login"
	"github.com/ncuhome/GeniusAuthoritarian/internal/api/middlewares"
)

func Router(G *gin.RouterGroup) {
	G.POST("/", controllers.Login) // 个人页面登录
	G.POST("/verify", middlewares.RequireAppSignature, controllers.VerifyToken)
	G.POST("/mfa", controllers.VerifyMfa)
	G.POST("/refresh", controllers.RefreshToken)
	G.POST("/access", controllers.VerifyAccessToken)

	routerLoginPasskey(G.Group("passkey"))

	thirdParty := G.Group(":app")
	thirdParty.POST("/", controllers.ThirdPartySelfLogin) // 登录鉴权控制系统
	thirdParty.POST("/:appCode", controllers.ThirdPartyLogin)

	thirdPartyLink := thirdParty.Group("link")
	thirdPartyLink.GET("/", controllers.GetSelfLoginLink) // 登录鉴权控制系统
	thirdPartyLink.GET("/:appCode", controllers.GetLoginLink)
}
