package controllers

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/go-webauthn/webauthn/protocol"
	"github.com/go-webauthn/webauthn/webauthn"
	"github.com/ncuhome/GeniusAuthoritarian/internal/api/callback"
	"github.com/ncuhome/GeniusAuthoritarian/internal/db/redis"
	"github.com/ncuhome/GeniusAuthoritarian/internal/pkg/webAuthn"
	"strconv"
	"time"
	"unsafe"
)

func BeginPasskeyLogin(c *gin.Context) {
	options, sessionData, err := webAuthn.Client.BeginDiscoverableLogin()
	if err != nil {
		callback.Error(c, callback.ErrUnexpected, err)
		return
	}

	err = redis.NewPasskey().
		StoreSession(context.Background(), c.ClientIP(), sessionData, time.Minute*5)
	if err != nil {
		callback.Error(c, callback.ErrDBOperation, err)
		return
	}

	callback.Success(c, options)
}

func FinishPasskeyLogin(c *gin.Context) {
	//todo 获取应用信息、验证权限等
	var f struct {
		Credential *protocol.CredentialAssertionResponse `json:"credential" binding:"required"`
	}
	if err := c.ShouldBind(&f); err != nil {
		callback.Error(c, callback.ErrForm, err)
		return
	}
	parsedCredential, err := f.Credential.Parse()
	if err != nil {
		callback.Error(c, callback.ErrForm, err)
		return
	}

	var sessionData webauthn.SessionData
	err = redis.NewPasskey().
		ReadSession(context.Background(), c.ClientIP(), &sessionData)
	if err != nil {
		if err == redis.Nil {
			callback.Error(c, callback.ErrLoginSessionExpired)
			return
		}
		callback.Error(c, callback.ErrDBOperation, err)
		return
	}

	_, err = webAuthn.Client.ValidateDiscoverableLogin(func(_, userHandle []byte) (user webauthn.User, err error) {
		userId, err := strconv.ParseUint(unsafe.String(unsafe.SliceData(userHandle), len(userHandle)), 10, 64)
		if err != nil {
			return nil, err
		}
		return webAuthn.NewUser(uint(userId))
	}, sessionData, parsedCredential)
	if err != nil {
		callback.Error(c, callback.ErrUnauthorized, err)
		return
	}

	//todo 返回 token 与 callback

	callback.Default(c)
}
