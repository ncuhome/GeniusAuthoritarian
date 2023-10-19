package controllers

import (
	"encoding/json"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/go-webauthn/webauthn/protocol"
	"github.com/go-webauthn/webauthn/webauthn"
	"github.com/ncuhome/GeniusAuthoritarian/internal/api/callback"
	"github.com/ncuhome/GeniusAuthoritarian/internal/pkg/webAuthn"
	log "github.com/sirupsen/logrus"
	"unsafe"
)

func BeginPasskeyLogin(c *gin.Context) {
	options, sessionData, err := webAuthn.Client.BeginDiscoverableLogin()
	if err != nil {
		callback.Error(c, callback.ErrUnexpected, err)
		return
	}

	sessionDataBytes, err := json.Marshal(&sessionData)
	if err != nil {
		callback.Error(c, callback.ErrUnexpected, err)
		return
	}
	sessionDataStr := unsafe.String(unsafe.SliceData(sessionDataBytes), len(sessionDataBytes))

	session := sessions.Default(c)
	session.Set("passkey-login", sessionDataStr)
	if err = session.Save(); err != nil {
		callback.Error(c, callback.ErrUnexpected, err)
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

	session := sessions.Default(c)
	sessionDataInterface := session.Get("passkey-login")
	if sessionDataInterface == nil {
		callback.Error(c, callback.ErrLoginSessionExpired)
		return
	}
	sessionDataStr, ok := sessionDataInterface.(string)
	if !ok {
		callback.Error(c, callback.ErrUnexpected)
		return
	}
	sessionDataBytes := unsafe.Slice(unsafe.StringData(sessionDataStr), len(sessionDataStr))
	var sessionData webauthn.SessionData
	if err = json.Unmarshal(sessionDataBytes, &sessionData); err != nil {
		callback.Error(c, callback.ErrUnexpected, err)
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
