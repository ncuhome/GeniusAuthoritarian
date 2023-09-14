package controllers

import (
	"errors"
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
	if err := c.ShouldBind(&f); err != nil {
		callback.Error(c, callback.ErrForm, err)
		return
	}

	callbackStr, err := service.App.FirstAppCallbackByID(f.ID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			callback.Error(c, callback.ErrAppNotFound)
			return
		}
		callback.Error(c, callback.ErrDBOperation, err)
		return
	}

	callbackUrl, err := url.Parse(callbackStr)
	if err != nil {
		callback.Error(c, callback.ErrUnexpected, err)
		return
	}

	callbackUrl.Path = ""
	callbackUrl.RawQuery = ""

	callback.Success(c, gin.H{
		"url": callbackUrl.String(),
	})
}
