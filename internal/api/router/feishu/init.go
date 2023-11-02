package feishu

import (
	"github.com/gin-gonic/gin"
	controllers "github.com/ncuhome/GeniusAuthoritarian/internal/api/controllers/feishu"
)

func Router(G *gin.RouterGroup) {
	G.POST("webhook", controllers.Webhook)
}
