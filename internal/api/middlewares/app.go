package middlewares

import (
	"crypto/sha256"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/ncuhome/GeniusAuthoritarian/internal/api/callback"
	"github.com/ncuhome/GeniusAuthoritarian/internal/pkg/jwt"
	"github.com/ncuhome/GeniusAuthoritarian/internal/service"
	"github.com/ncuhome/GeniusAuthoritarian/internal/tools"
	"gorm.io/gorm"
	"sort"
	"strings"
	"time"
)

func RequireAppSignature(c *gin.Context) {
	var header struct {
		AppCode   string `json:"appCode" form:"appCode" binding:"required"`
		TimeStamp int64  `json:"timeStamp" form:"timeStamp" binding:"required"`
		Signature string `json:"signature" form:"signature" binding:"required"`
	}
	if err := c.ShouldBind(&header); err != nil {
		callback.Error(c, callback.ErrForm, err)
		return
	}

	if time.Now().Sub(time.Unix(header.TimeStamp, 0)) > time.Minute*5 {
		callback.Error(c, callback.ErrSignatureExpired)
		return
	}

	_, appSecret, err := service.App.FirstAppKeyPairByAppCode(header.AppCode)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			callback.Error(c, callback.ErrAppCodeNotFound)
			return
		}
		callback.Error(c, callback.ErrDBOperation, err)
		return
	}

	var payload map[string]string
	if err = c.ShouldBindJSON(&payload); err != nil {
		callback.Error(c, callback.ErrForm, err)
		return
	}
	delete(payload, "signature")
	payload["appSecret"] = appSecret

	keys := make([]string, len(payload))
	var signStrLen = len(payload)*2 - 1
	i := 0
	for key, value := range payload {
		keys[i] = key
		signStrLen += len(key) + len(value)
		i++
	}
	sort.Strings(keys)

	signBuilder := strings.Builder{}
	signBuilder.Grow(signStrLen)
	for i, key := range keys {
		if i != 0 {
			signBuilder.Write([]byte("&"))
		}
		signBuilder.Write([]byte(key + "=" + payload[key]))
	}

	h := sha256.New()
	h.Write([]byte(signBuilder.String()))
	if header.Signature != fmt.Sprintf("%x", h.Sum(nil)) {
		callback.Error(c, callback.ErrUnauthorized)
		return
	}

	c.Set("appCode", header.AppCode)
}

// RequireAccessToken 解析并将 access claims 写入上下文
// 需要在 RequireAppSignature 之后调用
func RequireAccessToken(c *gin.Context) {
	var f struct {
		Token string `json:"token" form:"token" binding:"required"`
	}
	if err := c.ShouldBind(&f); err != nil {
		callback.Error(c, callback.ErrForm, err)
		return
	}

	claims, valid, err := jwt.ParseAccessToken(f.Token)
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

	c.Set("access", claims)
}
