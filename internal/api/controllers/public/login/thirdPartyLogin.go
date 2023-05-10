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

func genLoginLink(c *gin.Context, appCode string) string {
	switch c.Param("app") {
	case appFeishu:
		return feishu.Api.LoginLink(c.Request.Host, appCode)
	case appDingTalk:
		return dingTalk.Api.LoginLink(c.Request.Host, appCode)
	default:
		callback.Error(c, nil, callback.ErrForm)
		return ""
	}
}
func GetSelfLoginLink(c *gin.Context) {
	callback.Success(c, gin.H{
		"url": genLoginLink(c, ""),
	})
}
func GetLoginLink(c *gin.Context) {
	appCode := c.Param("appCode")

	if ok, e := service.App.CheckAppCode(appCode); e != nil {
		callback.Error(c, e, callback.ErrDBOperation)
		return
	} else if !ok {
		callback.Error(c, e, callback.ErrAppCodeNotFound)
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

func loadUserPhone(c *gin.Context, code string) string {
	switch c.Param("app") {
	case appDingTalk:
		return loginDingTalk(c, code)
	case appFeishu:
		return loginFeishu(c, code)
	default:
		callback.Error(c, nil, callback.ErrForm)
		return ""
	}
}
func callThirdPartyLoginResult(c *gin.Context, user *dao.User, appInfo *dao.App, groups []dao.Group, ip string) {
	var groupSlice = make([]string, len(groups))
	for i, g := range groups {
		groupSlice[i] = g.Name
	}

	token, e := jwt.GenerateLoginToken(user.ID, appInfo.ID, user.Name, ip, groupSlice)
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
func ThirdPartySelfLogin(c *gin.Context) {
	var f struct {
		Code string `json:"code" form:"code" binding:"required"`
	}
	if e := c.ShouldBind(&f); e != nil {
		callback.Error(c, e, callback.ErrForm)
		return
	}

	userPhone := loadUserPhone(c, f.Code)
	if c.IsAborted() {
		return
	} else if userPhone == "" {
		callback.Error(c, nil, callback.ErrUnexpected)
		return
	}

	appInfo := service.App.This(c.Request.Host)

	user, e := service.User.UserInfo(userPhone)
	if e != nil {
		callback.Error(c, e, callback.ErrDBOperation)
		return
	}

	callThirdPartyLoginResult(c, user, appInfo, nil, c.ClientIP())
}
func ThirdPartyLogin(c *gin.Context) {
	var f struct {
		Code string `json:"code" form:"code" binding:"required"`
	}
	if e := c.ShouldBind(&f); e != nil {
		callback.Error(c, e, callback.ErrForm)
		return
	}

	appCode := c.Param("appCode")

	if ok, e := service.App.CheckAppCode(appCode); e != nil {
		callback.Error(c, e, callback.ErrDBOperation)
		return
	} else if !ok {
		callback.Error(c, e, callback.ErrAppCodeNotFound)
		return
	}

	userPhone := loadUserPhone(c, f.Code)
	if c.IsAborted() {
		return
	} else if userPhone == "" {
		callback.Error(c, nil, callback.ErrUnexpected)
		return
	}

	appInfo, e := service.App.FistAppForLogin(appCode)
	if e != nil {
		callback.Error(c, e, callback.ErrDBOperation)
		return
	}

	user, groups, e := service.User.UserInfoForAppCode(userPhone, appCode)
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

	callThirdPartyLoginResult(c, user, appInfo, groups, c.ClientIP())
}
