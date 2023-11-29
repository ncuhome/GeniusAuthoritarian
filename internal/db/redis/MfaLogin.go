package redis

import (
	"fmt"
	"github.com/ncuhome/GeniusAuthoritarian/internal/pkg/jwt/jwtClaims"
	"github.com/ncuhome/GeniusAuthoritarian/pkg/tokenStore"
)

func NewMfaLogin(uid uint) tokenStore.TokenStore[*jwtClaims.MfaRedis] {
	return tokenStore.NewTokenStore[*jwtClaims.MfaRedis](Client, keyUserMfaLogin.String()+fmt.Sprint(uid)+"-")
}
