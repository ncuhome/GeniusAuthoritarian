package redis

import (
	"github.com/ncuhome/GeniusAuthoritarian/internal/pkg/jwt/jwtClaims"
	"github.com/ncuhome/GeniusAuthoritarian/pkg/tokenStore"
)

func NewMfaLogin() tokenStore.TokenStore[jwtClaims.MfaRedis] {
	return tokenStore.NewTokenStore[jwtClaims.MfaRedis](Client, keyUserMfaLogin.String())
}
