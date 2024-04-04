package ca

import (
	"crypto/rand"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"github.com/ncuhome/GeniusAuthoritarian/pkg/ed25519"
	"github.com/ncuhome/GeniusAuthoritarian/pkg/keypair"
	"math/big"
	"time"
)

func NewRoot(valid time.Duration) (public []byte, private []byte, err error) {
	ed25519Keypair, err := ed25519.Generate()
	if err != nil {
		return nil, nil, err
	}

	template := x509.Certificate{
		SerialNumber: big.NewInt(time.Now().UnixNano()),
		Subject: pkix.Name{
			CommonName: "GeniusAuthoritarian Root",
		},
		NotBefore:             time.Now(),
		NotAfter:              time.Now().Add(valid),
		KeyUsage:              x509.KeyUsageCertSign | x509.KeyUsageCRLSign,
		IsCA:                  true,
		BasicConstraintsValid: true,
	}
	certBytes, err := x509.CreateCertificate(rand.Reader, &template, &template, ed25519Keypair.Public, ed25519Keypair.Private)
	if err != nil {
		return nil, nil, err
	}
	certPem := pem.EncodeToMemory(&pem.Block{
		Type:  "CERTIFICATE",
		Bytes: certBytes,
	})
	privateKeyPem, err := keypair.PemMarshalPrivate(ed25519Keypair.Private)
	if err != nil {
		return nil, nil, err
	}
	return certPem, privateKeyPem, nil
}
