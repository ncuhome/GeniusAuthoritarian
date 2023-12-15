package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/ncuhome/GeniusAuthoritarian/internal/api/callback"
	"github.com/ncuhome/GeniusAuthoritarian/internal/service"
	"github.com/ncuhome/GeniusAuthoritarian/internal/tools"
)

func ProfileData(c *gin.Context) {
	claims := tools.GetUserInfo(c)
	uid := claims.UID
	profile, err := service.User.UserProfile(uid)
	if err != nil {
		callback.Error(c, callback.ErrDBOperation, err)
		return
	}
	loginRecord, err := service.LoginRecord.UserHistory(uid, 10)
	if err != nil {
		callback.Error(c, callback.ErrDBOperation, err)
		return
	}
	onlineLogin, err := service.LoginRecord.UserOnline(uid, uint(claims.ID))
	if err != nil {
		callback.Error(c, callback.ErrDBOperation, err)
		return
	}
	callback.Success(c, gin.H{
		"user": profile,
		"loginRecord": gin.H{
			"online":  onlineLogin,
			"history": loginRecord,
		},
	})
}
