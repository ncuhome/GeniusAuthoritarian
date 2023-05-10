package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/ncuhome/GeniusAuthoritarian/internal/api/callback"
	"github.com/ncuhome/GeniusAuthoritarian/internal/db/dao"
	"github.com/ncuhome/GeniusAuthoritarian/internal/service"
)

func AppInfo(c *gin.Context) {
	var f struct {
		AppCode string `json:"appCode" form:"appCode"`
	}
	if e := c.ShouldBind(&f); e != nil {
		callback.Error(c, e, callback.ErrForm)
		return
	}

	var appInfo *dao.App
	if f.AppCode != "" {
		if ok, e := service.App.CheckAppCode(f.AppCode); e != nil {
			callback.Error(c, e, callback.ErrDBOperation)
			return
		} else if !ok {
			callback.Error(c, e, callback.ErrAppCodeNotFound)
			return
		}

		var e error
		appInfo, e = service.App.FistAppForLogin(f.AppCode)
		if e != nil {
			callback.Error(c, e, callback.ErrDBOperation)
			return
		}
	} else {
		appInfo = service.App.This(c.Request.Host)
	}

	callback.Success(c, gin.H{
		"name":      appInfo.Name,
		"createdAt": appInfo.CreatedAt,
	})
}
