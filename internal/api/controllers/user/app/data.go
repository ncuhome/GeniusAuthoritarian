package app

import (
	"github.com/gin-gonic/gin"
	"github.com/ncuhome/GeniusAuthoritarian/internal/api/callback"
	"github.com/ncuhome/GeniusAuthoritarian/internal/service"
	"github.com/ncuhome/GeniusAuthoritarian/internal/tools"
)

func ListOwnedApp(c *gin.Context) {
	uid := tools.GetUserInfo(c).ID

	apps, e := service.App.GetUserOwnedApp(uid)
	if e != nil {
		callback.Error(c, e, callback.ErrDBOperation)
		return
	}

	callback.Success(c, apps)
}

func ListAccessAbleApp(c *gin.Context) {
	uid := tools.GetUserInfo(c).ID
}
