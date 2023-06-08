package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/ncuhome/GeniusAuthoritarian/internal/api/callback"
	"github.com/ncuhome/GeniusAuthoritarian/internal/global"
	"github.com/ncuhome/GeniusAuthoritarian/internal/service"
	"net/url"
)

func AppInfo(c *gin.Context) {
	var f struct {
		AppCode string `json:"appCode" form:"appCode"`
	}
	if e := c.ShouldBind(&f); e != nil {
		callback.Error(c, callback.ErrForm, e)
		return
	}

	var appName, appHost string
	if f.AppCode != "" {
		if ok, e := service.App.CheckAppCode(f.AppCode); e != nil {
			callback.Error(c, callback.ErrDBOperation, e)
			return
		} else if !ok {
			callback.Error(c, callback.ErrAppCodeNotFound, e)
			return
		}

		appInfo, e := service.App.FirstAppByAppCode(f.AppCode)
		if e != nil {
			callback.Error(c, callback.ErrDBOperation, e)
			return
		}

		callbackUrl, e := url.Parse(appInfo.Callback)
		if e != nil {
			callback.Error(c, callback.ErrUnexpected, e)
			return
		}

		appName = appInfo.Name
		appHost = callbackUrl.Host
	} else {
		appName = global.ThisAppName
		appHost = c.Request.Host
	}

	callback.Success(c, gin.H{
		"name": appName,
		"host": appHost,
	})
}
