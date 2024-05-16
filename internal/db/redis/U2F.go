package redis

import (
	"github.com/ncuhome/GeniusAuthoritarian/pkg/tokenStore"
	"go/types"
)

var NewU2F = tokenStore.NewTokenStoreFactory[types.Nil](Client, keyU2F.String())
