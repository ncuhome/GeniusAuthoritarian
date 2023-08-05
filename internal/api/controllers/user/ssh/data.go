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
	var f struct {
		Code string `json:"code" form:"code" binding:"required,len=6,numeric"`
	}
	if e := c.ShouldBind(&f); e != nil {
		callback.Error(c, callback.ErrForm, e)
		return
	}

	uid := tools.GetUserInfo(c).ID

	mfaSecret, err := service.User.FirstMfa(uid)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			callback.Error(c, callback.ErrMfaRequired)
			return
		}
		callback.Error(c, callback.ErrDBOperation, err)
		return
	}

	ok, err := tools.VerifyMfa(f.Code, mfaSecret)
	if err != nil {
		callback.Error(c, callback.ErrUnexpected, err)
		return
	} else if !ok {
		callback.Error(c, callback.ErrMfaCode)
		return
	}

	data, err := service.UserSsh.FirstSshSecretsForUserShow(uid)
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
