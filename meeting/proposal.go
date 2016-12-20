package meeting

type (
	Signature struct {
		Digest    []byte // 摘要
		Signature []byte // 签名
		PublicKey []byte // 公钥
	}

	ProposalCore struct {
		Who  string // 议员
		Type uint64 // 类型
		Data []byte // 数据
	}

	Proposal struct {
		ProposalCore
		Signature
	}

	Signatureable interface {
		GetSignature() *Signature
		SetSignature(*Signature)
	}
)

func (s *Signature) GetSignature() *Signature {
	return s
}

func (s *Signature) SetSignature(x *Signature) {
	*s = *x
}
