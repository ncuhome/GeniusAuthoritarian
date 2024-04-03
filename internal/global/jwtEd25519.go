package global

import (
	"crypto/ed25519"
	ed25519Pkg "github.com/ncuhome/GeniusAuthoritarian/pkg/ed25519"
	log "github.com/sirupsen/logrus"
	"path"
)

func initJwtEd25519() {
	// load from pem format files
	const (
		PublicKeyName  = "public.pem"
		PrivateKeyName = "private.pem"
	)

	var err error
	var loader = ed25519Pkg.Loader{
		Dir: path.Join(ConfigDir, "jwt"),
	}
	JwtEd25519.PrivateKey, err = loader.LoadPrivateKey(PrivateKeyName)
	if err != nil {
		log.Fatalln("load jwt private key failed:", err)
	}

	JwtEd25519.PublicKey, err = loader.LoadPublicKey(PublicKeyName)
	if err != nil {
		JwtEd25519.PublicKey = JwtEd25519.PrivateKey.Public().(ed25519.PublicKey)
		log.Warnln("load jwt public key failed, generated from private key:", err)
	}
}

type _JwtEd25519 struct {
	PublicKey  ed25519.PublicKey
	PrivateKey ed25519.PrivateKey
}

var JwtEd25519 _JwtEd25519
