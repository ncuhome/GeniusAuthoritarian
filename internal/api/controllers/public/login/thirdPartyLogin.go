package controllers

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/ncuhome/GeniusAuthoritarian/internal/api/callback"
	"github.com/ncuhome/GeniusAuthoritarian/internal/api/models/response"
	"github.com/ncuhome/GeniusAuthoritarian/internal/db/dao"
	"github.com/ncuhome/GeniusAuthoritarian/internal/db/dto"
	"github.com/ncuhome/GeniusAuthoritarian/internal/pkg/dingTalk"
	"github.com/ncuhome/GeniusAuthoritarian/internal/pkg/feishu"
	"github.com/ncuhome/GeniusAuthoritarian/internal/pkg/jwt"
	"github.com/ncuhome/GeniusAuthoritarian/internal/service"
	"github.com/ncuhome/GeniusAuthoritarian/internal/tools"
	"gorm.io/gorm"
)

// 第三方登录 app 路由名称
const (
	appDingTalk = "dingTalk"
	appFeishu   = "feishu"
)

// 生成登录地址
func genLoginLink(c *gin.Context, appCode string) string {
	switch c.Param("app") {
	case appFeishu:
		return feishu.Api.LoginLink(c.Request.Host, appCode)
	case appDingTalk:
		return dingTalk.Api.LoginLink(c.Request.Host, appCode)
	default:
		callback.Error(c, callback.ErrForm)
		return ""
	}
}

// GetSelfLoginLink 获取登录控制系统地址
func GetSelfLoginLink(c *gin.Context) {
	callback.Success(c, gin.H{
		"url": genLoginLink(c, ""),
	})
}

// GetLoginLink 获取登录指定 APP 地址
func GetLoginLink(c *gin.Context) {
	appCode := c.Param("appCode")

	if ok, e := service.App.CheckAppCode(appCode); e != nil {
		callback.Error(c, callback.ErrDBOperation, e)
		return
	} else if !ok {
		callback.Error(c, callback.ErrAppCodeNotFound, e)
		return
	}

	link := genLoginLink(c, appCode)
	if c.IsAborted() {
		return
	}

	callback.Success(c, gin.H{
		"url": link,
	})
}

// 根据前端传来的 code 获取用户电话号码
func loadUserIdentity(c *gin.Context, code string) *dto.UserThirdPartyIdentity {
	switch c.Param("app") {
	case appDingTalk:
		return loginDingTalk(c, code)
	case appFeishu:
		return loginFeishu(c, code)
	default:
		callback.Error(c, callback.ErrForm)
		return nil
	}
}

type ThirdPartyLoginContext struct {
	User      *dao.User
	AppInfo   *dao.App
	Groups    []string
	Ip        string
	AvatarUrl string
}

// 根据数据完成请求响应
func callThirdPartyLoginResult(c *gin.Context, info ThirdPartyLoginContext) {
	claims := jwt.LoginTokenClaims{
		UID:       info.User.ID,
		Name:      info.User.Name,
		IP:        info.Ip,
		Groups:    info.Groups,
		AppID:     info.AppInfo.ID,
		AvatarUrl: info.AvatarUrl,
	}

	if info.User.MFA == "" {
		token, e := jwt.GenerateLoginToken(claims)
		if e != nil {
			callback.Error(c, callback.ErrUnexpected, e)
			return
		}

		callbackUrl, e := tools.GenCallback(info.AppInfo.Callback, token)
		if e != nil {
			callback.Error(c, callback.ErrUnexpected, e)
			return
		}

		callback.Success(c, response.ThirdPartyLogin{
			Token:    token,
			Mfa:      false,
			Callback: callbackUrl,
		})
	} else {
		token, e := jwt.GenerateMfaToken(claims, info.User.MFA, info.AppInfo.Callback)
		if e != nil {
			callback.Error(c, callback.ErrUnexpected, e)
			return
		}

		callback.Success(c, response.ThirdPartyLogin{
			Token: token,
			Mfa:   true,
		})
	}
}

// ThirdPartySelfLogin 校验第三方登录回调结果，生成控制系统回调链接
func ThirdPartySelfLogin(c *gin.Context) {
	var f struct {
		Code string `json:"code" form:"code" binding:"required"`
	}
	if e := c.ShouldBind(&f); e != nil {
		callback.Error(c, callback.ErrForm, e)
		return
	}

	userIdentity := loadUserIdentity(c, f.Code)
	if c.IsAborted() {
		return
	} else if userIdentity == nil {
		callback.Error(c, callback.ErrUnexpected)
		return
	}

	appInfo := service.App.This(c.Request.Host)

	user, e := service.User.UserInfo(userIdentity.Phone)
	if e != nil {
		callback.Error(c, callback.ErrDBOperation, e)
		return
	}

	callThirdPartyLoginResult(c, ThirdPartyLoginContext{
		User:      user,
		AppInfo:   appInfo,
		Groups:    nil,
		Ip:        c.ClientIP(),
		AvatarUrl: userIdentity.AvatarUrl,
	})
}

// ThirdPartyLogin 校验第三方登录回调结果，生成目标 APP 回调链接
func ThirdPartyLogin(c *gin.Context) {
	var f struct {
		Code string `json:"code" form:"code" binding:"required"`
	}
	if e := c.ShouldBind(&f); e != nil {
		callback.Error(c, callback.ErrForm, e)
		return
	}

	appCode := c.Param("appCode")

	if ok, e := service.App.CheckAppCode(appCode); e != nil {
		callback.Error(c, callback.ErrDBOperation, e)
		return
	} else if !ok {
		callback.Error(c, callback.ErrAppCodeNotFound, e)
		return
	}

	userIdentity := loadUserIdentity(c, f.Code)
	if c.IsAborted() {
		return
	} else if userIdentity == nil {
		callback.Error(c, callback.ErrUnexpected)
		return
	}

	appInfo, e := service.App.FirstAppByAppCode(appCode)
	if e != nil {
		callback.Error(c, callback.ErrDBOperation, e)
		return
	}

	user, e := service.User.UserInfo(userIdentity.Phone)
	if e != nil {
		callback.Error(c, callback.ErrDBOperation, e)
		return
	}

	var groups []string
	isCenterMember, e := service.UserGroups.IsCenterMember(user.ID)
	if e != nil {
		callback.Error(c, callback.ErrDBOperation, e)
		return
	} else if isCenterMember {
		groups, e = service.UserGroups.GetForUser(user.ID)
		if e != nil {
			callback.Error(c, callback.ErrDBOperation, e)
			return
		}
	} else {
		var baseGroups []dao.BaseGroup
		baseGroups, e = service.UserGroups.GetForAppCode(user.ID, appCode)
		if e != nil {
			if errors.Is(e, gorm.ErrRecordNotFound) {
				callback.ErrorWithTip(c, callback.ErrUnauthorized, "没有找到角色，请尝试使用其他登录方式或联系管理员")
				return
			}
			callback.Error(c, callback.ErrDBOperation, e)
			return
		}

		if len(groups) == 0 && !appInfo.PermitAllGroup {
			callback.Error(c, callback.ErrFindUnit)
			return
		}

		groups = make([]string, len(baseGroups))
		for i, g := range baseGroups {
			groups[i] = g.Name
		}
	}

	callThirdPartyLoginResult(c, ThirdPartyLoginContext{
		User:      user,
		AppInfo:   appInfo,
		Groups:    groups,
		Ip:        c.ClientIP(),
		AvatarUrl: userIdentity.AvatarUrl,
	})
}
