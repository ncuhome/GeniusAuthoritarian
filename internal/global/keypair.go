package global

import (
	"crypto/ed25519"
	"github.com/ncuhome/GeniusAuthoritarian/pkg/ca"
	ed25519Pkg "github.com/ncuhome/GeniusAuthoritarian/pkg/ed25519"
	log "github.com/sirupsen/logrus"
	"path"
)

const (
	_PublicKeyName  = "public.pem"
	_PrivateKeyName = "private.pem"

	_PublicTlsKeyName  = "tls.crt"
	_PrivateTlsKeyName = "tls.key"

	dirJwt = "jwt"
	dirCa  = "ca"
)

func initJwtEd25519() {
	var err error
	var loader = ed25519Pkg.Loader{
		Dir: path.Join(ConfigDir, dirJwt),
	}
	JwtEd25519.PrivateKey, err = loader.LoadPrivateKey(_PrivateKeyName)
	if err != nil {
		log.Fatalln("load jwt private key failed:", err)
	}

	JwtEd25519.PublicKey, err = loader.LoadPublicKey(_PublicKeyName)
	if err != nil {
		JwtEd25519.PublicKey = JwtEd25519.PrivateKey.Public().(ed25519.PublicKey)
		log.Warnln("load jwt public key failed, generated from private key:", err)
	}
}

func initCaIssuer() {
	var loader = ed25519Pkg.Loader{
		Dir: path.Join(ConfigDir, dirCa),
	}
	certBytes, err := loader.LoadPem(_PublicTlsKeyName)
	if err != nil {
		log.Fatalln("read ca cert failed:", err)
	}
	keyBytes, err := loader.LoadPem(_PrivateTlsKeyName)
	if err != nil {
		log.Fatalln("read ca key failed:", err)
	}
	CaIssuer, err = ca.NewIssuer(certBytes, keyBytes)
	if err != nil {
		log.Fatalln("init ca issuer failed:", err)
	}
}

type _Ed25519Keypair struct {
	PublicKey  ed25519.PublicKey
	PrivateKey ed25519.PrivateKey
}

var JwtEd25519 _Ed25519Keypair
var CaIssuer *ca.Issuer
