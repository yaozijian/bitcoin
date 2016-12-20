package core

import (
	"bytes"
	"time"

	"golang.org/x/net/context"

	"meeting"

	log "github.com/cihub/seelog"
)

func (u *user) dig(proposal *meeting.Proposal, ctx context.Context) {

	bill := &Bill{
		When:       time.Now().UTC(),
		PayDigest:  proposal.Signature.Digest,
		PrevDigest: u.last_bill_digest,
		Collier:    u.name,
	}

	core := &meeting.ProposalCore{
		Who:  u.name,
		Type: PrososalType_Bill,
	}

	var data []byte
	var sign *meeting.Signature
	var err error

	for err == nil {

		if data, err = u.encode(bill); err != nil {
			break
		}

		core.Data = data

		if sign, err = u.signature(core); err != nil {
			break
		}

		// 摘要前缀满足要求则停止
		if bytes.HasPrefix(sign.Digest, validBillDigestPrefix) {
			break
		} else {
			select {
			// 别的矿工已经从这个账单挖到矿了,停止继续尝试
			case <-ctx.Done():
				log.Infof("用户%v: 别人已经挖到矿石了,停止继续尝试", u.name)
				return
			default:
				// 否则增加MagicNumber,继续尝试
				bill.MagicNumber++
			}
		}
	}

	if err != nil {
		// 内部错误,无法完成工作
		log.Infof("xxx 用户%v: 内部错误: %v", u.name, err)
		u.sendDecide(proposal, err)
	} else {
		log.Infof("用户%v: 挖到矿石了 Magic=%v", u.name, bill.MagicNumber)
		// 发送账单
		billctx := &proposalContext{
			proposal: meeting.Proposal{
				ProposalCore: *core,
				Signature:    *sign,
			},
		}
		u.proposal_chnl <- billctx
	}
}
