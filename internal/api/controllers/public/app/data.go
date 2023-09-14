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
	if err := c.ShouldBind(&f); err != nil {
		callback.Error(c, callback.ErrForm, err)
		return
	}

	var appName, appHost string
	if f.AppCode != "" {
		if ok, err := service.App.CheckAppCode(f.AppCode); err != nil {
			callback.Error(c, callback.ErrDBOperation, err)
			return
		} else if !ok {
			callback.Error(c, callback.ErrAppCodeNotFound, err)
			return
		}

		appInfo, err := service.App.FirstAppByAppCode(f.AppCode)
		if err != nil {
			callback.Error(c, callback.ErrDBOperation, err)
			return
		}

		callbackUrl, err := url.Parse(appInfo.Callback)
		if err != nil {
			callback.Error(c, callback.ErrUnexpected, err)
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
