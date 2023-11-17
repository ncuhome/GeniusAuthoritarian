package redis

import (
	"github.com/ncuhome/GeniusAuthoritarian/internal/pkg/tokenStore"
)

func NewThirdPartyLogin() tokenStore.TokenStore {
	return tokenStore.NewTokenStore(Client, keyThirdPartyLogin.String())
}
