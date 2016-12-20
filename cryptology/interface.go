package cryptology

import (
	"crypto"
	"crypto/x509"
)

type (
	PrivateKey interface {
		Sign(hasnType crypto.Hash, hashVal []byte) (signature []byte, err error)
	}

	PublicKey interface {
		VerifySignature(algo x509.SignatureAlgorithm, digest, signature []byte) error
		Certificate() []byte
		certList() []*x509.Certificate
	}

	CA interface {
		PublicKey
		VerifyPublicKey(PublicKey) error
	}
)
