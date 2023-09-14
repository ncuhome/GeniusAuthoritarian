package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/ncuhome/GeniusAuthoritarian/internal/api/callback"
	"github.com/ncuhome/GeniusAuthoritarian/internal/db/dto"
	"github.com/ncuhome/GeniusAuthoritarian/internal/pkg/feishu"
)

func loginFeishu(c *gin.Context, code string) *dto.UserThirdPartyIdentity {
	user, err := feishu.Api.GetUser(code)
	if err != nil {
		callback.Error(c, callback.ErrRemoteOperationFailed, err)
		return nil
	}
	return &dto.UserThirdPartyIdentity{
		Phone:     user.Mobile,
		AvatarUrl: user.AvatarUrl,
	}
}
