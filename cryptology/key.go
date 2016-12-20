package cryptology

import (
	"crypto"
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"io/ioutil"
)

type (
	privateKey struct {
		key interface{}
	}
)

func LoadPrivateKey(keyfile string) (private PrivateKey, err error) {

	cont, x := ioutil.ReadFile(keyfile)
	if x != nil {
		return nil, x
	}

	var key interface{}

	if key, err = parsePrivateKey(cont); err == nil && key != nil {
		private = &privateKey{key: key}
	} else if key = loadPEMKey(cont); key != nil {
		err = nil
		private = &privateKey{key: key}
	} else {
		err = fmt.Errorf("解析私钥文件失败")
	}

	return
}

func (key *privateKey) Sign(hasnType crypto.Hash, hashVal []byte) (signature []byte, err error) {

	switch x := key.key.(type) {
	case *ecdsa.PrivateKey:
		opts := &rsa.PSSOptions{SaltLength: 10, Hash: hasnType}
		signature, err = x.Sign(rand.Reader, hashVal, opts)
	case *rsa.PrivateKey:
		signature, err = rsa.SignPKCS1v15(rand.Reader, x, hasnType, hashVal)
	default:
		err = fmt.Errorf("内部错误,无法执行签名操作")
	}
	return
}

func loadPEMKey(pemConts []byte) interface{} {

	var block *pem.Block

	for len(pemConts) > 0 {

		if block, pemConts = pem.Decode(pemConts); block == nil {
			break
		}

		if block.Type == "PRIVATE KEY" && len(block.Bytes) > 0 {
			if key, err := parsePrivateKey(block.Bytes); err == nil && key != nil {
				return key
			}
		}
	}

	return nil
}

// 从标准包tls中tls.go文件末尾抄过来的
func parsePrivateKey(der []byte) (crypto.PrivateKey, error) {
	if key, err := x509.ParsePKCS1PrivateKey(der); err == nil {
		return key, nil
	}
	if key, err := x509.ParsePKCS8PrivateKey(der); err == nil {
		switch key := key.(type) {
		case *rsa.PrivateKey, *ecdsa.PrivateKey:
			return key, nil
		default:
			return nil, errors.New("crypto/tls: found unknown private key type in PKCS#8 wrapping")
		}
	}
	if key, err := x509.ParseECPrivateKey(der); err == nil {
		return key, nil
	}

	return nil, errors.New("crypto/tls: failed to parse private key")
}
