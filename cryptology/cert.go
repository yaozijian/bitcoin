package cryptology

import (
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io/ioutil"
)

type (
	publicKey struct {
		certlist  []*x509.Certificate
		certbytes []byte
	}
)

func LoadPublicKeyFromFile(certfile string) (PublicKey, error) {
	if certbytes, list, err := LoadCertsFromFile(certfile); err == nil {
		return &publicKey{certlist: list, certbytes: certbytes}, nil
	} else {
		return nil, err
	}
}

func LoadPublicKeyFromBytes(certbytes []byte) (PublicKey, error) {
	if list, err := LoadCertsFromBytes(certbytes); err == nil {
		return &publicKey{certlist: list, certbytes: certbytes}, nil
	} else {
		return nil, err
	}
}

func (p *publicKey) VerifySignature(algo x509.SignatureAlgorithm, data, signature []byte) error {
	for _, cert := range p.certlist {
		if err := cert.CheckSignature(algo, data, signature); err == nil {
			return nil
		}
	}
	return fmt.Errorf("签名错误")
}

func (p *publicKey) Certificate() []byte {
	t := make([]byte, len(p.certbytes))
	copy(t, p.certbytes)
	return t
}

func (p *publicKey) certList() []*x509.Certificate {
	return p.certlist
}

//--------------------------------------------------------------------------

func LoadCertsFromFile(certfile string) ([]byte, []*x509.Certificate, error) {

	certbytes, err := ioutil.ReadFile(certfile)
	if err != nil {
		return nil, nil, err
	}

	a, b := LoadCertsFromBytes(certbytes)
	return certbytes, a, b
}

func LoadCertsFromBytes(certbytes []byte) ([]*x509.Certificate, error) {

	list := loadPEMCerts(certbytes)
	list = append(list, loadDERCerts(certbytes)...)

	if len(list) == 0 {
		return nil, fmt.Errorf("空证书")
	} else {
		return list, nil
	}
}

//--------------------------------------------------------------------------

func loadPEMCerts(pemCerts []byte) (list []*x509.Certificate) {

	var block *pem.Block

	for len(pemCerts) > 0 {

		if block, pemCerts = pem.Decode(pemCerts); block == nil {
			break
		}

		if block.Type == "CERTIFICATE" && len(block.Bytes) > 0 {
			cert, err := x509.ParseCertificate(block.Bytes)
			if err == nil {
				list = append(list, cert)
			}
		}
	}

	return
}

func loadDERCerts(derCerts []byte) (list []*x509.Certificate) {
	if certlist, err := x509.ParseCertificates(derCerts); err == nil {
		for _, cert := range certlist {
			list = append(list, cert)
		}
	}
	return
}
