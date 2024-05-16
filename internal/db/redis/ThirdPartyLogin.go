package redis

import (
	"github.com/go-redis/redis/v8"
	"github.com/ncuhome/GeniusAuthoritarian/internal/pkg/jwt/jwtClaims"
	"github.com/ncuhome/GeniusAuthoritarian/pkg/tokenStore"
)

var NewThirdPartyLogin = tokenStore.NewTokenStoreFactory[jwtClaims.LoginRedis](keyThirdPartyLogin.String(), func() *redis.Client {
	return Client
})
