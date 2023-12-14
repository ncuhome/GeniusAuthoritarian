package public

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/ncuhome/GeniusAuthoritarian/internal/api/callback"
	"github.com/ncuhome/GeniusAuthoritarian/internal/api/models/response"
	"github.com/ncuhome/GeniusAuthoritarian/internal/db/redis"
	"github.com/ncuhome/GeniusAuthoritarian/internal/pkg/jwt"
	"github.com/ncuhome/GeniusAuthoritarian/internal/service"
	"github.com/ncuhome/GeniusAuthoritarian/internal/tools"
)

func RefreshToken(c *gin.Context) {
	var f struct {
		Token string `json:"token" binding:"required"`
	}
	if err := tools.ShouldBindReused(c, &f); err != nil {
		callback.Error(c, callback.ErrForm, err)
		return
	}

	claims, valid, err := jwt.ParseRefreshToken(f.Token)
	if err != nil {
		callback.Error(c, callback.ErrTokenInvalid, err)
		return
	} else if !valid {
		callback.Error(c, callback.ErrTokenInvalid)
		return
	}

	appCode := tools.GetAppCode(c)
	if appCode != claims.AppCode {
		callback.Error(c, callback.ErrUnauthorized)
		return
	}

	accessToken, err := jwt.GenerateAccessToken(claims.ID, claims.UID, appCode, claims.Payload)
	if err != nil {
		callback.Error(c, callback.ErrUnexpected, err)
		return
	}

	callback.Success(c, &response.RefreshToken{
		AccessToken: accessToken,
		Payload:     claims.Payload,
	})
}

func DestroyRefreshToken(c *gin.Context) {
	var f struct {
		Token string `json:"token" binding:"required"`
	}
	if err := tools.ShouldBindReused(c, &f); err != nil {
		callback.Error(c, callback.ErrForm, err)
		return
	}

	claims, valid, err := jwt.ParseRefreshToken(f.Token)
	if err != nil {
		callback.Error(c, callback.ErrTokenInvalid, err)
		return
	} else if !valid {
		callback.Error(c, callback.ErrTokenInvalid)
		return
	}

	err = redis.NewRecordedToken().NewStorePoint(claims.ID).Destroy(context.Background())
	if err != nil {
		callback.Error(c, callback.ErrUnexpected, err)
		return
	}

	callback.Default(c)
}

func VerifyAccessToken(c *gin.Context) {
	claims := tools.GetAccessClaims(c)
	callback.Success(c, &response.VerifyAccessToken{
		UID:     claims.UID,
		Payload: claims.Payload,
	})
}

func GetUserInfo(c *gin.Context) {
	claims := tools.GetAccessClaims(c)
	user, err := service.User.UserInfoByID(claims.UID)
	if err != nil {
		callback.Error(c, callback.ErrDBOperation, err)
		return
	}

	groups, err := service.UserGroups.GetNamesForUser(claims.UID)
	if err != nil {
		callback.Error(c, callback.ErrDBOperation, err)
		return
	}

	callback.Success(c, gin.H{
		"userID":    user.ID,
		"name":      user.Name,
		"groups":    groups,
		"avatarUrl": user.AvatarUrl,
	})
}
