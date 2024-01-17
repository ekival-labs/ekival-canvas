package viewmodel

import (
	"github.com/Salvionied/apollo/serialization/Address"
	"github.com/Salvionied/apollo/serialization/PlutusData"
)

type TreasuryInfo struct {
	Address Address.Address
	Datum   *PlutusData.PlutusData
}
