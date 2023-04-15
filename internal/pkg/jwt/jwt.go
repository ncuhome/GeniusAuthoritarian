package jwt

import (
	"github.com/golang-jwt/jwt/v4"
	"github.com/ncuhome/GeniusAuthoritarian/internal/db/redis"
	"github.com/ncuhome/GeniusAuthoritarian/internal/global"
	"time"
)

var key = []byte(global.Config.Jwt.SignKey)

func GenerateToken(claims jwt.Claims) (string, error) {
	return jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(key)
}

func ParseToken[C jwt.Claims](token string, target C) (claims C, valid bool, e error) {
	var t *jwt.Token
	t, e = jwt.ParseWithClaims(token, target, func(t *jwt.Token) (interface{}, error) {
		return key, nil
	})
	if e != nil {
		return
	}

	claims, _ = t.Claims.(C)
	valid = t.Valid
	return
}

// GenerateRefreshToken 生成有效期 15 天的 Refresh Token
func GenerateRefreshToken(name string, groups []string) (string, error) {
	return GenerateToken(&RefreshToken{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24 * 15)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
		Name:   name,
		Groups: groups,
	})
}

// GenerateAuthToken 生成有效期 5 分钟的校验 Token
func GenerateAuthToken(name string, groups []string) (string, error) {
	now := time.Now()
	valid := time.Minute * 5
	id, e := redis.Jwt.NewAuthPoint(now.Unix(), valid)
	if e != nil {
		return "", e
	}
	return GenerateToken(&AuthToken{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(valid)),
			IssuedAt:  jwt.NewNumericDate(now),
		},
		ID:     id,
		Name:   name,
		Groups: groups,
	})
}

func ParseRefreshToken(token string) (*RefreshToken, bool, error) {
	return ParseToken(token, &RefreshToken{})
}

// ParseAuth 解析后自动销毁
func ParseAuth(token string) (*AuthToken, bool, error) {
	claims, valid, e := ParseToken(token, &AuthToken{})
	if e != nil || !valid {
		return claims, false, e
	}

	valid, e = redis.Jwt.VerifyAuthPoint(claims.ID, claims.IssuedAt.Unix())
	if e != nil {
		if e == redis.Nil {
			return claims, false, nil
		}
		return claims, false, e
	}
	return claims, valid, DestroyAuthToken(claims.ID)
}

func DestroyAuthToken(cID uint64) error {
	return redis.Jwt.DelAuthPoint(cID)
}
