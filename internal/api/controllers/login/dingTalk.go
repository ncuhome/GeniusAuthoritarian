package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/ncuhome/GeniusAuthoritarian/internal/api/callback"
	"github.com/ncuhome/GeniusAuthoritarian/internal/pkg/dingTalk"
	"github.com/ncuhome/GeniusAuthoritarian/internal/service"
	"strings"
)

var DingTalkLoginLink = GetLoginLink(dingTalk.Api.LoginLink)

var DingTalkLogin = Login(func(c *gin.Context, code string) (string, []string) {
	userToken, e := dingTalk.Api.GetUserToken(code)
	if e != nil {
		callback.Error(c, callback.ErrRemoteOperationFailed)
		return "", nil
	}
	userInfo, e := dingTalk.Api.GetUserInfo(*userToken.Body.AccessToken)
	if e != nil {
		callback.Error(c, callback.ErrRemoteOperationFailed)
		return "", nil
	}

	if userInfo.Body.Mobile == nil || *userInfo.Body.Mobile == "" {
		callback.Error(c, callback.ErrUnexpected)
		return "", nil
	}

	phone := *userInfo.Body.Mobile
	if !strings.HasPrefix(phone, "+") {
		phone = "+86" + phone
	}
	user, groups, e := service.User.UserInfo(phone)
	if e != nil {
		callback.Error(c, callback.ErrDBOperation)
		return "", nil
	}

	var groupSlice = make([]string, len(groups))
	for i, group := range groups {
		groupSlice[i] = group.Name
	}

	return user.Name, groupSlice
})
