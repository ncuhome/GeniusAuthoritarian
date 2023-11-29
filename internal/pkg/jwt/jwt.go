package jwt

import (
	"context"
	"fmt"
	"github.com/Mmx233/daoUtil"
	"github.com/golang-jwt/jwt/v5"
	"github.com/ncuhome/GeniusAuthoritarian/internal/db/redis"
	"github.com/ncuhome/GeniusAuthoritarian/internal/global"
	"github.com/ncuhome/GeniusAuthoritarian/internal/pkg/jwt/jwtClaims"
	"github.com/ncuhome/GeniusAuthoritarian/internal/pkg/jwt/jwtVerify"
	"github.com/ncuhome/GeniusAuthoritarian/internal/service"
	"time"
)

var key = []byte(global.Config.Jwt.SignKey)

const (
	User    = "User"
	Login   = "Login"
	Mfa     = "Mfa"
	U2F     = "U2F"
	Refresh = "Refresh"
	Access  = "Access"
)

func NewTypedClaims(Type string, valid time.Duration) jwtClaims.TypedClaims {
	now := time.Now()
	return jwtClaims.TypedClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(valid)),
			IssuedAt:  jwt.NewNumericDate(now),
		},
		Type: Type,
	}
}

func NewUserClaims(uid uint, Type string, valid time.Duration) (claims jwtClaims.UserClaims, err error) {
	redisOperator := redis.NewUserJwt().NewOperator(uid)
	oid, err := redisOperator.GetOperateID(context.Background())
	if err == redis.Nil {
		var userSrv service.UserSrv
		userSrv, err = service.User.Begin()
		if err != nil {
			return
		}
		defer userSrv.Rollback()

		// service 操作会写入 redis，不用再操作创建 hash 项
		var exist bool
		exist, err = userSrv.UserIdExist(uid, daoUtil.LockForShare)
		if err != nil {
			return
		} else if !exist {
			err = fmt.Errorf("user %d not exist", uid)
			return
		}

		oid, err = redisOperator.GetOperateID(context.Background())
		if err != nil {
			return
		}
	} else if err != nil {
		return
	}

	return jwtClaims.UserClaims{
		TypedClaims:   NewTypedClaims(Type, valid),
		UID:           uid,
		UserOperateID: oid,
	}, nil
}

func GenerateToken(claims jwtClaims.Claims) (string, error) {
	return jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(key)
}

