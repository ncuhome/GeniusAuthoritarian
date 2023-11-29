package redis

import (
	"github.com/ncuhome/GeniusAuthoritarian/internal/pkg/jwt/jwtClaims"
	"github.com/ncuhome/GeniusAuthoritarian/pkg/tokenStore"
)

func NewThirdPartyLogin() tokenStore.TokenStore[jwtClaims.LoginRedis] {
	return tokenStore.NewTokenStore[jwtClaims.LoginRedis](Client, keyThirdPartyLogin.String())
}
