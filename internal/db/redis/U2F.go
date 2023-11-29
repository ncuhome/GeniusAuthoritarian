package redis

import (
	"fmt"
	"github.com/ncuhome/GeniusAuthoritarian/pkg/tokenStore"
	"go/types"
)

func NewU2F(uid uint) tokenStore.TokenStore[types.Nil] {
	return tokenStore.NewTokenStore[types.Nil](Client, keyU2F.String()+fmt.Sprint(uid)+"-")
}
