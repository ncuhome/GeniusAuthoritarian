package controllers

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
		callback.Error(c, callback.ErrDBOperation, e)
		return
	}

	callback.Success(c, apps)
}

func ListAccessibleApp(c *gin.Context) {
	uid := tools.GetUserInfo(c).ID

	permitAllApps, e := service.App.GetPermitAll()
	if e != nil {
		callback.Error(c, callback.ErrDBOperation, e)
		return
	}

	isCenterMember, e := service.UserGroups.IsCenterMember(uid)
	if e != nil {
		callback.Error(c, callback.ErrDBOperation, e)
		return
	}

	accessibleApps, e := service.App.GetUserAccessible(uid, isCenterMember)
	if e != nil {
		callback.Error(c, callback.ErrDBOperation, e)
		return
	}

	callback.Success(c, gin.H{
		"permitAll":  permitAllApps,
		"accessible": accessibleApps,
	})
}
