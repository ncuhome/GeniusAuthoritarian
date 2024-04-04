package ed25519

import (
	"crypto/ed25519"
	"crypto/rand"
	"github.com/ncuhome/GeniusAuthoritarian/pkg/keypair"
)

func Generate() (*KeyPair, error) {
	publicKey, privateKey, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		return nil, err
	}

	return &KeyPair{
		Public:  publicKey,
		Private: privateKey,
	}, nil
}

type KeyPair struct {
	Public  ed25519.PublicKey
	Private ed25519.PrivateKey
}

func (a KeyPair) MarshalPem() (public, private []byte, err error) {
	publicPem, err := keypair.PemMarshalPublic(keypair.FormatECDSA, a.Public)
	if err != nil {
		return nil, nil, err
	}
	privatePem, err := keypair.PemMarshalPrivate(keypair.FormatECDSA, a.Private)
	if err != nil {
		return nil, nil, err
	}
	return publicPem, privatePem, nil
}

func (a KeyPair) MarshalSSH(comment string) (public, private []byte, err error) {
	publicSshKey, err := keypair.SshMarshalPublic(a.Public)
	if err != nil {
		return nil, nil, err
	}
	privatePem, err := keypair.SshMarshalPrivate(a.Private, comment)
	if err != nil {
		return nil, nil, err
	}
	return publicSshKey, privatePem, nil
}
