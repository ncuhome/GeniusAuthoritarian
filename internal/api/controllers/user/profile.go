package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/ncuhome/GeniusAuthoritarian/internal/api/callback"
	"github.com/ncuhome/GeniusAuthoritarian/internal/service"
	"github.com/ncuhome/GeniusAuthoritarian/tools"
)

func ProfileData(c *gin.Context) {
	uid := tools.GetUID(c)
	profile, e := service.User.UserProfile(uid)
	if e != nil {
		callback.Error(c, e, callback.ErrDBOperation)
		return
	}
	loginRecord, e := service.LoginRecord.UserHistory(uid, 10)
	callback.Success(c, gin.H{
		"user":        profile,
		"loginRecord": loginRecord,
	})
}
