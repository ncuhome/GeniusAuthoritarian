package controllers

import (
	"errors"
	"github.com/Mmx233/daoUtil"
	"github.com/gin-gonic/gin"
	"github.com/ncuhome/GeniusAuthoritarian/internal/api/callback"
	"github.com/ncuhome/GeniusAuthoritarian/internal/service"
	"github.com/ncuhome/GeniusAuthoritarian/internal/tools"
	"gorm.io/gorm"
)

func DeletePasskey(c *gin.Context) {
	var f struct {
		ID uint `json:"id" form:"id" binding:"required"`
	}
	if err := c.ShouldBind(&f); err != nil {
		callback.Error(c, callback.ErrForm, err)
		return
	}

	webAuthnSrv, err := service.WebAuthn.Begin()
	if err != nil {
		callback.Error(c, callback.ErrDBOperation, err)
		return
	}
	defer webAuthnSrv.Rollback()

	err = webAuthnSrv.Delete(f.ID, tools.GetUserInfo(c).UID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			callback.Error(c, callback.ErrPasskeyNotExist, err)
			return
		}
		callback.Error(c, callback.ErrDBOperation, err)
		return
	}

	if err = webAuthnSrv.Commit().Error; err != nil {
		callback.Error(c, callback.ErrDBOperation, err)
		return
	}

	callback.Default(c)
}

func RenamePasskey(c *gin.Context) {
	var f struct {
		ID   uint   `json:"id" form:"id" binding:"required"`
		Name string `json:"name" form:"name" binding:"required,max=15"`
	}
	if err := c.ShouldBind(&f); err != nil {
		callback.Error(c, callback.ErrForm, err)
		return
	}

	webAuthnSrv, err := service.WebAuthn.Begin()
	if err != nil {
		callback.Error(c, callback.ErrDBOperation, err)
		return
	}
	defer webAuthnSrv.Rollback()

	exist, err := webAuthnSrv.Exist(f.ID, tools.GetUserInfo(c).UID, daoUtil.LockForUpdate)
	if err != nil {
		callback.Error(c, callback.ErrDBOperation, err)
		return
	} else if !exist {
		callback.Error(c, callback.ErrPasskeyNotExist)
		return
	}

	if err = webAuthnSrv.Rename(f.ID, f.Name); err != nil {
		callback.Error(c, callback.ErrDBOperation, err)
		return
	}

	if err = webAuthnSrv.Commit().Error; err != nil {
		callback.Error(c, callback.ErrDBOperation, err)
		return
	}

	callback.Default(c)
}
