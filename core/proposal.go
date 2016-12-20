package core

import (
	"bytes"
	"encoding/hex"
	"fmt"

	"golang.org/x/net/context"

	"meeting"
)

// 处理议案
func (u *user) onProposal(proposal *meeting.Proposal) {

	var err error

	// 签名验证
	switch proposal.ProposalCore.Type {
	case ProposalType_Register:
		err = u.verifySignature(&proposal.ProposalCore, &proposal.Signature, true)
	default:
		err = u.verifySignature(&proposal.ProposalCore, &proposal.Signature, false)
	}

	if err != nil {
		u.sendDecide(proposal, err)
		return
	}

	// 分别处理
	switch proposal.ProposalCore.Type {
	case ProposalType_Register:
		reginfo := new(RegisterUser)
		if err = u.decode(proposal.ProposalCore.Data, reginfo); err == nil {
			u.onRegisterUser(proposal, reginfo)
		}
	case ProposalType_Pay:
		payinfo := new(Pay)
		if err = u.decode(proposal.ProposalCore.Data, payinfo); err == nil {
			u.onPayInfo(proposal, payinfo)
		}
	case PrososalType_Bill:
		bill := new(Bill)
		if err = u.decode(proposal.ProposalCore.Data, bill); err == nil {
			u.onBill(proposal, bill)
		}
	}

	if err != nil {
		u.sendDecide(proposal, err)
	}
}

// 账号注册
func (u *user) onRegisterUser(proposal *meeting.Proposal, reginfo *RegisterUser) {

	if user := u.accounts_map[reginfo.Name]; user != nil {
		u.sendDecide(proposal, fmt.Errorf("用户已经存在"))
	} else {
		user = &UserRecord{
			Name:   reginfo.Name,
			Amount: reginfo.Amount,
			Coins:  reginfo.Coins,
		}

		reqctx := u.add2ProposalMap(proposal, false)

		reqctx.onAgree = func() { u.accounts_map[user.Name] = user }

		u.sendDecide(proposal, nil)
	}
}

// 支付
func (u *user) onPayInfo(proposal *meeting.Proposal, payinfo *Pay) {

	from := u.accounts_map[payinfo.PayFrom]
	to := u.accounts_map[payinfo.PayTo]

	if to == nil {
		u.sendDecide(proposal, fmt.Errorf("找不到收款方"))
	} else if from.Amount < payinfo.Amount || from.Coins < payinfo.Coins {
		u.sendDecide(proposal, fmt.Errorf("付款方金额不足"))
	} else {

		reqctx := u.add2ProposalMap(proposal, true)

		reqctx.onAgree = func() {
			from.Amount -= payinfo.Amount
			from.Coins -= payinfo.Coins
			to.Amount += payinfo.Amount
			to.Coins += payinfo.Coins
		}

		go u.dig(proposal, reqctx.ctx)
	}
}

// 账单
func (u *user) onBill(proposal *meeting.Proposal, bill *Bill) {

	var err error

	payctx := u.proposals_map[hex.EncodeToString(bill.PayDigest)]
	who := u.accounts_map[bill.Collier]

	if payctx == nil {
		err = fmt.Errorf("账单错误: 找不到对应的支付议案")
	} else if who == nil {
		err = fmt.Errorf("找不到矿工记录: %v", bill.Collier)
	} else if bytes.Equal(bill.PrevDigest, FirstBillDigest()) {
		err = nil
	} else if prev := u.bills_map[hex.EncodeToString(bill.PrevDigest)]; prev != nil {
		err = fmt.Errorf("前一账单摘要错误: 找不到前一账单")
	}

	if err == nil {

		reqctx := u.add2ProposalMap(proposal, false)

		record := &BillRecord{
			Bill:      *bill,
			Signature: proposal.Signature,
		}

		reqctx.onAgree = func() {

			// 完成支付
			payctx.onAgree()

			delete(u.proposals_map, hex.EncodeToString(payctx.proposal.Signature.Digest))

			if payctx.reply != nil {
				payctx.reply <- nil
			}

			// 停止挖矿
			if payctx.cancel != nil {
				payctx.cancel()
			}

			// 账单奖赏
			who.Coins += RewardEveryBill()

			// 记录账单
			u.bills_map[hex.EncodeToString(reqctx.proposal.Signature.Digest)] = record
		}

		u.sendDecide(proposal, nil)
	} else {
		u.sendDecide(proposal, err)
	}
}

//----------------------------------------------------------------------------------------------

func (u *user) add2ProposalMap(proposal *meeting.Proposal, needCancel bool) (reqctx *proposalContext) {

	digest := hex.EncodeToString(proposal.Signature.Digest)
	reqctx = u.proposals_map[digest]

	if reqctx == nil {
		reqctx = &proposalContext{proposal: *proposal}
		if needCancel {
			ctx, cancel := context.WithCancel(context.Background())
			reqctx.ctx = ctx
			reqctx.cancel = cancel
		}
		u.proposals_map[digest] = reqctx
	}

	return
}

func (u *user) createProposal(data interface{}, typ uint64) (ctx *proposalContext, err error) {

	enc, e := u.encode(data)

	if e != nil {
		err = e
		return
	}

	core := &meeting.ProposalCore{
		Who:  u.name,
		Type: typ,
		Data: enc,
	}

	// 生成签名
	sign, e := u.signature(core)

	if e != nil {
		err = e
	} else {
		ctx = &proposalContext{
			proposal: meeting.Proposal{
				ProposalCore: *core,
				Signature:    *sign,
			},
			reply: make(chan error, 1),
		}
	}

	return
}
