package utility

import "github.com/Salvionied/apollo/serialization/PlutusData"

func CreateSimpleDatum(tag_num uint64, fee_type int) *PlutusData.PlutusData {
	return &PlutusData.PlutusData{
		PlutusDataType: PlutusData.PlutusArray,
		TagNr:          tag_num,
		Value: PlutusData.PlutusDefArray{
			PlutusData.PlutusData{
				TagNr:          0,
				Value:          fee_type,
				PlutusDataType: PlutusData.PlutusInt,
			},
		},
	}
}
