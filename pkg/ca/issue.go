package ca

import (
	"crypto/ed25519"
	"crypto/rand"
	"crypto/x509"
	"crypto/x509/pkix"
	"errors"
	ed25519Pkg "github.com/ncuhome/GeniusAuthoritarian/pkg/ed25519"
	"github.com/ncuhome/GeniusAuthoritarian/pkg/keypair"
	"math/big"
	"time"
)

func NewIssuer(certPem, keyPem []byte) (*Issuer, error) {
	caCert, err := keypair.PemUnmarshalCertificate(certPem)
	if err != nil {
		return nil, err
	}
	if caCert.NotAfter.Before(time.Now()) {
		return nil, errors.New("cert has expired")
	}

	caPrivateKey, err := keypair.PemUnmarshalPKCS8Private[ed25519.PrivateKey](keypair.FormatECDSA, keyPem)
	if err != nil {
		return nil, err
	}

	return &Issuer{
		CaCertPem: certPem,
		CaCert:    caCert,
		caKey:     caPrivateKey,
	}, nil
}

type Issuer struct {
	CaCertPem []byte
	CaCert    *x509.Certificate

	caKey ed25519.PrivateKey
}

func (i Issuer) issue(certificate *x509.Certificate) (fullChain, private []byte, err error) {
	ed25519Keypair, err := ed25519Pkg.Generate()
	if err != nil {
		return nil, nil, err
	}

	clientCert, err := x509.CreateCertificate(rand.Reader, certificate, i.CaCert, ed25519Keypair.Public, i.caKey)
	if err != nil {
		return nil, nil, err
	}

	certPem := keypair.PemEncodeCertificate(clientCert)
	fullChain = make([]byte, 0, len(certPem)+len(i.CaCertPem))
	fullChain = append(fullChain, certPem...)
	fullChain = append(fullChain, i.CaCertPem...)

	privatePem := keypair.PemEncodePrivate(keypair.FormatECDSA, ed25519Keypair.Private)

	return fullChain, privatePem, nil
}

func (i Issuer) IssueServer(serverName string, notAfter time.Time) (fullChain, private []byte, err error) {
	return i.issue(&x509.Certificate{
		SerialNumber: big.NewInt(time.Now().UnixNano()),
		Subject: pkix.Name{
			CommonName: serverName,
		},
		NotBefore:   time.Now().Add(-time.Minute),
		NotAfter:    notAfter,
		KeyUsage:    x509.KeyUsageDigitalSignature | x509.KeyUsageKeyEncipherment,
		ExtKeyUsage: []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
	})
}

func (i Issuer) IssueClient(serverName string, notAfter time.Time) (fullChain, private []byte, err error) {
	return i.issue(&x509.Certificate{
		SerialNumber: big.NewInt(time.Now().UnixNano()),
		Subject: pkix.Name{
			CommonName: serverName,
		},
		NotBefore:   time.Now().Add(-time.Minute),
		NotAfter:    notAfter,
		KeyUsage:    x509.KeyUsageDigitalSignature | x509.KeyUsageKeyEncipherment,
		ExtKeyUsage: []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth},
	})
}
