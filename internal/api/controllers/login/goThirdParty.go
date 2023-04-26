package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/ncuhome/GeniusAuthoritarian/internal/api/callback"
	"github.com/ncuhome/GeniusAuthoritarian/internal/service"
)

func GetLoginLink(linkGen func(host, callback string) string) gin.HandlerFunc {
	return func(c *gin.Context) {
		var f struct {
			Callback string `json:"callback" form:"callback" binding:"required,uri"`
		}
		if e := c.ShouldBind(&f); e != nil {
			callback.Error(c, callback.ErrForm)
			return
		}

		if ok, e := service.SiteWhiteList.CheckUrl(f.Callback); e != nil {
			callback.Error(c, callback.ErrDBOperation)
			return
		} else if !ok {
			callback.Error(c, callback.ErrSiteNotAllow)
			return
		}

		callback.Success(c, gin.H{
			"url": linkGen(c.Request.Host, f.Callback),
		})
	}
}
