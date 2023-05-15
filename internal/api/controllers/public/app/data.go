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
		callback.Error(c, e, callback.ErrForm)
		return
	}

	var appName, appHost string
	if f.AppCode != "" {
		if ok, e := service.App.CheckAppCode(f.AppCode); e != nil {
			callback.Error(c, e, callback.ErrDBOperation)
			return
		} else if !ok {
			callback.Error(c, e, callback.ErrAppCodeNotFound)
			return
		}

		appInfo, e := service.App.FirstAppByAppCode(f.AppCode)
		if e != nil {
			callback.Error(c, e, callback.ErrDBOperation)
			return
		}

		callbackUrl, e := url.Parse(appInfo.Callback)
		if e != nil {
			callback.Error(c, e, callback.ErrUnexpected)
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
