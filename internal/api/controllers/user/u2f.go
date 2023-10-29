package controllers

import (
	"context"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/go-webauthn/webauthn/webauthn"
	"github.com/ncuhome/GeniusAuthoritarian/internal/api/callback"
	"github.com/ncuhome/GeniusAuthoritarian/internal/db/redis"
	"github.com/ncuhome/GeniusAuthoritarian/internal/pkg/jwt"
	"github.com/ncuhome/GeniusAuthoritarian/internal/pkg/webAuthn"
	"github.com/ncuhome/GeniusAuthoritarian/internal/service"
	"github.com/ncuhome/GeniusAuthoritarian/internal/tools"
	"gorm.io/gorm"
)

// U2F 是后台已登录用户身份二次校验的总方法

func AvailableU2fMethod(c *gin.Context) {
	data, err := service.User.U2fStatus(tools.GetUserInfo(c).ID)
	if err != nil {
		callback.Error(c, callback.ErrDBOperation, err)
		return
	}

	callback.Success(c, data)
}

// BeginU2F 为用户分发 U2F 短效令牌，可以通过指定需要 U2F 的接口
// 各校验方式的前置准备在其他对应路由组的接口中，本接口直接获取其结果
func BeginU2F(c *gin.Context) {
	uid := tools.GetUserInfo(c).ID

	method := c.Param("method")
	switch method {
	case "phone":
		var f struct {
			Code string `json:"code" form:"code" binding:"required,len=6,numeric"`
		}
		if err := c.ShouldBind(&f); err != nil {
			callback.Error(c, callback.ErrForm, err)
			return
		}

		ok, err := redis.NewUserIdentityCode(uid).VerifyAndDestroy(f.Code)
		if err != nil {
			callback.Error(c, callback.ErrUnexpected, err)
			return
		} else if !ok {
			callback.Error(c, callback.ErrIdentityCodeNotCorrect)
			return
		}
	case "mfa":
		var f struct {
			Code string `json:"code" form:"code" binding:"required,len=6,numeric"`
		}
		if err := c.ShouldBind(&f); err != nil {
			callback.Error(c, callback.ErrForm, err)
			return
		}

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
	case "passkey":
		var sessionData webauthn.SessionData
		err := redis.NewPasskey(c.ClientIP()).
			ReadSession(context.Background(), &sessionData)
		if err != nil {
			if err == redis.Nil {
				callback.Error(c, callback.ErrLoginSessionExpired)
				return
			}
			callback.Error(c, callback.ErrDBOperation, err)
			return
		}
		user, err := webAuthn.NewUser(uid)
		if err != nil {
			callback.Error(c, callback.ErrDBOperation, err)
			return
		}
		credential, err := webAuthn.Client.FinishLogin(user, sessionData, c.Request)
		if err != nil {
			callback.Error(c, callback.ErrUnauthorized, err)
			return
		}

		webAuthnSrv, err := service.WebAuthn.Begin()
		if err != nil {
			callback.Error(c, callback.ErrDBOperation, err)
			return
		}
		defer webAuthnSrv.Rollback()

		err = webAuthnSrv.UpdateLastUsedAt(uid, credential.Descriptor().CredentialID.String())
		if err != nil {
			callback.Error(c, callback.ErrDBOperation, err)
			return
		}

		if err = webAuthnSrv.Commit().Error; err != nil {
			callback.Error(c, callback.ErrDBOperation, err)
			return
		}
	default:
		c.AbortWithStatus(404)
		return
	}

	token, claims, err := jwt.GenerateU2fToken(uid, c.ClientIP())
	if err != nil {
		callback.Error(c, callback.ErrUnexpected, err)
		return
	}

	callback.Success(c, gin.H{
		"token":        token,
		"valid_before": claims.ExpiresAt.Time.Unix(),
	})
}