func ParseToken[C jwtClaims.Claims](Type, token string, target C) (claims C, valid bool, err error) {
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

func ParseTokenAndVerify[C jwtClaims.ClaimsUser](Type, token string, target C) (claims C, valid bool, err error) {
	claims, valid, err = ParseToken(Type, token, target)
	if err == nil && valid {
		valid, err = jwtVerify.CheckUserClaims(context.Background(), claims)
	}
	return
}

// GenerateUserToken 生成后台 Token
func GenerateUserToken(uid uint, name string, groups []string, valid time.Duration) (string, error) {
	userClaims, err := NewUserClaims(uid, User, valid)
	if err != nil {
		return "", err
	}
	claims := &jwtClaims.UserToken{
		UserClaims: userClaims,
		Name:       name,
		Groups:     groups,
	}
	token, err := GenerateToken(claims)
	return token, err
}

func ParseUserToken(token string) (*jwtClaims.UserToken, bool, error) {
	return ParseTokenAndVerify(User, token, &jwtClaims.UserToken{})
}

// GenerateLoginToken 生成有效期 5 分钟的登录校验 Token
func GenerateLoginToken(claims jwtClaims.LoginRedis) (string, error) {
	valid := time.Minute * 5

	userClaims, err := NewUserClaims(claims.UID, Login, valid)
	if err != nil {
		return "", err
	}

	tokenClaims := &jwtClaims.LoginToken{
		UserClaims: userClaims,
	}
	tokenClaims.ID, err = redis.NewThirdPartyLogin().CreateStorePoint(context.Background(), valid, &claims)
	if err != nil {
		return "", err
	}

	return GenerateToken(tokenClaims)
}

// ParseLoginToken 解析后自动销毁
func ParseLoginToken(token string) (*jwtClaims.LoginRedis, bool, error) {
	claims, valid, err := ParseTokenAndVerify(Login, token, &jwtClaims.LoginToken{})
	if err != nil || !valid {
		return nil, false, err
	}

	var redisClaims jwtClaims.LoginRedis
	err = redis.NewThirdPartyLogin().NewStorePoint(claims.ID).GetAndDestroy(context.Background(), &redisClaims)
	if err != nil {
		if err == redis.Nil {
			err = nil
		}
		return nil, false, err
	}
	return &redisClaims, true, nil
}

// GenerateMfaToken 生成绑定 TOTP MFA 中间身份令牌，五分钟有效
func GenerateMfaToken(claims jwtClaims.LoginRedis, mfaSecret, appCallback string) (string, error) {
	valid := time.Minute * 5

	userClaims, err := NewUserClaims(claims.UID, Mfa, valid)
	if err != nil {
		return "", err
	}
	mfaTokenClaims := &jwtClaims.MfaToken{
		UserClaims: userClaims,
	}
	if mfaTokenClaims.ID, err = redis.NewMfaLogin().CreateStorePoint(context.Background(), valid, &jwtClaims.MfaRedis{
		LoginRedis:  claims,
		Mfa:         mfaSecret,
		AppCallback: appCallback,
	}); err != nil {
		return "", err
	}

	return GenerateToken(mfaTokenClaims)
}

// ParseMfaToken 不会销毁，允许多次验证尝试
func ParseMfaToken(token string) (*jwtClaims.MfaRedis, bool, error) {
	claims, valid, err := ParseTokenAndVerify(Mfa, token, &jwtClaims.MfaToken{})
	if err != nil || !valid {
		return nil, valid, err
	}

	var redisClaims jwtClaims.MfaRedis
	return &redisClaims, true, redis.NewMfaLogin().NewStorePoint(claims.ID).Get(context.Background(), &redisClaims)
}

// GenerateU2fToken 生成后台 U2F 身份令牌，五分钟有效
func GenerateU2fToken(uid uint, ip string) (string, *jwtClaims.U2fToken, error) {
	valid := time.Minute * 5

	userClaims, err := NewUserClaims(uid, U2F, valid)
	if err != nil {
		return "", nil, err
	}
	tokenClaims := &jwtClaims.U2fToken{
		UserClaims: userClaims,
		IP:         ip,
	}
	if tokenClaims.ID, err = redis.NewU2F().CreateStorePoint(context.Background(), valid, nil); err != nil {
		return "", nil, err
	}

	token, err := GenerateToken(tokenClaims)
	return token, tokenClaims, err
}

func ParseU2fToken(token, ip string) (bool, error) {
	claims, valid, err := ParseTokenAndVerify(U2F, token, &jwtClaims.U2fToken{})
	if err != nil || !valid || claims.IP != ip {
		return false, err
	}

	err = redis.NewU2F().NewStorePoint(claims.ID).Get(context.Background(), nil)
	if err != nil {
		if err == redis.Nil {
			err = nil
		}
		return false, err
	}
	return true, nil
}

func GenerateRefreshToken(uid uint, appCode, payload string, valid time.Duration) (string, error) {
	userClaims, err := NewUserClaims(uid, Refresh, valid)
	if err != nil {
		return "", err
	}
	claims := jwtClaims.RefreshToken{
		UserClaims: userClaims,
		AppCode:    appCode,
		Payload:    payload,
	}
	claims.ID, err = redis.NewRefreshToken().CreateStorePoint(context.Background(), valid, nil)
	if err != nil {
		return "", err
	}
	return GenerateToken(&claims)
}

func ParseRefreshToken(token string) (*jwtClaims.RefreshToken, bool, error) {
	claims, valid, err := ParseTokenAndVerify(Refresh, token, &jwtClaims.RefreshToken{})
	if err != nil || !valid {
		return claims, valid, err
	}
	return claims, true, redis.NewRefreshToken().NewStorePoint(claims.ID).Get(context.Background(), nil)
}

func GenerateAccessToken(uid uint, appCode, payload string) (string, error) {
	valid := time.Minute * 5

	userClaims, err := NewUserClaims(uid, Access, valid)
	if err != nil {
		return "", err
	}
	return GenerateToken(&jwtClaims.AccessToken{
		RefreshToken: jwtClaims.RefreshToken{
			UserClaims: userClaims,
			AppCode:    appCode,
			Payload:    payload,
		},
	})
}

func ParseAccessToken(token string) (*jwtClaims.AccessToken, bool, error) {
	return ParseTokenAndVerify(Access, token, &jwtClaims.AccessToken{})
}
