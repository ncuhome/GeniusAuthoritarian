package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/ncuhome/GeniusAuthoritarian/internal/api/callback"
	"github.com/ncuhome/GeniusAuthoritarian/internal/pkg/feishu"
	"github.com/ncuhome/GeniusAuthoritarian/internal/pkg/jwt"
	"github.com/ncuhome/GeniusAuthoritarian/internal/service"
	log "github.com/sirupsen/logrus"
	"net/url"
)

func FeishuLogin(c *gin.Context) {
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

	user, e := feishu.Api.GetUser(f.Code)
	if e != nil {
		callback.Error(c, callback.ErrRemoteOperationFailed)
		return
	}

	userInfo, e := user.Info()
	if e != nil {
		callback.Error(c, callback.ErrRemoteOperationFailed)
		return
	}

	groups, e := service.FeishuGroups.Search(userInfo.User.DepartmentIds)
	if e != nil {
		callback.Error(c, callback.ErrDBOperation)
		return
	} else if len(groups) == 0 {
		callback.Error(c, callback.ErrFindUnit)
		return
	}

	var groupSlice = make([]string, len(groups))
	for i, group := range groups {
		groupSlice[i] = group.Name
	}

	token, e := jwt.GenerateLoginToken(userInfo.User.Name, groupSlice)
	if e != nil {
		log.Debugln("jwt generate failed:", e)
		callback.Error(c, callback.ErrUnexpected)
		return
	}

	callbackUrl, e := url.Parse(f.Callback)
	if e != nil {
		log.Debugln(e)
		callback.Error(c, callback.ErrUnexpected)
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
