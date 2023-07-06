package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/ncuhome/GeniusAuthoritarian/internal/api/callback"
	"github.com/ncuhome/GeniusAuthoritarian/internal/service"
	"github.com/ncuhome/GeniusAuthoritarian/internal/tools"
	"gorm.io/gorm"
	"net/url"
)

func LandingApp(c *gin.Context) {
	var f struct {
		ID uint `json:"id" form:"id" binding:"required"`
	}
	if e := c.ShouldBind(&f); e != nil {
		callback.Error(c, callback.ErrForm, e)
		return
	}

	user := tools.GetUserInfo(c)

	app, e := service.App.FirstAppByID(f.ID)
	if e != nil {
		if e == gorm.ErrRecordNotFound {
			callback.Error(c, callback.ErrAppNotFound)
			return
		}
		callback.Error(c, callback.ErrDBOperation, e)
		return
	}

	if !app.PermitAllGroup {
		var yes bool
		yes, e = service.App.UserAccessible(f.ID, user.ID)
		if e != nil {
			callback.Error(c, callback.ErrDBOperation, e)
			return
		} else if !yes {
			callback.ErrorWithTip(c, callback.ErrOperationIllegal, "没有访问该应用的权限")
			return
		}
	}

	callbackUrl, e := url.Parse(app.Callback)
	if e != nil {
		callback.Error(c, callback.ErrUnexpected, e)
		return
	}

	callbackUrl.Path = ""
	callbackUrl.RawQuery = ""

	callback.Success(c, gin.H{
		"url": callbackUrl.String(),
	})
}
