package redis

import (
	"github.com/ncuhome/GeniusAuthoritarian/pkg/tokenStore"
	"go/types"
)

func NewRefreshToken() tokenStore.TokenStore[types.Nil] {
	return tokenStore.NewTokenStore[types.Nil](Client, keyRefreshToken.String())
}
