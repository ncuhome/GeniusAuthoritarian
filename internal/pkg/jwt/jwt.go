package jwt

import (
	"context"
	"github.com/golang-jwt/jwt/v5"
	"github.com/ncuhome/GeniusAuthoritarian/internal/db/redis"
	"github.com/ncuhome/GeniusAuthoritarian/internal/global"
	"time"
)

var key = []byte(global.Config.Jwt.SignKey)

func NewTypedClaims(Type string, valid time.Duration) TypedClaims {
	now := time.Now()
	return TypedClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(valid)),
			IssuedAt:  jwt.NewNumericDate(now),
		},
		Type: Type,
	}
}

func GenerateToken(claims Claims) (string, error) {
	return jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(key)
}

func ParseToken[C Claims](Type, token string, target C) (claims C, valid bool, err error) {
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
	valid = t.Valid && claims.GetType() == Type
	return
}

// GenerateUserToken 生成有效期 15 天的后台 Token
func GenerateUserToken(uid uint, name string, groups []string) (string, *UserToken, error) {
	claims := &UserToken{
		TypedClaims: NewTypedClaims("User", time.Hour*24*15),
		ID:          uid,
		Name:        name,
		Groups:      groups,
	}
	token, err := GenerateToken(claims)
	return token, claims, err
}

func ParseUserToken(token string) (*UserToken, bool, error) {
	return ParseToken("User", token, &UserToken{})
}

// GenerateLoginToken 生成有效期 5 分钟的登录校验 Token
func GenerateLoginToken(claims LoginRedisClaims) (string, error) {
	valid := time.Minute * 5

	tokenClaims := &LoginToken{
		TypedClaims: NewTypedClaims("Login", valid),
	}
	var err error
	tokenClaims.ID, err = redis.NewThirdPartyLogin().CreateStorePoint(context.Background(), tokenClaims.IssuedAt.Time, valid, claims)
	if err != nil {
		return "", err
	}

	return GenerateToken(tokenClaims)
}

// ParseLoginToken 解析后自动销毁
func ParseLoginToken(token string) (*LoginRedisClaims, bool, error) {
	claims, valid, err := ParseToken("Login", token, &LoginToken{})
	if err != nil || !valid {
		return nil, false, err
	}

	var redisClaims LoginRedisClaims
	err = redis.NewThirdPartyLogin().NewStorePoint(claims.ID).GetAndDestroy(context.Background(), claims.IssuedAt.Time, &redisClaims)
	if err != nil {
		if err == redis.Nil {
			err = nil
		}
		return nil, false, err
	}
	return &redisClaims, true, nil
}

// GenerateMfaToken 生成绑定 TOTP MFA 中间身份令牌，五分钟有效
func GenerateMfaToken(claims LoginRedisClaims, mfaSecret, appCallback string) (string, error) {
	valid := time.Minute * 5

	mfaTokenClaims := &MfaToken{
		TypedClaims: NewTypedClaims("Mfa", valid),
		UID:         claims.UID,
	}
	var err error
	if mfaTokenClaims.ID, err = redis.NewMfaLogin(claims.UID).CreateStorePoint(context.Background(), mfaTokenClaims.IssuedAt.Time, valid, &MfaRedisClaims{
		LoginRedisClaims: claims,
		Mfa:              mfaSecret,
		AppCallback:      appCallback,
	}); err != nil {
		return "", err
	}

	return GenerateToken(mfaTokenClaims)
}

// ParseMfaToken 不会销毁，允许多次验证尝试
func ParseMfaToken(token string) (*MfaRedisClaims, error) {
	claims, valid, err := ParseToken("Mfa", token, &MfaToken{})
	if err != nil || !valid {
		return nil, err
	}

	var redisClaims MfaRedisClaims
	return &redisClaims, redis.NewMfaLogin(claims.UID).NewStorePoint(claims.ID).Get(context.Background(), claims.IssuedAt.Time, &redisClaims)
}

// GenerateU2fToken 生成后台 U2F 身份令牌，五分钟有效
func GenerateU2fToken(uid uint, ip string) (string, *U2fToken, error) {
	valid := time.Minute * 5

	tokenClaims := &U2fToken{
		TypedClaims: NewTypedClaims("U2F", valid),
		UID:         uid,
		IP:          ip,
	}
	var err error
	if tokenClaims.ID, err = redis.NewU2F(uid).CreateStorePoint(context.Background(), tokenClaims.IssuedAt.Time, valid, nil); err != nil {
		return "", nil, err
	}

	token, err := GenerateToken(tokenClaims)
	return token, tokenClaims, err
}

func ParseU2fToken(token, ip string) (bool, error) {
	claims, valid, err := ParseToken("U2F", token, &U2fToken{})
	if err != nil || !valid || claims.IP != ip {
		return false, err
	}

	err = redis.NewU2F(claims.UID).NewStorePoint(claims.ID).Get(context.Background(), claims.IssuedAt.Time, nil)
	if err != nil {
		if err == redis.Nil {
			err = nil
		}
		return false, err
	}
	return true, nil
}
