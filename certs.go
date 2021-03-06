package main

import (
	"fmt"

	"github.com/square/certstrap/pkix"
)

const KeyBits = 2048
const CAExpiryYears = 10
const HostCertExpiryYears = 2

type CA struct {
	VarName_CA string
	CommonName string

	key  *pkix.Key
	cert *pkix.Certificate
}

type CertKeyPair struct {
	VarName_Cert string
	VarName_Key  string
	CommonName   string
	Domains      []string

	key  *pkix.Key
	cert *pkix.Certificate
}

type PlainKeyPair struct {
	VarName_PublicKey  string
	VarName_PrivateKey string
}

type exportable interface {
	Export() ([]byte, error)
}

func asString(e exportable) (string, error) {
	pemBytes, err := e.Export()
	if err != nil {
		return "", fmt.Errorf("export pem: %s", err)
	}

	return string(pemBytes), nil
}

func (ca *CA) Init() error {
	var err error
	ca.key, err = pkix.CreateRSAKey(KeyBits)
	if err != nil {
		return fmt.Errorf("create ca key: %s", err)
	}

	ca.cert, err = pkix.CreateCertificateAuthority(ca.key, "", CAExpiryYears,
		"", "", "", "", ca.CommonName)
	if err != nil {
		return fmt.Errorf("create ca cert: %s", err)
	}

	return nil
}

func (ca *CA) CACertAsString() (string, error) {
	return asString(ca.cert)
}

func (ca *CA) InitCertKeyPair(certKeyPair *CertKeyPair) error {
	var err error
	certKeyPair.key, err = pkix.CreateRSAKey(KeyBits)
	if err != nil {
		return fmt.Errorf("create host key: %s", err)
	}
	csr, err := pkix.CreateCertificateSigningRequest(certKeyPair.key, "", nil,
		certKeyPair.Domains, "", "", "", "", certKeyPair.CommonName)
	if err != nil {
		return fmt.Errorf("create host csr: %s", err)
	}

	certKeyPair.cert, err = pkix.CreateCertificateHost(ca.cert, ca.key, csr, HostCertExpiryYears)
	if err != nil {
		return fmt.Errorf("sign host csr: %s", err)
	}

	return nil
}

func (kp *CertKeyPair) PrivateKeyAsString() (string, error) {
	pemBytes, err := kp.key.ExportPrivate()
	if err != nil {
		return "", fmt.Errorf("export private key: %s", err)
	}

	return string(pemBytes), nil
}

func (kp *CertKeyPair) CertAsString() (string, error) {
	return asString(kp.cert)
}
