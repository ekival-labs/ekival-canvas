package tests

import (
	"ekival-canvas/config"
	"ekival-canvas/utility"
)

var (
	adminWallet = config.SetWallet("garment shadow into lab truck quiz file warm sheriff marriage voice icon thunder iron early jungle bird dash material strike desk mango deer letter")
	makerWallet = config.SetWallet("brief machine better party office phone thunder way swift plug cave praise total blue inject copy provide limb flower decrease matrix wreck host whale")
	takerWallet = config.SetWallet("bone miracle mother grocery rabbit decorate rain moment print empty harbor cinnamon resource desert roof attitude suspect cupboard allow hunt inhale praise sausage mom")

	COLLATERAL_PCT int64 = utility.ToFraction(10)
	MAKER_PCT      int64 = utility.ToFraction(0.25)
	TAKER_PCT      int64 = utility.ToFraction(0.75)
	CANCEL_PCT     int64 = utility.ToFraction(1)
	DISPUTE_PCT    int64 = utility.ToFraction(2)
)
