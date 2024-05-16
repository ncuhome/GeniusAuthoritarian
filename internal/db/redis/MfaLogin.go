package redis

import (
	"github.com/ncuhome/GeniusAuthoritarian/internal/pkg/jwt/jwtClaims"
	"github.com/ncuhome/GeniusAuthoritarian/pkg/tokenStore"
)

var NewMfaLogin = tokenStore.NewTokenStoreFactory[jwtClaims.MfaRedis](Client, keyUserMfaLogin.String())
