package redis

import (
	"github.com/ncuhome/GeniusAuthoritarian/pkg/tokenStore"
	"go/types"
)

func NewRecordedToken() tokenStore.TokenStore[types.Nil] {
	return tokenStore.NewTokenStore[types.Nil](Client, keyRecordedToken.String())
}
