package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/ncuhome/GeniusAuthoritarian/internal/api/callback"
	"github.com/ncuhome/GeniusAuthoritarian/internal/db/dao"
	"github.com/ncuhome/GeniusAuthoritarian/internal/pkg/dingTalk"
	"github.com/ncuhome/GeniusAuthoritarian/internal/pkg/feishu"
	"github.com/ncuhome/GeniusAuthoritarian/internal/pkg/jwt"
	"github.com/ncuhome/GeniusAuthoritarian/internal/service"
	"gorm.io/gorm"
	"net/url"
)

// 第三方登录 app 路由名称
const (
	appDingTalk = "dingTalk"
	appFeishu   = "feishu"
)

func GetLoginLink(c *gin.Context) {
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

	var link string
	switch c.Param("app") {
	case appFeishu:
		link = feishu.Api.LoginLink(c.Request.Host, f.AppCode)
	case appDingTalk:
		link = dingTalk.Api.LoginLink(c.Request.Host, f.AppCode)
	default:
		callback.Error(c, nil, callback.ErrForm)
		return
	}

	callback.Success(c, gin.H{
		"url": link,
	})
}

func ThirdPartyLogin(c *gin.Context) {
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

	var userPhone string
	switch c.Param("app") {
	case appDingTalk:
		loginDingTalk(c, f.Code)
	case appFeishu:
		loginFeishu(c, f.Code)
	default:
		callback.Error(c, nil, callback.ErrForm)
		return
	}
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
		appInfo = service.App.This(c.Request.Host)
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
