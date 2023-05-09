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

// GenerateUserToken 生成有效期 15 天的个人信息访问 Token
func GenerateUserToken(uid uint) (string, error) {
	return GenerateToken(&UserToken{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24 * 15)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
		ID: uid,
	})
}

// GenerateLoginToken 生成有效期 5 分钟的登录校验 Token
func GenerateLoginToken(uid uint, appCode, userMame, ip string) (string, error) {
	now := time.Now()
	valid := time.Minute * 5
	id, e := redis.Jwt.NewLoginPoint(now.Unix(), valid, LoginTokenClaims{
		UID:     uid,
		IP:      ip,
		Name:    userMame,
		AppCode: appCode,
	})
	if e != nil {
		return "", e
	}
	return GenerateToken(&LoginToken{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(valid)),
			IssuedAt:  jwt.NewNumericDate(now),
		},
		ID: id,
	})
}

func ParseUserToken(token string) (*UserToken, bool, error) {
	return ParseToken(token, &UserToken{})
}

// ParseLoginToken 解析后自动销毁
func ParseLoginToken(token string) (*LoginTokenClaims, bool, error) {
	claims, valid, e := ParseToken(token, &LoginToken{})
	if e != nil || !valid {
		return nil, false, e
	}

	var redisClaims LoginTokenClaims
	valid, e = redis.Jwt.VerifyLoginPoint(claims.ID, claims.IssuedAt.Unix(), &redisClaims)
	if e != nil {
		if e == redis.Nil {
			e = nil
		}
		return nil, false, e
	}
	return &redisClaims, valid, DestroyAuthToken(claims.ID)
}

func DestroyAuthToken(cID uint64) error {
	return redis.Jwt.DelLoginPoint(cID)
}
