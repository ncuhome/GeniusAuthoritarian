//go:build dev

package global

import (
	"github.com/ncuhome/GeniusAuthoritarian/pkg/ca"
	"github.com/ncuhome/GeniusAuthoritarian/pkg/ed25519"
	"os"
	"path"
	"time"
)

func initKeypair() {
	err := os.MkdirAll(path.Join(ConfigDir, dirJwt), 0600)
	if err != nil {
		panic(err)
	}
	err = os.Mkdir(path.Join(ConfigDir, dirCa), 0600)
	if err != nil {
		panic(err)
	}

	jwtKeypair, err := ed25519.Generate()
	if err != nil {
		panic(err)
	}
	jwtPublic, jwtPrivate, err := jwtKeypair.MarshalPem()
	if err != nil {
		panic(err)
	}
	err = os.WriteFile(path.Join(ConfigDir, dirJwt, _PublicKeyName), jwtPublic, 0600)
	if err != nil {
		panic(err)
	}
	err = os.WriteFile(path.Join(ConfigDir, dirJwt, _PrivateKeyName), jwtPrivate, 0600)
	if err != nil {
		panic(err)
	}

	caCert, caPrivate, err := ca.NewRoot(time.Now().AddDate(10, 0, 0))
	if err != nil {
		panic(err)
	}
	err = os.WriteFile(path.Join(ConfigDir, dirCa, _PublicTlsKeyName), caCert, 0600)
	if err != nil {
		panic(err)
	}
	err = os.WriteFile(path.Join(ConfigDir, dirCa, _PrivateTlsKeyName), caPrivate, 0600)
	if err != nil {
		panic(err)
	}
}
