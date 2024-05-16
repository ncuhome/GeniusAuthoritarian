package redis

import (
	"github.com/ncuhome/GeniusAuthoritarian/internal/pkg/jwt/jwtClaims"
	"github.com/ncuhome/GeniusAuthoritarian/pkg/tokenStore"
)

var NewThirdPartyLogin = tokenStore.NewTokenStoreFactory[jwtClaims.LoginRedis](Client, keyThirdPartyLogin.String())
