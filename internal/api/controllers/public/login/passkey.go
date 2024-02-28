package controllers

import (
	"context"
	"errors"
	"github.com/Mmx233/tool"
	"github.com/gin-gonic/gin"
	"github.com/go-webauthn/webauthn/protocol"
	"github.com/go-webauthn/webauthn/webauthn"
	"github.com/ncuhome/GeniusAuthoritarian/internal/api/callback"
	"github.com/ncuhome/GeniusAuthoritarian/internal/db/dao"
	"github.com/ncuhome/GeniusAuthoritarian/internal/db/redis"
	"github.com/ncuhome/GeniusAuthoritarian/internal/pkg/jwt"
	"github.com/ncuhome/GeniusAuthoritarian/internal/pkg/jwt/jwtClaims"
	"github.com/ncuhome/GeniusAuthoritarian/internal/pkg/webAuthn"
	"github.com/ncuhome/GeniusAuthoritarian/internal/service"
	"github.com/ncuhome/GeniusAuthoritarian/internal/tools"
	"gorm.io/gorm"
	"math/rand"
	"strconv"
	"time"
	"unsafe"
)

const passkeyCookieName = "session-key"
const passkeyKeyLength = 6

func BeginPasskeyLogin(c *gin.Context) {
	options, sessionData, err := webAuthn.Client.BeginDiscoverableLogin()
	if err != nil {
		callback.Error(c, callback.ErrUnexpected, err)
		return
	}

	identity := "r" + tool.NewRand(rand.NewSource(time.Now().UnixNano())).
		WithLetters("qwertyuiopasdfghjklzxcvbnmQWERTYUIOPASDFGHJKLZXCVBNM").
		String(passkeyKeyLength-1)

	err = redis.NewPasskey(c.ClientIP(), redis.PasskeyLogin, identity).
		StoreSession(context.Background(), sessionData, time.Minute*5)
	if err != nil {
		callback.Error(c, callback.ErrDBOperation, err)
		return
	}

	c.SetCookie(passkeyCookieName, identity, int((time.Minute * 5).Seconds()), "", "", true, true)
	callback.Success(c, options)
}

func FinishPasskeyLogin(c *gin.Context) {
	var f struct {
		AppCode    string                                `json:"app_code" form:"app_code"`
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

	identity, err := c.Cookie(passkeyCookieName)
	if err != nil || len(identity) != passkeyKeyLength {
		callback.Error(c, callback.ErrLoginSessionExpired, err)
		return
	}

	var sessionData webauthn.SessionData
	err = redis.NewPasskey(c.ClientIP(), redis.PasskeyLogin, identity).
		ReadSession(context.Background(), &sessionData)
	if err != nil {
		if errors.Is(err, redis.Nil) {
			callback.Error(c, callback.ErrLoginSessionExpired)
			return
		}
		callback.Error(c, callback.ErrDBOperation, err)
		return
	}

	var uid uint
	credential, err := webAuthn.Client.ValidateDiscoverableLogin(func(_, userHandle []byte) (user webauthn.User, err error) {
		userId, err := strconv.ParseUint(unsafe.String(unsafe.SliceData(userHandle), len(userHandle)), 10, 64)
		if err != nil {
			return nil, err
		}
		uid = uint(userId)
		return webAuthn.NewUser(uid)
	}, sessionData, parsedCredential)
	if err != nil {
		callback.Error(c, callback.ErrUnauthorized, err)
		return
	}

	var appInfo *dao.App
	if f.AppCode == "" {
		appInfo = service.App.This(c.Request.Host)
	} else {
		appInfo, err = service.App.FirstAppByAppCode(f.AppCode)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				callback.Error(c, callback.ErrAppCodeNotFound)
				return
			}
			callback.Error(c, callback.ErrDBOperation, err)
			return
		}
	}

	user, err := service.User.UserInfoByID(uid)
	if err != nil {
		callback.Error(c, callback.ErrDBOperation, err)
		return
	}

	groups, ok := checkUserPermission(c, user, f.AppCode, appInfo.PermitAllGroup)
	if !ok {
		return
	}

	token, err := jwt.GenerateLoginToken(jwtClaims.LoginRedis{
		UID:       user.ID,
		Name:      user.Name,
		IP:        c.ClientIP(),
		Useragent: c.Request.UserAgent(),
		Groups:    groups,
		AppID:     appInfo.ID,
		AvatarUrl: user.AvatarUrl,
		Method:    "passkey",
	})
	if err != nil {
		callback.Error(c, callback.ErrUnexpected, err)
		return
	}

	callbackUrl, err := tools.GenCallback(appInfo.Callback, token)
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

	err = webAuthnSrv.UpdateLastUsedAt(uid, credential.Descriptor().CredentialID.String())
	if err != nil {
		callback.Error(c, callback.ErrDBOperation, err)
		return
	}

	if err = webAuthnSrv.Commit().Error; err != nil {
		callback.Error(c, callback.ErrDBOperation, err)
		return
	}

	callback.Success(c, gin.H{
		"token":    token,
		"callback": callbackUrl,
	})
}
