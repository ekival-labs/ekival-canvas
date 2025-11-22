package utility

import (
	"encoding/hex"
	"errors"

	"ekival-canvas/config"
	"ekival-canvas/model"

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

func CreateEUTxOs(utxos []UTxO.UTxO) []model.EUTxO {
	var eUTxOs []model.EUTxO
	for _, utxo := range utxos {
		eUTxOs = append(eUTxOs, model.EUTxO{
			TxID:      hex.EncodeToString(utxo.Input.TransactionId),
			TxIDIndex: int(utxo.Input.Index),
		})
	}
	return eUTxOs
}
