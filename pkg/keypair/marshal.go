package keypair

import (
	"crypto/x509"
	"encoding/pem"
	"golang.org/x/crypto/ssh"
)

const (
	TypePemPublic  = "PUBLIC KEY"
	TypePemPrivate = "PRIVATE KEY"
)

func PemEncodePublic(content []byte) []byte {
	return pem.EncodeToMemory(&pem.Block{
		Type:  TypePemPublic,
		Bytes: content,
	})
}
func PemEncodePrivate(content []byte) []byte {
	return pem.EncodeToMemory(&pem.Block{
		Type:  TypePemPrivate,
		Bytes: content,
	})
}

func PemMarshalPublic(key any) ([]byte, error) {
	publicPem, err := x509.MarshalPKIXPublicKey(key)
	if err != nil {
		return nil, err
	}
	return PemEncodePublic(publicPem), nil
}
func PemMarshalPrivate(key any) ([]byte, error) {
	privatePem, err := x509.MarshalPKCS8PrivateKey(key)
	if err != nil {
		return nil, err
	}
	return PemEncodePrivate(privatePem), nil
}

func SshMarshalPublic(key any) ([]byte, error) {
	publicSshKey, err := ssh.NewPublicKey(key)
	if err != nil {
		return nil, err
	}
	return ssh.MarshalAuthorizedKey(publicSshKey), nil
}
func SshMarshalPrivate(key any, comment string) ([]byte, error) {
	privatePem, err := ssh.MarshalPrivateKey(key, comment)
	if err != nil {
		return nil, err
	}
	return pem.EncodeToMemory(privatePem), nil
}
