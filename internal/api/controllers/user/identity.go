package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/ncuhome/GeniusAuthoritarian/internal/api/callback"
	"github.com/ncuhome/GeniusAuthoritarian/internal/service"
	"github.com/ncuhome/GeniusAuthoritarian/internal/tools"
)

func SendVerifySms(c *gin.Context) {
	uid := tools.GetUserInfo(c).ID
	phone, err := service.User.FirstPhoneByID(uid)
	if err != nil {
		callback.Error(c, callback.ErrDBOperation)
		return
	}

}
