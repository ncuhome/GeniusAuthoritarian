package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/ncuhome/GeniusAuthoritarian/internal/api/callback"
	"github.com/ncuhome/GeniusAuthoritarian/internal/service"
	"github.com/ncuhome/GeniusAuthoritarian/internal/tools"
)

func ProfileData(c *gin.Context) {
	uid := tools.GetUserInfo(c).ID
	profile, e := service.User.UserProfile(uid)
	if e != nil {
		callback.Error(c, callback.ErrDBOperation, e)
		return
	}
	loginRecord, e := service.LoginRecord.UserHistory(uid, 10)
	if e != nil {
		callback.Error(c, callback.ErrDBOperation, e)
		return
	}
	callback.Success(c, gin.H{
		"user":        profile,
		"loginRecord": loginRecord,
	})
}
