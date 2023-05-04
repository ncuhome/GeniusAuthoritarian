package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/ncuhome/GeniusAuthoritarian/internal/api/callback"
	"github.com/ncuhome/GeniusAuthoritarian/internal/service"
	"github.com/ncuhome/GeniusAuthoritarian/tools"
)

func ProfileData(c *gin.Context) {
	profile, e := service.User.UserProfile(tools.GetUID(c))
	if e != nil {
		callback.Error(c, e, callback.ErrDBOperation)
		return
	}
	callback.Success(c, profile)
}

func SetAvatar(c *gin.Context) {
	file, e := c.FormFile("file")
	if e != nil {
		callback.Error(c, e, callback.ErrForm)
		return
	}

}
