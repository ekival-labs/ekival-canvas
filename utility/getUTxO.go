package utility

import (
	"errors"

	"ekival-canvas/config"
	"ekival-canvas/model"

	"github.com/Salvionied/apollo/serialization/UTxO"
)

func GetUserUTxOs(userUTxOs []model.EUTxO) ([]UTxO.UTxO, error) {
	var totalUTxOs []UTxO.UTxO
	for _, userUTxO := range userUTxOs {
		utxo := config.CHAIN_CTX.GetUtxoFromRef(userUTxO.TxID, userUTxO.TxIDIndex)
		if utxo == nil {
			return nil, errors.New("error at getUserUTxOs: UTxO not found")
		}
		totalUTxOs = append(totalUTxOs, *utxo)
	}
	return totalUTxOs, nil
}
