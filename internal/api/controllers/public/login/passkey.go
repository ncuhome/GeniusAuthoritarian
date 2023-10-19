package controllers

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/go-webauthn/webauthn/protocol"
	"github.com/go-webauthn/webauthn/webauthn"
	"github.com/ncuhome/GeniusAuthoritarian/internal/api/callback"
	"github.com/ncuhome/GeniusAuthoritarian/internal/db/redis"
	"github.com/ncuhome/GeniusAuthoritarian/internal/pkg/webAuthn"
	log "github.com/sirupsen/logrus"
	"time"
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
	parsedResponse, err := protocol.ParseCredentialRequestResponse(c.Request)
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

	// todo 可能需要处理返回的 cred 信息
	_, err = webAuthn.Client.ValidateDiscoverableLogin(func(rawID, userHandle []byte) (user webauthn.User, err error) {
		// todo find user
		log.Debugln(string(rawID))
		log.Debugln(string(userHandle))
		return nil, nil
	}, sessionData, parsedResponse)
	if err != nil {
		callback.Error(c, callback.ErrDBOperation, err)
		return
	}

	callback.Default(c)
}
