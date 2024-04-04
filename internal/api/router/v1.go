package router

import (
	"github.com/gin-gonic/gin"
	"github.com/ncuhome/GeniusAuthoritarian/internal/api/router/app"
	"github.com/ncuhome/GeniusAuthoritarian/internal/api/router/feishu"
	"github.com/ncuhome/GeniusAuthoritarian/internal/api/router/public"
	"github.com/ncuhome/GeniusAuthoritarian/internal/api/router/user"
)

func ApiV1(G *gin.RouterGroup) {
	public.Router(G.Group("public"))
	app.Router(G.Group("app"))
	user.Router(G.Group("user"))
	feishu.Router(G.Group("feishu"))
}
