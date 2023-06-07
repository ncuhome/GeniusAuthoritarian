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
		callback.Error(c, e, callback.ErrForm)
		return
	}

	user := tools.GetUserInfo(c)

	mfaSecret, e := service.User.FindMfa(user.ID)
	if e != nil {
		callback.Error(c, e, callback.ErrDBOperation)
		return
	} else if mfaSecret != "" {
		if f.Code == "" {
			callback.Error(c, nil, callback.ErrMfaRequired)
			return
		} else if len(f.Code) != 6 {
			callback.Error(c, nil, callback.ErrMfaCode)
			return
		}

		var valid bool
		valid, e = tools.VerifyMfa(f.Code, mfaSecret)
		if e != nil {
			callback.Error(c, e, callback.ErrUnexpected)
			return
		} else if !valid {
			callback.Error(c, nil, callback.ErrMfaCode)
			return
		}
	}

	app, e := service.App.FirstAppByID(f.ID)
	if e != nil {
		if e == gorm.ErrRecordNotFound {
			callback.Error(c, nil, callback.ErrAppNotFound)
			return
		}
		callback.Error(c, e, callback.ErrDBOperation)
		return
	}

	if !app.PermitAllGroup {
		var yes bool
		yes, e = service.App.UserAccessible(f.ID, user.ID)
		if e != nil {
			callback.Error(c, e, callback.ErrDBOperation)
			return
		} else if !yes {
			callback.ErrorWithTip(c, nil, callback.ErrOperationIllegal, "没有访问该应用的权限")
			return
		}
	}

	token, e := jwt.GenerateLoginToken(user.ID, app.ID, user.Name, c.ClientIP(), user.Groups)
	if e != nil {
		callback.Error(c, e, callback.ErrUnexpected)
		return
	}

	callbackUrl, e := tools.GenCallback(app.Callback, token)
	if e != nil {
		callback.Error(c, e, callback.ErrUnexpected)
		return
	}

	callback.Success(c, gin.H{
		"url": callbackUrl,
	})
}
