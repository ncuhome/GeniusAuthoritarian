package redis

import (
	"github.com/go-redis/redis/v8"
	"github.com/ncuhome/GeniusAuthoritarian/pkg/tokenStore"
	"go/types"
)

var NewU2F = tokenStore.NewTokenStoreFactory[types.Nil](keyU2F.String(), func() *redis.Client {
	return Client
})
