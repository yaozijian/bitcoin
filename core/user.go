package core

import (
	"fmt"

	"cryptology"
	"meeting"
	"router"

	"golang.org/x/net/context"

	log "github.com/cihub/seelog"
)

type (
	User interface {
		Register(who string, amount, coins int64) error
		PayTo(who string, amount, coins int64) error
		Show()
		OnBroadCast(router.BroadCastData)
	}

	user struct {
		ca      cryptology.CA
		public  cryptology.PublicKey
		private cryptology.PrivateKey

		router router.Router

		name    string // 名称
		onlines int

		broadcast_chnl chan router.BroadCastData // 接收广播
		proposal_chnl  chan *proposalContext     // 接收用户请求

		proposals_map map[string]*proposalContext // 表决中的议案: 议案摘要(hex编码) -->

		accounts_map map[string]*UserRecord // 用户
		bills_map    map[string]*BillRecord // 账单: 账单摘要(hex编码) --> 账单

		last_bill_digest []byte
	}

	proposalContext struct {
		proposal meeting.Proposal // 议案

		reply  chan error // 发送表决结果
		agree  int        // 赞成数
		oppose int        // 反对数

		onAgree func() // 成功后要执行的操作
		onFail  func() // 失败后要执行的操作

		ctx    context.Context
		cancel func()
	}
)

func NewUser(name, ca, cert, key string, r router.Router) (User, error) {

	u := &user{
		name:   name,
		router: r,

		broadcast_chnl: make(chan router.BroadCastData, 1000),
		proposal_chnl:  make(chan *proposalContext, 1000),

		proposals_map: make(map[string]*proposalContext),
		accounts_map:  make(map[string]*UserRecord),
		bills_map:     make(map[string]*BillRecord),

		last_bill_digest: FirstBillDigest(),
	}

	u.ca, _ = cryptology.LoadCAFromFile(ca)
	u.public, _ = cryptology.LoadPublicKeyFromFile(cert)
	u.private, _ = cryptology.LoadPrivateKey(key)

	if u.ca == nil || u.public == nil || u.private == nil {
		return nil, fmt.Errorf("证书或者密钥文件不正确")
	} else {
		go u.run()
		return u, nil
	}
}

func (u *user) run() {

	u.router.AddWatcher(u)

	for {
		select {
		case proposal := <-u.proposal_chnl:
			if err := u.doBroadCast(&proposal.proposal, BroadCastType_Proposal); err == nil {
				// 发起议案的用户需要首先添加议案,为了得到reply字段
				reqctx := u.add2ProposalMap(&proposal.proposal, proposal.proposal.ProposalCore.Type == ProposalType_Pay)
				reqctx.reply = proposal.reply
			} else if proposal.reply != nil {
				proposal.reply <- err
			}
		case data := <-u.broadcast_chnl:
			u.onBroadCast(data)
		}
	}
}

func (u *user) Show() {
	log.Infof("----- me=%v 用户记录begin -----", u.name)
	for _, k := range u.accounts_map {
		log.Infof("用户%v: u=%v a=%v c=%v", u.name, k.Name, k.Amount, k.Coins)
	}
	log.Infof("----- me=%v 用户记录end -----\n\n", u.name)
}
