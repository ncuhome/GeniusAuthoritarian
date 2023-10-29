package user

import (
	"github.com/gin-gonic/gin"
	"github.com/ncuhome/GeniusAuthoritarian/internal/api/middlewares"
	"github.com/ncuhome/GeniusAuthoritarian/internal/api/router/user/admin"
	"github.com/ncuhome/GeniusAuthoritarian/internal/api/router/user/dev"
	"github.com/ncuhome/GeniusAuthoritarian/internal/api/router/user/passkey"
)

func Router(G *gin.RouterGroup) {
	G.Use(middlewares.UserAuth)

	routerProfile(G.Group("profile"))
	routerApp(G.Group("app"))
	routerGroups(G.Group("group"))
	routerMfa(G.Group("mfa"))
	routerIdentity(G.Group("identity"))
	routerU2F(G.Group("u2f"))

	passkey.Router(G.Group("passkey"))
	dev.Router(G.Group("dev"))
	admin.Router(G.Group("admin"))
}
