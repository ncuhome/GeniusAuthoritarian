package public

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/ncuhome/GeniusAuthoritarian/internal/api/callback"
	"github.com/ncuhome/GeniusAuthoritarian/internal/api/models/response"
	"github.com/ncuhome/GeniusAuthoritarian/internal/pkg/jwt"
	"github.com/ncuhome/GeniusAuthoritarian/internal/tools"
)

func RefreshToken(c *gin.Context) {
	var f struct {
		Token string `json:"token" binding:"required"`
	}
	if err := c.ShouldBindBodyWith(&f, binding.JSON); err != nil {
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

	accessToken, err := jwt.GenerateAccessToken(claims.UID, appCode, claims.Payload)
	if err != nil {
		callback.Error(c, callback.ErrUnexpected, err)
		return
	}

	callback.Success(c, &response.RefreshToken{
		AccessToken: accessToken,
		Payload:     claims.Payload,
	})
}

func VerifyAccessToken(c *gin.Context) {
	claims := tools.GetAccessClaims(c)
	callback.Success(c, &response.VerifyAccessToken{
		UID:     claims.UID,
		Payload: claims.Payload,
	})
}
