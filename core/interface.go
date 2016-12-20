package core

import (
	"fmt"
	"time"
)

func (u *user) Register(who string, amount, coins int64) error {

	var err error

	if len(who) == 0 {
		err = fmt.Errorf("用户名不能为空")
	} else if amount < 0 || coins < 0 {
		err = fmt.Errorf("金额不能小于零")
	}

	if err != nil {
		return err
	}

	//---------------------------------------

	reginfo := &RegisterUser{
		When: time.Now().UTC(),
		UserRecord: UserRecord{
			Name:   who,
			Amount: amount,
			Coins:  coins,
		},
	}

	ctx, err := u.createProposal(reginfo, ProposalType_Register)

	if err != nil {
		return err
	}

	//---------------------------------------

	u.proposal_chnl <- ctx

	ack := <-ctx.reply

	close(ctx.reply)

	return ack
}

func (u *user) PayTo(who string, amount, coins int64) error {

	var err error

	if len(who) == 0 {
		err = fmt.Errorf("用户名不能为空")
	} else if amount < 0 || coins < 0 {
		err = fmt.Errorf("金额不能小于零")
	} else if amount == 0 && coins == 0 {
		err = fmt.Errorf("金额不能为零")
	}

	if err != nil {
		return err
	}

	//---------------------------------------

	payinfo := &Pay{
		When:    time.Now().UTC(),
		PayFrom: u.name,
		PayTo:   who,
		Amount:  amount,
		Coins:   coins,
	}

	ctx, err := u.createProposal(payinfo, ProposalType_Pay)

	if err != nil {
		return err
	}

	//---------------------------------------

	u.proposal_chnl <- ctx

	ack := <-ctx.reply

	close(ctx.reply)

	return ack
}
