package keypair

import (
	"crypto"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"golang.org/x/crypto/ssh"
)

const (
	TypePemPublic  = "PUBLIC KEY"
	TypePemPrivate = "PRIVATE KEY"
)

func DecodePemBlock(content []byte, targetType string) (*pem.Block, error) {
	block, _ := pem.Decode(content)
	if block == nil {
		return nil, errors.New("not pem format")
	} else if block.Type != targetType {
		return nil, fmt.Errorf("pem type should be %s, got %s", targetType, block.Type)
	}
	return block, nil
}

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

func PemUnmarshalPublic[T crypto.PublicKey](content []byte) (key T, err error) {
	block, err := DecodePemBlock(content, TypePemPublic)
	if err != nil {
		return
	}

	publicKey, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		err = fmt.Errorf("parse public key failed: %v", err)
		return
	}
	key, ok := publicKey.(T)
	if !ok {
		err = errors.New("public key format error")
	}
	return
}
func PemUnmarshalPrivate[T crypto.PrivateKey](content []byte) (key T, err error) {
	block, err := DecodePemBlock(content, TypePemPrivate)
	if err != nil {
		return
	}

	privateKey, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		err = fmt.Errorf("parse private key failed: %v", err)
		return
	}
	key, ok := privateKey.(T)
	if !ok {
		err = errors.New("private key format error")
	}
	return
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
