package ed25519

import (
	"crypto/ed25519"
	"github.com/ncuhome/GeniusAuthoritarian/pkg/keypair"
	"os"
	"path"
)

type Loader struct {
	Dir string
}

func (loader Loader) LoadPem(filename string) ([]byte, error) {
	keyPath := path.Join(loader.Dir, filename)
	return os.ReadFile(keyPath)
}

func (loader Loader) LoadPrivateKey(filename string) (ed25519.PrivateKey, error) {
	keyBytes, err := loader.LoadPem(filename)
	if err != nil {
		return nil, err
	}
	return keypair.PemUnmarshalPrivate[ed25519.PrivateKey](keyBytes)
}

func (loader Loader) LoadPublicKey(filename string) (ed25519.PublicKey, error) {
	keyBytes, err := loader.LoadPem(filename)
	if err != nil {
		return nil, err
	}
	return keypair.PemUnmarshalPublic[ed25519.PublicKey](keyBytes)
}
