package redis

import (
	"github.com/ncuhome/GeniusAuthoritarian/internal/pkg/tokenStorePoint"
	"time"
)

func NewThirdPartyLogin(iat time.Time) tokenStorePoint.TokenStore {
	return tokenStorePoint.NewTokenStore(Client, &idThirdPartyLogin, keyThirdPartyLogin.String(), iat)
}
