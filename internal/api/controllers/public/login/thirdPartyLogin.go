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
	"github.com/ncuhome/GeniusAuthoritarian/pkg/departments"
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

func checkUserPermission(c *gin.Context, user *dao.User, appCode string, permitAllGroup bool) (groups []string, ok bool) {
	isCenterMember, err := service.UserGroups.IsUnitMember(user.ID, departments.UCe)
	if err != nil {
		callback.Error(c, callback.ErrDBOperation, err)
		return
	} else if isCenterMember {
		groups, err = service.UserGroups.GetNamesForUser(user.ID)
		if err != nil {
			callback.Error(c, callback.ErrDBOperation, err)
			return
		}
	} else {
		groups, err = service.UserGroups.GetForAppCode(user.ID, appCode)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				callback.Error(c, callback.ErrUserIdentity)
				return
			}
			callback.Error(c, callback.ErrDBOperation, err)
			return
		}

		if len(groups) == 0 && !permitAllGroup {
			callback.Error(c, callback.ErrFindUnit)
			return
		}
	}

	ok = true
	return
}

type ThirdPartyLoginContext struct {
	User    *dao.User
	AppInfo *dao.App
	Groups  []string
	Ip      string
}

// 根据数据完成请求响应
func callThirdPartyLoginResult(c *gin.Context, info ThirdPartyLoginContext) {
	claims := jwt.LoginRedisClaims{
		UID:       info.User.ID,
		Name:      info.User.Name,
		IP:        info.Ip,
		Groups:    info.Groups,
		AppID:     info.AppInfo.ID,
		AvatarUrl: info.User.AvatarUrl,
	}

	if info.User.MFA == "" {
		token, err := jwt.GenerateLoginToken(claims)
		if err != nil {
			callback.Error(c, callback.ErrUnexpected, err)
			return
		}

		accessToken, err := jwt.GenerateAccessToken(info.User.ID, info.User.Name, info.AppInfo.AppCode, info.Groups)
		if err != nil {
			callback.Error(c, callback.ErrUnexpected, err)
			return
		}

		refreshToken, err := jwt.GenerateRefreshToken(info.User.ID, info.User.Name, info.AppInfo.AppCode, info.Groups)
		if err != nil {
			callback.Error(c, callback.ErrUnexpected, err)
			return
		}

		callbackUrl, err := tools.GenCallback(info.AppInfo.Callback, token)
		if err != nil {
			callback.Error(c, callback.ErrUnexpected, err)
			return
		}

		callback.Success(c, response.ThirdPartyLogin{
			Token:        token,
			Mfa:          false,
			Callback:     callbackUrl,
			AccessToken:  accessToken,
			RefreshToken: refreshToken,
		})
	} else {
		token, err := jwt.GenerateMfaToken(claims, info.User.MFA, info.AppInfo.Callback)
		if err != nil {
			callback.Error(c, callback.ErrUnexpected, err)
			return
		}

		accessToken, err := jwt.GenerateAccessToken(info.User.ID, info.User.Name, info.AppInfo.AppCode, info.Groups)
		if err != nil {
			callback.Error(c, callback.ErrUnexpected, err)
			return
		}

		refreshToken, err := jwt.GenerateRefreshToken(info.User.ID, info.User.Name, info.AppInfo.AppCode, info.Groups)
		if err != nil {
			callback.Error(c, callback.ErrUnexpected, err)
			return
		}

		callback.Success(c, response.ThirdPartyLogin{
			Token:        token,
			Mfa:          true,
			AccessToken:  accessToken,
			RefreshToken: refreshToken,
		})
	}
}

// ThirdPartySelfLogin 校验第三方登录回调结果，生成控制系统回调链接
func ThirdPartySelfLogin(c *gin.Context) {
	var f struct {
		Code string `json:"code" form:"code" binding:"required"`
	}
	if err := c.ShouldBind(&f); err != nil {
		callback.Error(c, callback.ErrForm, err)
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

	user, err := service.User.FirstByPhone(userIdentity.Phone)
	if err != nil {
		callback.Error(c, callback.ErrDBOperation, err)
		return
	}

	callThirdPartyLoginResult(c, ThirdPartyLoginContext{
		User:    user,
		AppInfo: appInfo,
		Groups:  nil,
		Ip:      c.ClientIP(),
	})
}

// ThirdPartyLogin 校验第三方登录回调结果，生成目标 APP 回调链接
func ThirdPartyLogin(c *gin.Context) {
	var f struct {
		Code string `json:"code" form:"code" binding:"required"`
	}
	if err := c.ShouldBind(&f); err != nil {
		callback.Error(c, callback.ErrForm, err)
		return
	}

	appCode := c.Param("appCode")

	if ok, err := service.App.CheckAppCode(appCode); err != nil {
		callback.Error(c, callback.ErrDBOperation, err)
		return
	} else if !ok {
		callback.Error(c, callback.ErrAppCodeNotFound)
		return
	}

	userIdentity := loadUserIdentity(c, f.Code)
	if c.IsAborted() {
		return
	} else if userIdentity == nil {
		callback.Error(c, callback.ErrUnexpected)
		return
	}

	appInfo, err := service.App.FirstAppByAppCode(appCode)
	if err != nil {
		callback.Error(c, callback.ErrDBOperation, err)
		return
	}

	user, err := service.User.FirstByPhone(userIdentity.Phone)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			callback.Error(c, callback.ErrUserIdentity)
			return
		}
		callback.Error(c, callback.ErrDBOperation, err)
		return
	}

	groups, ok := checkUserPermission(c, user, appCode, appInfo.PermitAllGroup)
	if !ok {
		return
	}

	callThirdPartyLoginResult(c, ThirdPartyLoginContext{
		User:    user,
		AppInfo: appInfo,
		Groups:  groups,
		Ip:      c.ClientIP(),
	})
}
