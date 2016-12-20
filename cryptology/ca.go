package cryptology

import (
	"crypto/x509"
	"fmt"
)

type (
	ca struct {
		*publicKey
		pool *x509.CertPool
	}
)

func LoadCAFromFile(certfile string) (CA, error) {
	p, e := LoadPublicKeyFromFile(certfile)
	if e != nil {
		return nil, e
	} else {
		return newCAFromCertsList(p)
	}
}

func LoadCAFromBytes(certbytes []byte) (CA, error) {
	p, e := LoadPublicKeyFromBytes(certbytes)
	if e != nil {
		return nil, e
	} else {
		return newCAFromCertsList(p)
	}
}

func newCAFromCertsList(p PublicKey) (CA, error) {

	count := 0
	oldlist := p.certList()
	newlist := make([]*x509.Certificate, len(oldlist))
	pool := x509.NewCertPool()

	for _, cert := range oldlist {
		if cert.IsCA {
			pool.AddCert(cert)
			newlist[count] = cert
			count++
		}
	}

	if count > 0 {
		newlist = newlist[:count]
		c := &ca{pool: pool, publicKey: p.(*publicKey)}
		c.certlist = newlist
		return c, nil
	} else {
		return nil, fmt.Errorf("没有CA证书")
	}
}

func (c *ca) VerifyPublicKey(p PublicKey) error {

	opts := x509.VerifyOptions{
		Roots: c.pool,
	}

	for _, cert := range p.certList() {
		if _, err := cert.Verify(opts); err != nil {
			return err
		}
	}

	return nil
}
