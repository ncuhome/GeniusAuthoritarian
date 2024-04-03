package ed25519

import (
	"crypto/ed25519"
	"github.com/ncuhome/GeniusAuthoritarian/pkg/keypair"
	"math/rand"
)

func Generate(randRand *rand.Rand) (*KeyPair, error) {
	publicKey, privateKey, err := ed25519.GenerateKey(rand.New(randRand))
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

func (a KeyPair) MarshalPem() (public []byte, private []byte, err error) {
	publicPem, err := keypair.PemMarshalPublic(a.Public)
	if err != nil {
		return nil, nil, err
	}
	privatePem, err := keypair.PemMarshalPrivate(a.Private)
	if err != nil {
		return nil, nil, err
	}
	return publicPem, privatePem, nil
}

func (a KeyPair) MarshalSSH() (public []byte, private []byte, err error) {
	publicSshKey, err := keypair.SshMarshalPublic(a.Public)
	if err != nil {
		return nil, nil, err
	}
	privatePem, err := keypair.SshMarshalPrivate(a.Private, "")
	if err != nil {
		return nil, nil, err
	}
	return publicSshKey, privatePem, nil
}
