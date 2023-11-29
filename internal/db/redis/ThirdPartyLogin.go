package redis

import (
	"github.com/ncuhome/GeniusAuthoritarian/pkg/tokenStore"
)

func NewThirdPartyLogin() tokenStore.TokenStore {
	return tokenStore.NewTokenStore(Client, keyThirdPartyLogin.String())
}
