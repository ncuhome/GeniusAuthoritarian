package user

import (
	"github.com/gin-gonic/gin"
	"github.com/ncuhome/GeniusAuthoritarian/internal/api/controllers/app"
)

func Router(G *gin.RouterGroup) {
	G.GET("info", controllers.GetUserPublicInfo)
}
