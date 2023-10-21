package middlewares

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/ncuhome/GeniusAuthoritarian/internal/api/callback"
	"github.com/ncuhome/GeniusAuthoritarian/internal/service"
	"github.com/ncuhome/GeniusAuthoritarian/internal/tools"
	"gorm.io/gorm"
)

func RequireMfa(c *gin.Context) {
	var f struct {
		Code string `json:"code" form:"code" binding:"required,len=6,numeric"`
	}
	if err := c.ShouldBind(&f); err != nil {
		callback.Error(c, callback.ErrForm, err)
		return
	}

	uid := tools.GetUserInfo(c).ID

	mfaSecret, err := service.User.FirstMfa(uid)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			callback.Error(c, callback.ErrUnexpected)
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
}

func ShouldMfa(c *gin.Context) {
	var f struct {
		Code string `json:"mfa_code" form:"mfa_code" binding:"numeric"`
	}
	if err := c.ShouldBind(&f); err != nil {
		callback.Error(c, callback.ErrForm, err)
		return
	}

	uid := tools.GetUserInfo(c).ID
	mfaSecret, err := service.User.FirstMfa(uid)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			callback.Error(c, callback.ErrUnexpected)
			return
		}
		callback.Error(c, callback.ErrDBOperation, err)
		return
	}

	if mfaSecret == "" {
		return
	} else if len(f.Code) != 6 {
		callback.Error(c, callback.ErrMfaRequired)
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
}
