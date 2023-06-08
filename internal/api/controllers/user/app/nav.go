package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/ncuhome/GeniusAuthoritarian/internal/api/callback"
	"github.com/ncuhome/GeniusAuthoritarian/internal/pkg/jwt"
	"github.com/ncuhome/GeniusAuthoritarian/internal/service"
	"github.com/ncuhome/GeniusAuthoritarian/internal/tools"
	"gorm.io/gorm"
)

func LandingApp(c *gin.Context) {
	var f struct {
		ID   uint   `json:"id" form:"id" binding:"required"`
		Code string `json:"code" form:"code"`
	}
	if e := c.ShouldBind(&f); e != nil {
		callback.Error(c, callback.ErrForm, e)
		return
	}

	user := tools.GetUserInfo(c)

	mfaSecret, e := service.User.FindMfa(user.ID)
	if e != nil {
		callback.Error(c, callback.ErrDBOperation, e)
		return
	} else if mfaSecret != "" {
		if f.Code == "" {
			callback.Error(c, callback.ErrMfaRequired)
			return
		} else if len(f.Code) != 6 {
			callback.Error(c, callback.ErrMfaCode)
			return
		}

		var valid bool
		valid, e = tools.VerifyMfa(f.Code, mfaSecret)
		if e != nil {
			callback.Error(c, callback.ErrUnexpected, e)
			return
		} else if !valid {
			callback.Error(c, callback.ErrMfaCode)
			return
		}
	}

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

	token, e := jwt.GenerateLoginToken(user.ID, app.ID, user.Name, c.ClientIP(), user.Groups)
	if e != nil {
		callback.Error(c, callback.ErrUnexpected, e)
		return
	}

	callbackUrl, e := tools.GenCallback(app.Callback, token)
	if e != nil {
		callback.Error(c, callback.ErrUnexpected, e)
		return
	}

	callback.Success(c, gin.H{
		"url": callbackUrl,
	})
}
