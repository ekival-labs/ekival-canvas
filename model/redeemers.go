package model

type ExtendingPenalty struct {
	_       struct{} `plutusType:"DefList" plutusConstr:"0"`
	Penalty int64    `plutusType:"Int"`
}

type ExtendDeadlineWithPenaltyRedeemer struct {
	_                struct{} `plutusType:"DefList" plutusConstr:"7"`
	ExtendingPenalty ExtendingPenalty
}

type DisputeInfo struct {
	_              struct{} `plutusType:"DefList" plutusConstr:"1"`
	DisputePct     int64    `plutusType:"Int"`
	DisputeMinFee  int64    `plutusType:"Int"`
	DisputePenalty int64    `plutusType:"Int"`
}

type MakerWinDisputeRedeemer struct {
	_           struct{} `plutusType:"DefList" plutusConstr:"8"`
	DisputeInfo DisputeInfo
}

type TakerWinDisputeRedeemer struct {
	_           struct{} `plutusType:"DefList" plutusConstr:"9"`
	DisputeInfo DisputeInfo
}
