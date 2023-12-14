package controllers

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/ncuhome/GeniusAuthoritarian/internal/api/callback"
	"github.com/ncuhome/GeniusAuthoritarian/internal/db/redis"
	"github.com/ncuhome/GeniusAuthoritarian/internal/tools"
)

func Logout(c *gin.Context) {
	loginID := tools.GetUserInfo(c).ID
	err := redis.NewRecordedToken().NewStorePoint(loginID).Destroy(context.Background())
	if err != nil {
		callback.Error(c, callback.ErrUnexpected, err)
		return
	}

	callback.Default(c)
}
