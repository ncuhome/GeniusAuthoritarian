package controllers

import (
	"context"
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
	user, err := webAuthn.NewUser(tools.GetUserInfo(c).UID)
	if err != nil {
		callback.Error(c, callback.ErrDBOperation, err)
		return
	}

	options, session, err := webAuthn.Client.BeginLogin(user)
	if err != nil {
		callback.Error(c, callback.ErrUnexpected, err)
		return
	}

	err = redis.NewPasskey(c.ClientIP()).StoreSession(context.Background(), session, time.Minute*5)
	if err != nil {
		callback.Error(c, callback.ErrDBOperation, err)
		return
	}

	callback.Success(c, options)
}
