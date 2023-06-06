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

// GenerateUserToken 生成有效期 15 天的后台 Token
func GenerateUserToken(uid uint, name string, groups []string) (string, error) {
	return GenerateToken(&UserToken{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24 * 15)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
		ID:     uid,
		Name:   name,
		Groups: groups,
	})
}

// GenerateLoginToken 生成有效期 5 分钟的登录校验 Token
func GenerateLoginToken(uid, appID uint, name, ip string, groups []string) (string, error) {
	now := time.Now()
	valid := time.Minute * 5

	id, e := redis.ThirdPartyLogin.NewLoginPoint(now.Unix(), valid, LoginTokenClaims{
		UID:    uid,
		IP:     ip,
		Name:   name,
		AppID:  appID,
		Groups: groups,
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

// GenerateMfaToken 生成 2FA 中间身份令牌，五分钟有效
func GenerateMfaToken(uid, appID uint, name, ip string, groups []string) (string, error) {
	now := time.Now()
	valid := time.Minute * 5

	token, e := GenerateToken(&MfaToken{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(valid)),
			IssuedAt:  jwt.NewNumericDate(now),
		},
		UID: uid,
	})
	if e != nil {
		return "", e
	}

	if e = redis.MfaLogin.Set(uid, token, valid, MfaLoginClaims{
		LoginTokenClaims: LoginTokenClaims{
			UID:    uid,
			Name:   name,
			IP:     ip,
			Groups: groups,
			AppID:  appID,
		},
	}); e != nil {
		return "", e
	}

	return token, nil
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
	valid, e = redis.ThirdPartyLogin.VerifyLoginPoint(claims.ID, claims.IssuedAt.Unix(), &redisClaims)
	if e != nil {
		if e == redis.Nil {
			e = nil
		}
		return nil, false, e
	}
	return &redisClaims, valid, DestroyAuthToken(claims.ID)
}

// ParseMfaToken 不会销毁，允许多次验证尝试
func ParseMfaToken(token string) (*MfaLoginClaims, error) {
	claims, valid, e := ParseToken(token, &MfaToken{})
	if e != nil || !valid {
		return nil, e
	}

	var redisClaims MfaLoginClaims
	return &redisClaims, redis.MfaLogin.Get(claims.UID, token, &redisClaims)
}

func DestroyAuthToken(cID uint64) error {
	return redis.ThirdPartyLogin.DelLoginPoint(cID)
}
