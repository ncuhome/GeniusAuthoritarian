package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/ncuhome/GeniusAuthoritarian/internal/api/callback"
	"github.com/ncuhome/GeniusAuthoritarian/internal/service"
	"github.com/ncuhome/GeniusAuthoritarian/internal/tools"
)

func ListOwnedApp(c *gin.Context) {
	uid := tools.GetUserInfo(c).ID

	apps, err := service.App.GetUserOwnedApp(uid)
	if err != nil {
		callback.Error(c, callback.ErrDBOperation, err)
		return
	}

	callback.Success(c, apps)
}

func ListAccessibleApp(c *gin.Context) {
	uid := tools.GetUserInfo(c).ID

	permitAllApps, err := service.App.GetPermitAll()
	if err != nil {
		callback.Error(c, callback.ErrDBOperation, err)
		return
	}

	isCenterMember, err := service.UserGroups.IsCenterMember(uid)
	if err != nil {
		callback.Error(c, callback.ErrDBOperation, err)
		return
	}

	accessibleApps, err := service.App.GetUserAccessible(uid, isCenterMember)
	if err != nil {
		callback.Error(c, callback.ErrDBOperation, err)
		return
	}

	callback.Success(c, gin.H{
		"permitAll":  permitAllApps,
		"accessible": accessibleApps,
	})
}
