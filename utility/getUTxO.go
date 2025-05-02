package utility

import (
	"errors"

	"github.com/ekival-labs/ekival-canvas/config"
	"github.com/ekival-labs/ekival-canvas/model"

	"github.com/Salvionied/apollo/serialization/UTxO"
)

func GetUserUTxOs(userUTxOs []model.EUTxO) ([]UTxO.UTxO, error) {
	var totalUTxOs []UTxO.UTxO
	for _, userUTxO := range userUTxOs {
		utxo, err := config.CHAIN_CTX.GetUtxoFromRef(userUTxO.TxID, userUTxO.TxIDIndex)
		if err != nil {
			return nil, err
		}
		if utxo == nil {
			return nil, errors.New("error at getUserUTxOs: UTxO not found")
		}
		totalUTxOs = append(totalUTxOs, *utxo)
	}
	return totalUTxOs, nil
}
