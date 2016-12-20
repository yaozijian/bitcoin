package core

import (
	"meeting"
	"router"
)

func (u *user) doBroadCast(data meeting.Signatureable, typ uint64) error {

	sign, err := u.signature(data)

	if err != nil {
		return err
	}

	data.SetSignature(sign)

	content, err := u.encode(data)

	if err == nil {
		b := router.BroadCastData{
			Type: typ,
			Data: content,
		}
		u.router.BroadCast(b)
	}

	return err
}

func (u *user) OnBroadCast(d router.BroadCastData) {

	t := make([]byte, len(d.Data))
	copy(t, d.Data)
	d.Data = t
	select {
	case u.broadcast_chnl <- d:
	default:
	}
}

func (u *user) onBroadCast(d router.BroadCastData) {

	u.onlines = d.Users

	switch d.Type {
	case BroadCastType_Proposal:
		proposal := new(meeting.Proposal)
		err := u.decode(d.Data, proposal)
		if err == nil {
			u.onProposal(proposal)
		}
	case BroadCastType_Decide:
		decide := new(meeting.Decide)
		err := u.decode(d.Data, decide)
		if err == nil {
			u.onDecide(decide)
		}
	}
}
