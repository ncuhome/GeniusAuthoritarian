package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/ncuhome/GeniusAuthoritarian/internal/api/callback"
	"github.com/ncuhome/GeniusAuthoritarian/internal/pkg/jwt"
	"github.com/ncuhome/GeniusAuthoritarian/internal/service"
	"github.com/ncuhome/GeniusAuthoritarian/tools"
	log "github.com/sirupsen/logrus"
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

func Login(userInfo func(c *gin.Context, code string) (name string, groups []string)) gin.HandlerFunc {
	return func(c *gin.Context) {
		var f struct {
			Code     string `json:"code" form:"code" binding:"required"`
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

		name, groups := userInfo(c, f.Code)
		if c.IsAborted() {
			return
		}

		token, e := jwt.GenerateLoginToken(name, groups)
		if e != nil {
			log.Debugln("jwt generate failed:", e)
			callback.Error(c, callback.ErrUnexpected)
			return
		}

		callbackUrl, e := tools.AddTokenToUrlQuery(f.Callback, token)
		if e != nil {
			log.Debugln(e)
			callback.Error(c, callback.ErrUnexpected)
			return
		}

		callback.Success(c, gin.H{
			"token":    token,
			"callback": callbackUrl,
		})
	}
}
