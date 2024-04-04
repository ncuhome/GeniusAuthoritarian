package keypair

import (
	"crypto"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"golang.org/x/crypto/ssh"
)

func DecodePemBlock(header PemType, content []byte) (*pem.Block, error) {
	block, _ := pem.Decode(content)
	if block == nil {
		return nil, errors.New("not pem format")
	} else if headerStr := header.String(); block.Type != headerStr {
		return nil, fmt.Errorf("pem type should be %s, got %s", headerStr, block.Type)
	}
	return block, nil
}

func PemEncodeCertificate(content []byte) []byte {
	return pem.EncodeToMemory(&pem.Block{
		Type:  TypeCertificate.String(),
		Bytes: content,
	})
}
func PemEncodePublic(format Format, content []byte) []byte {
	return pem.EncodeToMemory(&pem.Block{
		Type:  TypePemPublic.Format(format).String(),
		Bytes: content,
	})
}
func PemEncodePrivate(format Format, content []byte) []byte {
	return pem.EncodeToMemory(&pem.Block{
		Type:  TypePemPrivate.Format(format).String(),
		Bytes: content,
	})
}

func PemMarshalPublic(format Format, key crypto.PublicKey) ([]byte, error) {
	publicPem, err := x509.MarshalPKIXPublicKey(key)
	if err != nil {
		return nil, err
	}
	return PemEncodePublic(format, publicPem), nil
}
func PemMarshalPrivate(format Format, key crypto.PrivateKey) ([]byte, error) {
	privatePem, err := x509.MarshalPKCS8PrivateKey(key)
	if err != nil {
		return nil, err
	}
	return PemEncodePrivate(format, privatePem), nil
}

func PemUnmarshalCertificate(content []byte) (*x509.Certificate, error) {
	block, err := DecodePemBlock(TypeCertificate, content)
	if err != nil {
		return nil, err
	}
	return x509.ParseCertificate(block.Bytes)
}
func PemUnmarshalPublic[T crypto.PublicKey](format Format, content []byte) (key T, err error) {
	block, err := DecodePemBlock(TypePemPublic.Format(format), content)
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
func PemUnmarshalPrivate[T crypto.PrivateKey](format Format, content []byte) (key T, err error) {
	block, err := DecodePemBlock(TypePemPrivate.Format(format), content)
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

func SshMarshalPublic(key crypto.PublicKey) ([]byte, error) {
	publicSshKey, err := ssh.NewPublicKey(key)
	if err != nil {
		return nil, err
	}
	return ssh.MarshalAuthorizedKey(publicSshKey), nil
}
func SshMarshalPrivate(key crypto.PrivateKey, comment string) ([]byte, error) {
	privatePem, err := ssh.MarshalPrivateKey(key, comment)
	if err != nil {
		return nil, err
	}
	return pem.EncodeToMemory(privatePem), nil
}
