package ca

import (
	"crypto/ed25519"
	"crypto/rand"
	"crypto/x509"
	"crypto/x509/pkix"
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
	caPrivateKey, err := keypair.PemUnmarshalPrivate[ed25519.PrivateKey](keypair.FormatECDSA, keyPem)
	if err != nil {
		return nil, err
	}

	return &Issuer{
		caCertPem: certPem,
		caCert:    caCert,
		caKey:     caPrivateKey,
	}, nil
}

type Issuer struct {
	caCertPem []byte
	caCert    *x509.Certificate

	caKey ed25519.PrivateKey
}

func (i Issuer) Issue(dnsNames []string, notAfter time.Time) (fullChain, private []byte, err error) {
	ed25519Keypair, err := ed25519Pkg.Generate()
	if err != nil {
		return nil, nil, err
	}

	clientCert, err := x509.CreateCertificate(rand.Reader, &x509.Certificate{
		SerialNumber: big.NewInt(time.Now().UnixNano()),
		Subject: pkix.Name{
			CommonName: "GeniusAuthoritarian Client Cert",
		},
		NotBefore:   time.Now().Add(-time.Minute),
		NotAfter:    notAfter,
		KeyUsage:    x509.KeyUsageDigitalSignature | x509.KeyUsageKeyEncipherment,
		ExtKeyUsage: []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth},
		DNSNames:    dnsNames,
	}, i.caCert, ed25519Keypair.Public, i.caKey)
	if err != nil {
		return nil, nil, err
	}

	certPem := keypair.PemEncodeCertificate(clientCert)
	fullChain = make([]byte, 0, len(certPem)+len(i.caCertPem))
	fullChain = append(fullChain, certPem...)
	fullChain = append(fullChain, i.caCertPem...)

	privatePem := keypair.PemEncodePrivate(keypair.FormatECDSA, ed25519Keypair.Private)

	return fullChain, privatePem, nil
}
