package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/ncuhome/GeniusAuthoritarian/internal/api/callback"
	"github.com/ncuhome/GeniusAuthoritarian/internal/pkg/feishu"
)

func loginFeishu(c *gin.Context, code string) string {
	user, e := feishu.Api.GetUser(code)
	if e != nil {
		callback.Error(c, e, callback.ErrRemoteOperationFailed)
		return ""
	}
	return user.Mobile
}
