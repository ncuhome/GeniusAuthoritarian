package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/ncuhome/GeniusAuthoritarian/internal/api/callback"
	"github.com/ncuhome/GeniusAuthoritarian/internal/service"
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

	callbackStr, e := service.App.FirstAppCallbackByID(f.ID)
	if e != nil {
		if e == gorm.ErrRecordNotFound {
			callback.Error(c, callback.ErrAppNotFound)
			return
		}
		callback.Error(c, callback.ErrDBOperation, e)
		return
	}

	callbackUrl, e := url.Parse(callbackStr)
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
