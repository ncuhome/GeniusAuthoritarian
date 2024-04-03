package ed25519

import (
	"crypto/ed25519"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"os"
	"path"
)

type Loader struct {
	Dir string
}

func (loader Loader) LoadPemBlock(filename, targetType string) (*pem.Block, error) {
	keyPath := path.Join(loader.Dir, filename)
	keyBytes, err := os.ReadFile(keyPath)
	if err != nil {
		return nil, fmt.Errorf("read %s failed: %v", keyPath, err)
	}

	block, _ := pem.Decode(keyBytes)
	if block == nil {
		return nil, errors.New("not pem format")
	} else if block.Type != targetType {
		return nil, fmt.Errorf("pem type should be %s, got %s", targetType, block.Type)
	}
	return block, nil
}

func (loader Loader) LoadPrivateKey(filename string) (ed25519.PrivateKey, error) {
	block, err := loader.LoadPemBlock(filename, "PRIVATE KEY")
	if err != nil {
		return nil, err
	}

	privateKey, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("parse private key failed: %v", err)
	}

	ed25519PrivateKey, ok := privateKey.(ed25519.PrivateKey)
	if !ok {
		return nil, errors.New("not an ed25519 private key")
	}
	return ed25519PrivateKey, nil
}

func (loader Loader) LoadPublicKey(filename string) (ed25519.PublicKey, error) {
	block, err := loader.LoadPemBlock(filename, "PUBLIC KEY")
	if err != nil {
		return nil, err
	}

	publicKey, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("parse public key failed: %v", err)
	}

	ed25519PublicKey, ok := publicKey.(ed25519.PublicKey)
	if !ok {
		return nil, errors.New("not an ed25519 public key")
	}
	return ed25519PublicKey, nil
}
