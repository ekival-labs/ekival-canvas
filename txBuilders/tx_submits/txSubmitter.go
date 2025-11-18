package tx_submits

import (
	"encoding/hex"
	"fmt"

	"ekival-canvas/config"
	"ekival-canvas/viewmodel"

	"github.com/Salvionied/apollo"
	fiberLogger "github.com/gofiber/fiber/v2/log"
)

func TxSubmitter(txResponse *viewmodel.TxResponse) (string, error) {
	defer func() {
		if err := recover(); err != nil {
			fiberLogger.Panicf("Panic occurred: %v", err)
		}
	}()

	apolloBE := apollo.New(&config.CHAIN_CTX)
	apolloBE, err := apolloBE.LoadTxCbor(txResponse.TxCBOR)
	if err != nil {
		fiberLogger.Errorf("Failed to load Tx CBOR: %v", err)
		return "", err
	}
	// txBytes1, err := apolloBE.GetTx().Bytes()
	// if err != nil {
	// 	fiberLogger.Error(err)
	// }
	// fmt.Println("Tx CBOR1:", hex.EncodeToString(txBytes1))

	// apolloBE, err = apolloBE.Complete()
	// if err != nil {
	// 	fiberLogger.Error(err)
	// }

	txHash, err := apolloBE.Submit()
	if err != nil {
		fiberLogger.Errorf("Submit Error: %T - %+v", err, err)
		return "", fmt.Errorf("apollo submit failed: %w", err)
	}

	fiberLogger.Debugf("TxID: %s", hex.EncodeToString(txHash.Payload))
	return "Tx successfully submitted", nil
}
