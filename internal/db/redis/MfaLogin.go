package redis

import (
	"fmt"
	"github.com/ncuhome/GeniusAuthoritarian/internal/pkg/tokenStore"
)

func NewMfaLogin(uid uint) tokenStore.TokenStore {
	return tokenStore.NewTokenStore(Client, keyUserMfaLogin.String()+fmt.Sprint(uid)+"-")
}
