package jwt

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/ncuhome/GeniusAuthoritarian/internal/db/redis"
	"github.com/ncuhome/GeniusAuthoritarian/internal/global"
	"time"
)

var key = []byte(global.Config.Jwt.SignKey)

func GenerateToken(claims jwt.Claims) (string, error) {
	return jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(key)
}

func ParseToken[C jwt.Claims](token string, target C) (claims C, valid bool, err error) {
	var t *jwt.Token
	t, err = jwt.ParseWithClaims(
		token, target, func(t *jwt.Token) (interface{}, error) {
			return key, nil
		},
		jwt.WithLeeway(time.Second*3),
	)
	if err != nil {
		return
	}

	claims, _ = t.Claims.(C)
	valid = t.Valid
	return
}

// GenerateUserToken 生成有效期 15 天的后台 Token
func GenerateUserToken(uid uint, name string, groups []string) (string, error) {
	now := time.Now()
	return GenerateToken(&UserToken{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(time.Hour * 24 * 15)),
			IssuedAt:  jwt.NewNumericDate(now),
		},
		ID:     uid,
		Name:   name,
		Groups: groups,
	})
}

// GenerateLoginToken 生成有效期 5 分钟的登录校验 Token
func GenerateLoginToken(clams LoginTokenClaims) (string, error) {
	now := time.Now()
	valid := time.Minute * 5

	id, err := redis.ThirdPartyLogin.NewLoginPoint(now.Unix(), valid, clams)
	if err != nil {
		return "", err
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
func GenerateMfaToken(clams LoginTokenClaims, mfaSecret, appCallback string) (string, error) {
	now := time.Now()
	valid := time.Minute * 5

	token, err := GenerateToken(&MfaToken{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(valid)),
			IssuedAt:  jwt.NewNumericDate(now),
		},
		UID: clams.UID,
	})
	if err != nil {
		return "", err
	}

	if err = redis.MfaLogin.Set(clams.UID, token, valid, MfaLoginClaims{
		LoginTokenClaims: clams,
		Mfa:              mfaSecret,
		AppCallback:      appCallback,
	}); err != nil {
		return "", err
	}

	return token, nil
}

func ParseUserToken(token string) (*UserToken, bool, error) {
	return ParseToken(token, &UserToken{})
}

// ParseLoginToken 解析后自动销毁
func ParseLoginToken(token string) (*LoginTokenClaims, bool, error) {
	claims, valid, err := ParseToken(token, &LoginToken{})
	if err != nil || !valid {
		return nil, false, err
	}

	var redisClaims LoginTokenClaims
	valid, err = redis.ThirdPartyLogin.VerifyLoginPoint(claims.ID, claims.IssuedAt.Unix(), &redisClaims)
	if err != nil {
		if err == redis.Nil {
			err = nil
		}
		return nil, false, err
	}
	return &redisClaims, valid, nil
}

// ParseMfaToken 不会销毁，允许多次验证尝试
func ParseMfaToken(token string) (*MfaLoginClaims, error) {
	claims, valid, err := ParseToken(token, &MfaToken{})
	if err != nil || !valid {
		return nil, err
	}

	var redisClaims MfaLoginClaims
	return &redisClaims, redis.MfaLogin.Get(claims.UID, token, &redisClaims)
}
