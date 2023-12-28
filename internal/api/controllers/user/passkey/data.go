package controllers

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/ncuhome/GeniusAuthoritarian/internal/api/callback"
	"github.com/ncuhome/GeniusAuthoritarian/internal/db/redis"
	"github.com/ncuhome/GeniusAuthoritarian/internal/pkg/webAuthn"
	"github.com/ncuhome/GeniusAuthoritarian/internal/service"
	"github.com/ncuhome/GeniusAuthoritarian/internal/tools"
	"time"
)

func ListPasskey(c *gin.Context) {
	list, err := service.WebAuthn.ListUserCredForShow(tools.GetUserInfo(c).UID)
	if err != nil {
		callback.Error(c, callback.ErrDBOperation, err)
		return
	}

	callback.Success(c, list)
}

func PasskeyOptions(c *gin.Context) {
	uid := tools.GetUserInfo(c).UID
	user, err := webAuthn.NewUser(uid)
	if err != nil {
		callback.Error(c, callback.ErrDBOperation, err)
		return
	}

	options, session, err := webAuthn.Client.BeginLogin(user)
	if err != nil {
		callback.Error(c, callback.ErrUnexpected, err)
		return
	}

	err = redis.NewPasskey(c.ClientIP(), fmt.Sprint(uid)).StoreSession(context.Background(), session, time.Minute*5)
	if err != nil {
		callback.Error(c, callback.ErrDBOperation, err)
		return
	}

	callback.Success(c, options)
}
