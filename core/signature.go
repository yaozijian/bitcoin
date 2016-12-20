package core

import (
	"bytes"
	"crypto"
	"crypto/sha512"
	"crypto/x509"
	"encoding/gob"
	"fmt"

	"cryptology"
	"meeting"
)

// 为数据生成签名
func (u *user) signature(data interface{}) (*meeting.Signature, error) {

	cont, err := u.encode(data)

	if err != nil {
		return nil, err
	}

	// 注意: 两次执行摘要操作,使用私钥对摘要的摘要进行签名
	hash := sha512.Sum512(cont)
	digest := sha512.Sum512(hash[:])

	sign, err := u.private.Sign(crypto.SHA512, digest[:])

	if err != nil {
		return nil, fmt.Errorf("私钥签名失败: %v", err)
	} else {
		retval := &meeting.Signature{
			Digest:    hash[:],
			Signature: sign,
			PublicKey: u.public.Certificate(),
		}
		return retval, nil
	}
}

func (u *user) verifySignature(data interface{}, sign *meeting.Signature, useca bool) error {
	if p, e := cryptology.LoadPublicKeyFromBytes(sign.PublicKey); e != nil {
		return fmt.Errorf("载入签名携带的公钥失败: %v", e)
	} else if e = u.ca.VerifyPublicKey(p); e != nil {
		return fmt.Errorf("签名携带的公钥未通过CA验证: %v", e)
	} else if useca {
		if e = u.ca.VerifySignature(x509.SHA512WithRSA, sign.Digest, sign.Signature); e != nil {
			return fmt.Errorf("使用CA证书验证摘要和签名失败: %v", e)
		} else {
			return nil
		}
	} else if e = p.VerifySignature(x509.SHA512WithRSA, sign.Digest, sign.Signature); e != nil {
		return fmt.Errorf("使用签名携带的公钥验证摘要和签名失败: %v", e)
	} else {
		return nil
	} /* else if mysign, e := u.signature(data); e == nil {
		if bytes.Equal(mysign.Digest, sign.Digest) {
			return nil
		} else {
			return fmt.Errorf("签名携带的摘要与数据摘要不匹配")
		}
	} else {
		return fmt.Errorf("内部错误: %v", e)
	}*/
}

// 数据编码
func (u *user) encode(data interface{}) ([]byte, error) {
	buf := bytes.NewBuffer(nil)
	enc := gob.NewEncoder(buf)
	err := enc.Encode(data)
	if err == nil {
		return buf.Bytes(), nil
	} else {
		return nil, fmt.Errorf("Gob编码失败: %v", err)
	}
}

// 数据解码
func (u *user) decode(data []byte, ptr interface{}) error {
	buf := bytes.NewBuffer(data)
	dec := gob.NewDecoder(buf)
	err := dec.Decode(ptr)
	if err != nil {
		err = fmt.Errorf("Gob解码失败: %v", err)
	}
	return err
}
