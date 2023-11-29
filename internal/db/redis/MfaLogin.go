package redis

import (
	"fmt"
	"github.com/ncuhome/GeniusAuthoritarian/pkg/tokenStore"
)

func NewMfaLogin(uid uint) tokenStore.TokenStore {
	return tokenStore.NewTokenStore(Client, keyUserMfaLogin.String()+fmt.Sprint(uid)+"-")
}
