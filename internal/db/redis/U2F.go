package redis

import (
	"fmt"
	"github.com/ncuhome/GeniusAuthoritarian/internal/pkg/tokenStore"
)

func NewU2F(uid uint) tokenStore.TokenStore {
	return tokenStore.NewTokenStore(Client, keyU2F.String()+fmt.Sprint(uid)+"-")
}
