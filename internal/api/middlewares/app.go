package middlewares

import (
	"crypto/sha256"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/ncuhome/GeniusAuthoritarian/internal/api/callback"
	"github.com/ncuhome/GeniusAuthoritarian/internal/service"
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

	var payload map[string]interface{}
	if err := c.ShouldBind(&payload); err != nil {
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
		valueStr := fmt.Sprint(value)
		payload[key] = valueStr
		signStrLen += len(key) + len(valueStr)
		i++
	}
	sort.Strings(keys)

	signBuilder := strings.Builder{}
	signBuilder.Grow(signStrLen)
	for i, key := range keys {
		if i != 0 {
			signBuilder.Write([]byte("&"))
		}
		signBuilder.Write([]byte(key + "=" + payload[key].(string)))
	}

	h := sha256.New()
	h.Write([]byte(signBuilder.String()))
	if header.Signature != fmt.Sprintf("%x", h.Sum(nil)) {
		callback.Error(c, callback.ErrUnauthorized)
		return
	}
}
