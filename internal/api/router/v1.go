package router

import (
	"github.com/gin-gonic/gin"
	"github.com/ncuhome/GeniusAuthoritarian/internal/api/router/public"
	"github.com/ncuhome/GeniusAuthoritarian/internal/api/router/user"
)

func ApiV1(G *gin.RouterGroup) {
	public.Router(G.Group("public"))
	user.Router(G.Group("user"))
}
