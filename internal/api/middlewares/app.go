package middlewares

import (
	"bytes"
	"crypto/sha256"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/ncuhome/GeniusAuthoritarian/internal/api/callback"
	"github.com/ncuhome/GeniusAuthoritarian/internal/pkg/jwt"
	"github.com/ncuhome/GeniusAuthoritarian/internal/service"
	"github.com/ncuhome/GeniusAuthoritarian/internal/tools"
	"gorm.io/gorm"
	"io"
	"sort"
	"strings"
	"time"
	"unsafe"
)

func RequireAppSignature(c *gin.Context) {
	var header struct {
		AppCode   string `json:"appCode" binding:"required"`
		TimeStamp int64  `json:"timeStamp" binding:"required"`
		Signature string `json:"signature" binding:"required"`
	}
	var form map[string]interface{}
	var err error
	if c.Request.Method != "GET" && c.Request.Method != "DELETE" {
		if c.Request.Body != nil {
			var bodyBytes []byte
			bodyBytes, err = io.ReadAll(c.Request.Body)
			if err != nil {
				callback.Error(c, callback.ErrForm, err)
				return
			}

			if err = binding.JSON.BindBody(bodyBytes, &header); err != nil {
				callback.Error(c, callback.ErrForm, err)
				return
			}

			jsonDecoder := json.NewDecoder(bytes.NewReader(bodyBytes))
			jsonDecoder.UseNumber()
			if err = jsonDecoder.Decode(&form); err != nil {
				callback.Error(c, callback.ErrForm, err)
				return
			}

			c.Set(gin.BodyBytesKey, bodyBytes)
		} else {
			callback.Error(c, callback.ErrForm, "signature required")
			return
		}
	} else {
		if err = c.ShouldBindQuery(&header); err != nil {
			callback.Error(c, callback.ErrForm, err)
			return
		}

		query := c.Request.URL.Query()
		form = make(map[string]interface{}, len(query))
		for k, v := range query {
			form[k] = v[0]
		}
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

	delete(form, "signature")
	formStrMap := make(map[string]string, len(form)+1)
	for key, value := range form {
		formStrMap[key] = fmt.Sprint(value)
	}
	formStrMap["appSecret"] = appSecret

	keys := make([]string, len(formStrMap))
	var signStrLen = len(formStrMap)*2 - 1
	i := 0
	for key, value := range formStrMap {
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
		signBuilder.Write([]byte(key + "=" + formStrMap[key]))
	}

	h := sha256.New()
	signStr := signBuilder.String()
	h.Write(unsafe.Slice(unsafe.StringData(signStr), len(signStr)))
	if header.Signature != fmt.Sprintf("%x", h.Sum(nil)) {
		callback.Error(c, callback.ErrUnauthorized, "signature invalid")
		return
	}

	c.Set("appCode", header.AppCode)
}

// RequireAccessToken 解析并将 access claims 写入上下文
// 需要在 RequireAppSignature 之后调用
func RequireAccessToken(c *gin.Context) {
	var f struct {
		Token string `json:"token" binding:"required"`
	}
	if err := tools.ShouldBindReused(c, &f); err != nil {
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
