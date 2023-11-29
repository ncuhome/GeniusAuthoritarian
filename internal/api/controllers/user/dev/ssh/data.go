package controllers

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/ncuhome/GeniusAuthoritarian/internal/api/callback"
	"github.com/ncuhome/GeniusAuthoritarian/internal/service"
	"github.com/ncuhome/GeniusAuthoritarian/internal/tools"
	"gorm.io/gorm"
)

func ShowSshKey(c *gin.Context) {
	data, err := service.UserSsh.FirstSshSecretsForUserShow(tools.GetUserInfo(c).UID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			callback.Error(c, callback.ErrSshNotFound)
			return
		}
		callback.Error(c, callback.ErrDBOperation, err)
		return
	}

	callback.Success(c, data)
}
