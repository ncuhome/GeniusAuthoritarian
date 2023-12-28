package controllers

import (
	"context"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-webauthn/webauthn/protocol"
	"github.com/go-webauthn/webauthn/webauthn"
	"github.com/ncuhome/GeniusAuthoritarian/internal/api/callback"
	"github.com/ncuhome/GeniusAuthoritarian/internal/db/redis"
	"github.com/ncuhome/GeniusAuthoritarian/internal/pkg/webAuthn"
	"github.com/ncuhome/GeniusAuthoritarian/internal/service"
	"github.com/ncuhome/GeniusAuthoritarian/internal/tools"
	"time"
)

func BeginPasskeyRegistration(c *gin.Context) {
	uid := tools.GetUserInfo(c).UID

	user, err := webAuthn.NewUser(uid)
	if err != nil {
		callback.Error(c, callback.ErrDBOperation, err)
		return
	}
	options, session, err := webAuthn.Client.BeginRegistration(
		user,
		webauthn.WithAuthenticatorSelection(protocol.AuthenticatorSelection{
			ResidentKey:        protocol.ResidentKeyRequirementRequired,
			RequireResidentKey: protocol.ResidentKeyRequired(),
			UserVerification:   protocol.VerificationRequired,
		}),
		webauthn.WithExtensions(map[string]interface{}{
			"credProps": true,
		}),
		webauthn.WithConveyancePreference(protocol.PreferNoAttestation),
	)
	if err != nil {
		callback.Error(c, callback.ErrUnexpected, err)
	}

	err = redis.NewPasskey(c.ClientIP(), redis.PasskeyUserRegister, fmt.Sprint(uid)).
		StoreSession(context.Background(), session, time.Minute*10)
	if err != nil {
		callback.Error(c, callback.ErrDBOperation, err)
		return
	}

	callback.Success(c, options)
}

func FinishPasskeyRegistration(c *gin.Context) {
	uid := tools.GetUserInfo(c).UID

	var session webauthn.SessionData
	err := redis.NewPasskey(c.ClientIP(), redis.PasskeyUserRegister, fmt.Sprint(uid)).
		ReadSession(context.Background(), &session)
	if err != nil {
		if errors.Is(err, redis.Nil) {
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

	credential, err := webAuthn.Client.FinishRegistration(user, session, c.Request)
	if err != nil {
		callback.Error(c, callback.ErrUnexpected, err)
		return
	}

	webAuthnSrv, err := service.WebAuthn.Begin()
	if err != nil {
		callback.Error(c, callback.ErrDBOperation, err)
		return
	}
	defer webAuthnSrv.Rollback()

	credDto, err := webAuthnSrv.Add(uid, credential)
	if err != nil {
		callback.Error(c, callback.ErrDBOperation, err)
		return
	}

	if err = webAuthnSrv.Commit().Error; err != nil {
		callback.Error(c, callback.ErrDBOperation, err)
		return
	}

	callback.Success(c, credDto)
}
