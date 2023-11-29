package redis

import (
	"fmt"
	"github.com/ncuhome/GeniusAuthoritarian/pkg/tokenStore"
)

func NewU2F(uid uint) tokenStore.TokenStore {
	return tokenStore.NewTokenStore(Client, keyU2F.String()+fmt.Sprint(uid)+"-")
}
