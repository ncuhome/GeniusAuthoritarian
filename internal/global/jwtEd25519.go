package global

import (
	"crypto/ed25519"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	log "github.com/sirupsen/logrus"
	"os"
	"path"
)

func initJwtEd25519() {
	// load from pem format files
	const (
		PublicKeyName  = "public.pem"
		PrivateKeyName = "private.pem"
	)

	var err error
	var loader = jwtKeyLoader{
		Dir: path.Join(ConfigDir, "jwt"),
	}
	JwtEd25519.PrivateKey, err = loader.LoadPrivateKey(PrivateKeyName)
	if err != nil {
		log.Fatalln("load jwt private key failed:", err)
	}

	JwtEd25519.PublicKey, err = loader.LoadPublicKey(PublicKeyName)
	if err != nil {
		JwtEd25519.PublicKey = JwtEd25519.PrivateKey.Public().(ed25519.PublicKey)
		if os.IsNotExist(err) {
			log.Debugln("jwt public key not exist, generated from private key")
		} else {
			log.Warnln("jwt public key load failed, generated from private key:", err)
		}
	}
}

type _JwtEd25519 struct {
	PublicKey  ed25519.PublicKey
	PrivateKey ed25519.PrivateKey
}

var JwtEd25519 _JwtEd25519

type jwtKeyLoader struct {
	Dir string
}

func (loader jwtKeyLoader) LoadPemBlock(filename, targetType string) (*pem.Block, error) {
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

func (loader jwtKeyLoader) LoadPrivateKey(filename string) (ed25519.PrivateKey, error) {
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

func (loader jwtKeyLoader) LoadPublicKey(filename string) (ed25519.PublicKey, error) {
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
