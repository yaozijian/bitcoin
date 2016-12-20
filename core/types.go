package core

import (
	"bytes"
	"time"

	"github.com/yaozijian/bitcoin/meeting"
)

type (
	RegisterUser struct {
		When time.Time
		UserRecord
	}

	Pay struct {
		When    time.Time
		PayFrom string // 付款方
		PayTo   string // 收款方
		Amount  int64  // 现实世界的金额
		Coins   int64  // 虚拟货币额
	}

	Bill struct {
		When        time.Time
		PayDigest   []byte // 支付提案摘要
		PrevDigest  []byte // 上一账单摘要
		Collier     string // 矿工名称
		MagicNumber uint64 // 幸运号
	}

	BillRecord struct {
		Bill
		meeting.Signature
	}

	UserRecord struct {
		Name   string
		Amount int64
		Coins  int64
	}
)

const (
	ProposalType_Error    = iota
	ProposalType_Register // 账户注册
	ProposalType_Pay      // 支付
	PrososalType_Bill     // 账单
)

const (
	BroadCastType_Proposal = iota
	BroadCastType_Decide
)

var (
	validBillDigestPrefix []byte
	firstBillDigest       []byte
	rewardEveryBill       int64
)

func init() {
	firstBillDigest = bytes.Repeat([]byte{0}, 64)
	validBillDigestPrefix = bytes.Repeat([]byte{0xA}, 1)
	rewardEveryBill = 50
}

func FirstBillDigest() []byte {
	t := make([]byte, len(firstBillDigest))
	copy(t, firstBillDigest)
	return t
}

func ValidBillDigestPrefix() []byte {
	t := make([]byte, len(validBillDigestPrefix))
	copy(t, validBillDigestPrefix)
	return t
}

func RewardEveryBill() int64 {
	return rewardEveryBill
}
