package constants

import (
	"github.com/Salvionied/apollo/serialization/PlutusData"
	"github.com/Salvionied/apollo/serialization/Redeemer"
)

var (
	INDEX_ONE_MINT_REDEEMER *Redeemer.Redeemer = &Redeemer.Redeemer{
		Tag:   Redeemer.MINT,
		Index: 0,
		Data: PlutusData.PlutusData{
			PlutusDataType: PlutusData.PlutusArray,
			TagNr:          INDEX_ONE,
			Value:          PlutusData.PlutusDefArray{},
		},
		// ExUnits: Redeemer.ExecutionUnits{
		// 	Mem:   450_000,
		// 	Steps: 200_000_000,
		// },
	}

	INDEX_TWO_MINT_REDEEMER *Redeemer.Redeemer = &Redeemer.Redeemer{
		Tag:   Redeemer.MINT,
		Index: 0,
		Data: PlutusData.PlutusData{
			PlutusDataType: PlutusData.PlutusArray,
			TagNr:          INDEX_TWO,
			Value:          PlutusData.PlutusDefArray{},
		},
	}

	INDEX_ONE_SPEND_REDEEMER *Redeemer.Redeemer = &Redeemer.Redeemer{
		Tag:   Redeemer.SPEND,
		Index: 0,
		Data: PlutusData.PlutusData{
			PlutusDataType: PlutusData.PlutusArray,
			TagNr:          INDEX_ONE,
			Value:          PlutusData.PlutusDefArray{},
		},
	}

	INDEX_TWO_SPEND_REDEEMER *Redeemer.Redeemer = &Redeemer.Redeemer{
		Tag:   Redeemer.SPEND,
		Index: 0,
		Data: PlutusData.PlutusData{
			PlutusDataType: PlutusData.PlutusArray,
			TagNr:          INDEX_TWO,
			Value:          PlutusData.PlutusDefArray{},
		},
	}

	INDEX_THREE_SPEND_REDEEMER *Redeemer.Redeemer = &Redeemer.Redeemer{
		Tag:   Redeemer.SPEND,
		Index: 0,
		Data: PlutusData.PlutusData{
			PlutusDataType: PlutusData.PlutusArray,
			TagNr:          INDEX_THREE,
			Value:          PlutusData.PlutusDefArray{},
		},
	}

	INDEX_FOUR_SPEND_REDEEMER *Redeemer.Redeemer = &Redeemer.Redeemer{
		Tag:   Redeemer.SPEND,
		Index: 0,
		Data: PlutusData.PlutusData{
			PlutusDataType: PlutusData.PlutusArray,
			TagNr:          INDEX_FOUR,
			Value:          PlutusData.PlutusDefArray{},
		},
	}

	INDEX_FIVE_SPEND_REDEEMER *Redeemer.Redeemer = &Redeemer.Redeemer{
		Tag:   Redeemer.SPEND,
		Index: 0,
		Data: PlutusData.PlutusData{
			PlutusDataType: PlutusData.PlutusArray,
			TagNr:          INDEX_FIVE,
			Value:          PlutusData.PlutusDefArray{},
		},
	}

	INDEX_SIX_SPEND_REDEEMER *Redeemer.Redeemer = &Redeemer.Redeemer{
		Tag:   Redeemer.SPEND,
		Index: 0,
		Data: PlutusData.PlutusData{
			PlutusDataType: PlutusData.PlutusArray,
			TagNr:          INDEX_SIX,
			Value:          PlutusData.PlutusDefArray{},
		},
	}

	EXTEND_DEADLINES_WITH_NO_PENALTY_REDEEMER *Redeemer.Redeemer = &Redeemer.Redeemer{
		Tag:   Redeemer.SPEND,
		Index: 0,
		Data: PlutusData.PlutusData{
			PlutusDataType: PlutusData.PlutusArray,
			TagNr:          INDEX_SEVEN,
			Value: PlutusData.PlutusDefArray{
				PlutusData.PlutusData{
					TagNr: 0,
					Value: PlutusData.PlutusData{
						PlutusDataType: PlutusData.PlutusArray,
						TagNr:          INDEX_ONE,
						Value:          PlutusData.PlutusDefArray{},
					},
				},
			},
		},
	}
)
