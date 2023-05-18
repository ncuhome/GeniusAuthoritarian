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
		ID uint `json:"id" form:"id" binding:"required"`
	}
	if e := c.ShouldBind(&f); e != nil {
		callback.Error(c, e, callback.ErrForm)
		return
	}

	user := tools.GetUserInfo(c)

	app, e := service.App.FirstAccessibleAppByID(f.ID, user.ID)
	if e != nil {
		if e == gorm.ErrRecordNotFound {
			callback.Error(c, nil, callback.ErrAppNotFound)
			return
		}
		callback.Error(c, e, callback.ErrDBOperation)
		return
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
