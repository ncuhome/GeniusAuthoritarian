package ca

import (
	"crypto/rand"
	"crypto/x509"
	"crypto/x509/pkix"
	"github.com/ncuhome/GeniusAuthoritarian/pkg/ed25519"
	"github.com/ncuhome/GeniusAuthoritarian/pkg/keypair"
	"math/big"
	"time"
)

func NewRoot(notAfter time.Time) (public []byte, private []byte, err error) {
	ed25519Keypair, err := ed25519.Generate()
	if err != nil {
		return nil, nil, err
	}

	template := x509.Certificate{
		SerialNumber: big.NewInt(time.Now().UnixNano()),
		Subject: pkix.Name{
			CommonName: "GeniusAuthoritarian Root",
		},
		NotBefore:             time.Now().Add(-time.Minute),
		NotAfter:              notAfter,
		KeyUsage:              x509.KeyUsageCertSign | x509.KeyUsageCRLSign,
		IsCA:                  true,
		BasicConstraintsValid: true,
	}
	certBytes, err := x509.CreateCertificate(rand.Reader, &template, &template, ed25519Keypair.Public, ed25519Keypair.Private)
	if err != nil {
		return nil, nil, err
	}
	certPem := keypair.PemEncodeCertificate(certBytes)
	privateKeyPem, err := keypair.PemMarshalPrivate(keypair.FormatECDSA, ed25519Keypair.Private)
	if err != nil {
		return nil, nil, err
	}
	return certPem, privateKeyPem, nil
}
