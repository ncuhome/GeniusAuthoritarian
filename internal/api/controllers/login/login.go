package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/ncuhome/GeniusAuthoritarian/internal/api/callback"
	"github.com/ncuhome/GeniusAuthoritarian/internal/pkg/jwt"
	"github.com/ncuhome/GeniusAuthoritarian/internal/service"
	"gorm.io/gorm"
	"net/url"
)

func GetLoginLink(linkGen func(host, callback string) string) gin.HandlerFunc {
	return func(c *gin.Context) {
		var f struct {
			Callback string `json:"callback" form:"callback" binding:"required,uri"`
		}
		if e := c.ShouldBind(&f); e != nil {
			callback.Error(c, e, callback.ErrForm)
			return
		}

		if ok, e := service.SiteWhiteList.CheckUrl(f.Callback); e != nil {
			callback.Error(c, e, callback.ErrDBOperation)
			return
		} else if !ok {
			callback.Error(c, e, callback.ErrSiteNotAllow)
			return
		}

		callback.Success(c, gin.H{
			"url": linkGen(c.Request.Host, f.Callback),
		})
	}
}

func Login(userInfo func(c *gin.Context, code string) (phone string)) gin.HandlerFunc {
	return func(c *gin.Context) {
		var f struct {
			Code     string `json:"code" form:"code" binding:"required"`
			Callback string `json:"callback" form:"callback" binding:"required,uri"`
		}
		if e := c.ShouldBind(&f); e != nil {
			callback.Error(c, e, callback.ErrForm)
			return
		}

		if ok, e := service.SiteWhiteList.CheckUrl(f.Callback); e != nil {
			callback.Error(c, e, callback.ErrDBOperation)
			return
		} else if !ok {
			callback.Error(c, e, callback.ErrSiteNotAllow)
			return
		}

		userPhone := userInfo(c, f.Code)
		if c.IsAborted() {
			return
		} else if userPhone == "" {
			callback.Error(c, nil, callback.ErrUnexpected)
			return
		}

		user, groups, e := service.User.UserInfo(userPhone)
		if e != nil {
			if e == gorm.ErrRecordNotFound {
				callback.Error(c, nil, callback.ErrUnauthorized)
				return
			}
			callback.Error(c, e, callback.ErrDBOperation)
			return
		} else if len(groups) == 0 {
			callback.Error(c, nil, callback.ErrFindUnit)
			return
		}
		var groupSlice = make([]string, len(groups))
		for i, group := range groups {
			groupSlice[i] = group.Name
		}

		callbackUrl, e := url.Parse(f.Callback)
		if e != nil {
			callback.Error(c, e, callback.ErrUnexpected)
			return
		}

		token, e := jwt.GenerateLoginToken(user.ID, user.Name, c.ClientIP(), callbackUrl.Host, groupSlice)
		if e != nil {
			callback.Error(c, e, callback.ErrUnexpected)
			return
		}

		callbackQuery := callbackUrl.Query()
		callbackQuery.Set("token", token)
		callbackUrl.RawQuery = callbackQuery.Encode()

		callback.Success(c, gin.H{
			"token":    token,
			"callback": callbackUrl.String(),
		})
	}
}
