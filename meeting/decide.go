package meeting

type (
	DecideCore struct {
		Proposal
		Who   string // 议员名称
		Error string // 错误描述
	}

	Decide struct {
		DecideCore
		Signature
	}
)
