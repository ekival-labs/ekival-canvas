package model

import (
	"github.com/Salvionied/apollo/serialization/Address"
)

type AdaOfferTxInfo struct {
	MakerAddress      Address.Address
	ChangeAddress     Address.Address
	UserUtxos         []EUTxO
	CollateralUtxo    EUTxO
	EkivalFeeLovelace int
}

type TokenOfferTxInfo struct {
	TradeTokenName     string
	TradeTokenPolicyId string
	EkivalFeeToken     int
	MakerAddress       Address.Address
	ChangeAddress      Address.Address
	UserUtxos          []EUTxO
	CollateralUtxo     EUTxO
	EkivalFeeLovelace  int
}
