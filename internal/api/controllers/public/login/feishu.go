package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/ncuhome/GeniusAuthoritarian/internal/api/callback"
	"github.com/ncuhome/GeniusAuthoritarian/internal/db/dto"
	"github.com/ncuhome/GeniusAuthoritarian/internal/pkg/feishu"
)

func loginFeishu(c *gin.Context, code string) *dto.UserThirdPartyIdentity {
	user, e := feishu.Api.GetUser(code)
	if e != nil {
		callback.Error(c, callback.ErrRemoteOperationFailed, e)
		return nil
	}
	return &dto.UserThirdPartyIdentity{
		Phone:     user.Mobile,
		AvatarUrl: user.AvatarUrl,
	}
}
