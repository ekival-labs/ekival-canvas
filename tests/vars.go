package tests

import (
	"ekival-canvas/config"
	"ekival-canvas/utility"
)

var (
	orderUTxOTxId      string = "b8d70feea6bca4bc51e4478490753a5020d3e2aecc1d83b1b307543a47b6a6a6"
	orderUTxOTxIdIndex int    = 0

	orderId    string = "AwesomeID_TT_RR_1"
	toLovelace int64  = 1000000
	// orderAmount int64  = 719854678691
	orderAmount      int64 = 197 * toLovelace
	precision        int64 = 10
	collateralPct    int64 = utility.ToFraction(10)
	makerPct         int64 = utility.ToFraction(0.25)
	takerPct         int64 = utility.ToFraction(0.75)
	cancelPct        int64 = utility.ToFraction(1)
	minCollateral    int64 = 25 * toLovelace
	makerMinFee      int64 = 1250000
	takerMinFee      int64 = 3750000
	cancelMinFee     int64 = 3 * toLovelace
	minOrderAmount   int64 = 10 * toLovelace
	orderThreshold   int64 = 500 * toLovelace
	cancelPenalty    int64 = 2 * cancelMinFee
	adaCollateral    int64 = 3 * toLovelace
	extendingPenalty int64 = 5 * toLovelace
	disputePct       int64 = utility.ToFraction(2)
	disputeMinFee    int64 = 5 * toLovelace
	disputePenalty   int64 = 2 * disputeMinFee

	makerDeadline int64 = 7200
	takerDeadline int64 = 3600

	adminWallet = config.SetWallet("garment shadow into lab truck quiz file warm sheriff marriage voice icon thunder iron early jungle bird dash material strike desk mango deer letter")
	makerWallet = config.SetWallet("bone miracle mother grocery rabbit decorate rain moment print empty harbor cinnamon resource desert roof attitude suspect cupboard allow hunt inhale praise sausage mom")
)
