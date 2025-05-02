package vars

import "github.com/ekival-labs/ekival-canvas/utility"

var (
	toLovelace int64 = 1000000

	TEKITokenName string = "tMoney"
	TEKIPolicyId  string = "6f56a63701536ff2b5a8fc53c4a72d02491c1ea0c1e863ce37a1ffda"

	OrderId       string = "AwesomeID_TT_RR_1"
	// OrderAmount   int64  = 197 * toLovelace
	// MakerAddress  string = "addr_test1qr43kyyys0sg0d8khhjq4y348zuc4mnzln836hxfrqjasx4rqghlrj99l5vrdmyrtg6mhkyxa88kwq5yf225a4m9pkes2a4xz3"
	// TakerAddress  string = "addr_test1qqt85kcauy3uktlfldmhqn4dnn5vxgv8s97t729akje92f3z5uathwrlk6dwzy9j89lmsy2qp7evugtxmhf9pycw44rqj76x7l"
	// MakerDeadline int64  = 3600
	// TakerDeadline int64  = 1800

	Precision      int64 = 10
	CollateralPct  int64 = utility.ToFraction(10)
	MakerPct       int64 = utility.ToFraction(0.25)
	TakerPct       int64 = utility.ToFraction(0.75)
	CancelPct      int64 = utility.ToFraction(1)
	MinCollateral  int64 = 25 * toLovelace
	MakerMinFee    int64 = 1250000
	TakerMinFee    int64 = 3750000
	CancelMinFee   int64 = 3 * toLovelace
	MinOrderAmount int64 = 10 * toLovelace
	OrderThreshold int64 = 500 * toLovelace
	CancelPenalty  int64 = 2 * CancelMinFee
	AdaCollateral  int64 = 3 * toLovelace
)
