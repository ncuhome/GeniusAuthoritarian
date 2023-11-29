package redis

import (
	"fmt"
	"github.com/ncuhome/GeniusAuthoritarian/pkg/tokenStore"
)

func NewU2F(uid uint) tokenStore.TokenStore[interface{}] {
	return tokenStore.NewTokenStore[interface{}](Client, keyU2F.String()+fmt.Sprint(uid)+"-")
}
