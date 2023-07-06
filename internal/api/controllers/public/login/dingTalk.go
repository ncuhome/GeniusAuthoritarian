package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/ncuhome/GeniusAuthoritarian/internal/api/callback"
	"github.com/ncuhome/GeniusAuthoritarian/internal/db/dto"
	"github.com/ncuhome/GeniusAuthoritarian/internal/pkg/dingTalk"
	"strings"
)

func loginDingTalk(c *gin.Context, code string) *dto.UserThirdPartyIdentity {
	userToken, e := dingTalk.Api.GetUserToken(code)
	if e != nil {
		callback.Error(c, callback.ErrRemoteOperationFailed, e)
		return nil
	}
	userInfo, e := dingTalk.Api.GetUserInfo(*userToken.Body.AccessToken)
	if e != nil {
		callback.Error(c, callback.ErrRemoteOperationFailed, e)
		return nil
	}

	if userInfo.Body.Mobile == nil || *userInfo.Body.Mobile == "" {
		callback.Error(c, callback.ErrUnexpected)
		return nil
	}

	phone := *userInfo.Body.Mobile
	if !strings.HasPrefix(phone, "+") {
		phone = "+86" + phone
	}

	var avatarUrl string
	if userInfo.Body.AvatarUrl != nil {
		avatarUrl = *userInfo.Body.AvatarUrl
	}

	return &dto.UserThirdPartyIdentity{
		Phone:     phone,
		AvatarUrl: avatarUrl,
	}
}
