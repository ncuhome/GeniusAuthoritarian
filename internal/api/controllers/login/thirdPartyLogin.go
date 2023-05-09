package controllers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/ncuhome/GeniusAuthoritarian/internal/api/callback"
	"github.com/ncuhome/GeniusAuthoritarian/internal/db/dao"
	"github.com/ncuhome/GeniusAuthoritarian/internal/pkg/jwt"
	"github.com/ncuhome/GeniusAuthoritarian/internal/service"
	"gorm.io/gorm"
	"net/url"
)

func GetLoginLink(linkGen func(host, appCode string) string) gin.HandlerFunc {
	return func(c *gin.Context) {
		var f struct {
			AppCode string `json:"appCode" form:"appCode"`
		}
		if e := c.ShouldBind(&f); e != nil {
			callback.Error(c, e, callback.ErrForm)
			return
		}

		if f.AppCode != "" {
			if ok, e := service.App.CheckAppCode(f.AppCode); e != nil {
				callback.Error(c, e, callback.ErrDBOperation)
				return
			} else if !ok {
				callback.Error(c, e, callback.ErrAppCodeNotFound)
				return
			}
		}

		callback.Success(c, gin.H{
			"url": linkGen(c.Request.Host, f.AppCode),
		})
	}
}

func ThirdPartyLogin(userInfo func(c *gin.Context, code string) (phone string)) gin.HandlerFunc {
	return func(c *gin.Context) {
		var f struct {
			Code    string `json:"code" form:"code" binding:"required"`
			AppCode string `json:"appCode" form:"appCode"` // 为空时为登录后台
		}
		if e := c.ShouldBind(&f); e != nil {
			callback.Error(c, e, callback.ErrForm)
			return
		}

		if f.AppCode != "" {
			if ok, e := service.App.CheckAppCode(f.AppCode); e != nil {
				callback.Error(c, e, callback.ErrDBOperation)
				return
			} else if !ok {
				callback.Error(c, e, callback.ErrAppCodeNotFound)
				return
			}
		}

		userPhone := userInfo(c, f.Code)
		if c.IsAborted() {
			return
		} else if userPhone == "" {
			callback.Error(c, nil, callback.ErrUnexpected)
			return
		}

		var appInfo *dao.App
		var e error
		if f.AppCode != "" {
			appInfo, e = service.App.FistAppForLogin(f.AppCode)
			if e != nil {
				callback.Error(c, e, callback.ErrDBOperation)
				return
			}
		} else {
			appInfo = &dao.App{
				Name:           "统一鉴权控制系统",
				Callback:       fmt.Sprintf("https://%s/login", c.Request.Host),
				PermitAllGroup: true,
			}
		}

		var user *dao.User
		var groups []dao.Group
		if f.AppCode != "" {
			user, groups, e = service.User.UserInfoForAppCode(userPhone, f.AppCode)
			if e != nil {
				if e == gorm.ErrRecordNotFound {
					callback.ErrorWithTip(c, nil, callback.ErrUnauthorized, "没有找到角色，请尝试使用其他登录方式或联系管理员")
					return
				}
				callback.Error(c, e, callback.ErrDBOperation)
				return
			} else if len(groups) == 0 && !appInfo.PermitAllGroup {
				callback.Error(c, nil, callback.ErrFindUnit)
				return
			}
		} else {
			user, e = service.User.UserInfo(userPhone)
			if e != nil {
				callback.Error(c, e, callback.ErrDBOperation)
				return
			}
		}

		var groupSlice = make([]string, len(groups))
		for i, g := range groups {
			groupSlice[i] = g.Name
		}

		token, e := jwt.GenerateLoginToken(user.ID, appInfo.ID, user.Name, c.ClientIP(), groupSlice)
		if e != nil {
			callback.Error(c, e, callback.ErrUnexpected)
			return
		}

		callbackUrl, e := url.Parse(appInfo.Callback)
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
