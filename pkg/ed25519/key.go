package ed25519

import (
	"crypto/ed25519"
	"crypto/x509"
	"encoding/pem"
	"golang.org/x/crypto/ssh"
)

type Key struct {
	Public  ed25519.PublicKey
	Private ed25519.PrivateKey
}

func (a Key) MarshalPem() (public []byte, private []byte, err error) {
	publicPem, err := x509.MarshalPKIXPublicKey(a.Public)
	if err != nil {
		return nil, nil, err
	}

	privatePem, err := x509.MarshalPKCS8PrivateKey(a.Private)
	if err != nil {
		return nil, nil, err
	}

	return pem.EncodeToMemory(&pem.Block{
			Type:  "PUBLIC KEY",
			Bytes: publicPem,
		}), pem.EncodeToMemory(&pem.Block{
			Type:  "PRIVATE KEY",
			Bytes: privatePem,
		}), nil
}

func (a Key) MarshalSSH() (public []byte, private []byte, err error) {
	publicSshKey, err := ssh.NewPublicKey(a.Public)
	if err != nil {
		return nil, nil, err
	}

	return ssh.MarshalAuthorizedKey(publicSshKey), pem.EncodeToMemory(&pem.Block{
		Type:  "OPENSSH PRIVATE KEY",
		Bytes: MarshalOpenSshPrivateKey(a.Private),
	}), nil
}
