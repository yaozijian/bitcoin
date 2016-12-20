package core

import (
	"encoding/hex"
	"fmt"

	"meeting"
)

func (u *user) onDecide(decide *meeting.Decide) {

	err := u.verifySignature(&decide.DecideCore, &decide.Signature, false)

	if err == nil {
		switch decide.DecideCore.Proposal.ProposalCore.Type {
		case ProposalType_Register, ProposalType_Pay, PrososalType_Bill:
			u.onAck(decide)
		}
	}
}

func (u *user) onAck(decide *meeting.Decide) {

	digest := hex.EncodeToString(decide.DecideCore.Proposal.Signature.Digest)
	reqctx := u.proposals_map[digest]

	if reqctx == nil {
		return
	}

	var err error
	var op func()

	if len(decide.DecideCore.Error) > 0 {
		if reqctx.oppose++; reqctx.oppose > u.onlines/2 {
			op = reqctx.onFail
			err = fmt.Errorf(decide.DecideCore.Error)
		} else {
			return
		}
	} else if reqctx.agree++; reqctx.agree > u.onlines/2 {
		op = reqctx.onAgree
	} else {
		return
	}

	if op != nil {
		op()
	}

	if reqctx.reply != nil {
		reqctx.reply <- err
	}

	delete(u.proposals_map, digest)
}

//----------------------------------------------------------------------------------------------

// 表决
func (u *user) sendDecide(proposal *meeting.Proposal, err error) {

	core := &meeting.DecideCore{
		Proposal: *proposal,
		Who:      u.name,
	}

	if err != nil {
		core.Error = err.Error()
	}

	sign, err := u.signature(core)

	if err == nil {
		decide := &meeting.Decide{
			DecideCore: *core,
			Signature:  *sign,
		}
		u.doBroadCast(decide, BroadCastType_Decide)
	}
}
