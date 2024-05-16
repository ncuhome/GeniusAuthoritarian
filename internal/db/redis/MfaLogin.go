package redis

import (
	"github.com/go-redis/redis/v8"
	"github.com/ncuhome/GeniusAuthoritarian/internal/pkg/jwt/jwtClaims"
	"github.com/ncuhome/GeniusAuthoritarian/pkg/tokenStore"
)

var NewMfaLogin = tokenStore.NewTokenStoreFactory[jwtClaims.MfaRedis](keyUserMfaLogin.String(), func() *redis.Client {
	return Client
})
