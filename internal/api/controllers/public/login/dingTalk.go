package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/ncuhome/GeniusAuthoritarian/internal/api/callback"
	"github.com/ncuhome/GeniusAuthoritarian/internal/db/dto"
	"github.com/ncuhome/GeniusAuthoritarian/internal/pkg/dingTalk"
	"strings"
)

func loginDingTalk(c *gin.Context, code string) *dto.UserThirdPartyIdentity {
	userToken, err := dingTalk.Api.GetUserToken(code)
	if err != nil {
		callback.Error(c, callback.ErrRemoteOperationFailed, err)
		return nil
	}
	userInfo, err := dingTalk.Api.GetUserInfo(*userToken.Body.AccessToken)
	if err != nil {
		callback.Error(c, callback.ErrRemoteOperationFailed, err)
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

	return &dto.UserThirdPartyIdentity{
		Phone: phone,
	}
}
