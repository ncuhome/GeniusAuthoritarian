package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/ncuhome/GeniusAuthoritarian/internal/api/callback"
	"github.com/ncuhome/GeniusAuthoritarian/internal/pkg/dingTalk"
	"strings"
)

func loginDingTalk(c *gin.Context, code string) string {
	userToken, e := dingTalk.Api.GetUserToken(code)
	if e != nil {
		callback.Error(c, e, callback.ErrRemoteOperationFailed)
		return ""
	}
	userInfo, e := dingTalk.Api.GetUserInfo(*userToken.Body.AccessToken)
	if e != nil {
		callback.Error(c, e, callback.ErrRemoteOperationFailed)
		return ""
	}

	if userInfo.Body.Mobile == nil || *userInfo.Body.Mobile == "" {
		callback.Error(c, nil, callback.ErrUnexpected)
		return ""
	}

	phone := *userInfo.Body.Mobile
	if !strings.HasPrefix(phone, "+") {
		phone = "+86" + phone
	}
	return phone
}
