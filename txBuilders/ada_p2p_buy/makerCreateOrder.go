package ada_p2p_buy

import (
	"ekival-canvas/config"
	"ekival-canvas/constants"
	"ekival-canvas/model"
	"ekival-canvas/plutusEncoder"
	"ekival-canvas/utility"
	"encoding/hex"
	"fmt"
	"log"

	"github.com/Salvionied/apollo"
	"github.com/Salvionied/apollo/serialization"
	"github.com/Salvionied/apollo/txBuilding/Utils"
)

// cfg := config.GetGlobalConfig()

// 	ma, err := Address.DecodeAddress(database.MakerAddress)
// 	if err != nil {
// 		log.Println(err)
// 		return err
// 	}

// 	ta, err := Address.DecodeAddress(database.TakerAddress)
// 	if err != nil {
// 		log.Println(err)
// 		return err
// 	}

// 	changeAddress, err := Address.DecodeAddress(body.ChangeAddress)
// 	if err != nil {
// 		log.Println(err)
// 		return err
// 	}

// 	escrowContractAddress, err := Address.DecodeAddress(database.Ada_P2PBuyEscrow.Address)
// 	if err != nil {
// 		log.Println(err)
// 		return err
// 	}

// 	currentTime := time.Now().UnixNano() / int64(time.Millisecond)
// // MakerDeadline is whole Order Time in Hours for example 5 days = 5 * 24 * 60 * 60 * 1000
// 	var md int64 = currentTime + (MakerDeadline * 1000)
// // TakerDeadline less than MakerDeadline in Hours for example 4 day = 4 * 24 * 60 * 60 * 1000
// 	var td int64 = currentTime + (TakerDeadline * 1000)

// makerFee := utility.CalculateFee(database.Precision, database.OrderAmount, database.OrderThreshold, database.MakerPct, database.MakerMinFee)
// collateralAmount := utility.CalculateFee(database.Precision, database.OrderAmount, database.OrderThreshold, database.CollateralPct, database.MinCollateral)

// orderInfo := &viewmodel.Order{
// 	OrderInfo: viewmodel.OrderInfo{
// 		OrderId:            database.OrderId,
// 		OrderAmount:        body.OrderAmount,
// 		MakerAddress:       ma,
// 		TakerAddress:       ta,
// 		MakerDeadline:      md,
// 		TakerDeadline:      td,
// 	},
// 	BrokerageInfo: viewmodel.BrokerageInfo{
// 		Precision:      database.Precision,
// 		CollateralPct:  database.CollateralPct,
// 		MakerPct:       database.MakerPct,
// 		TakerPct:       database.TakerPct,
// 		CancelPct:      database.CancelPct,
// 		MinCollateral:  database.MinCollateral,
// 		MakerMinFee:    database.MakerMinFee,
// 		TakerMinFee:    database.TakerMinFee,
// 		CancelMinFee:   database.CancelMinFee,
// 		MinOrderAmount: database.MinOrderAmount,
// 		OrderThreshold: database.OrderThreshold,
// 		CancelPenalty:  database.CancelPenalty,
// 	},
// 	TradeState: constants.UNCOMMITTED_ORDER_STATUS,
// 	OrderTxInfo: viewmodel.OrderTxInfo{
// 		EscrowContractAddress: escrowContractAddress,
// 		EscrowContractRefUtxo: model.EUTxO{
// 			TxID:      database.Ada_P2PBuyEscrow.RefTxID,
// 			TxIDIndex: database.Ada_P2PBuyEscrow.RefTxIDIndex,
// 		},
// 		StateTokenPolicyId: database.Ada_BST.PolicyID,
// 		StateTokenRefUtxo: model.EUTxO{
// 			TxID:      database.Ada_BST.RefTxID,
// 			TxIDIndex: database.Ada_BST.RefTxIDIndex,
// 		},
// 		MakerFee:         makerFee,
// 		CollateralAmount: collateralAmount,
// 		ChangeAddress:    changeAddress,
// 		UserUtxos:        body.UserUTxOs,
// 		CollateralUtxo:   body.CollateralUTxO,
// 	},
// }

// cborString, txHash, err := ada_p2p_buy.MakerCreateOrder(orderInfo, config.GetAda_BSTAdminWallet())
func MakerCreateOrder(order *model.Order, adminWallet *config.AdminWallet) (string, string, error) {

	// defer func() {
	// 	if err := recover(); err != nil {
	// 		log.Printf("Panic occurred: %v", err)
	// 		return
	// 	}
	// }()

	apolloBE := apollo.New(&config.CHAIN_CTX)
	apolloBE = apolloBE.SetWalletFromBech32(order.OrderInfo.MakerAddress.String())

	makerCommittingAmount := order.OrderInfo.OrderAmount + order.OrderTxInfo.MakerFee + order.OrderTxInfo.CollateralAmount

	orderDatumMarshaled, err := plutusEncoder.MarshalPlutus(*order)
	if err != nil {
		log.Println(err)
		return "", "", err
	}

	userUtxos, err := utility.GetUserUTxOs(order.OrderTxInfo.UserUtxos)
	if err != nil {
		log.Println(err)
		return "", "", err
	}

	lastSlot := config.CHAIN_CTX.LastBlockSlot()

	collateralUtxo := config.CHAIN_CTX.GetUtxoFromRef(order.OrderTxInfo.CollateralUtxo.TxID, order.OrderTxInfo.CollateralUtxo.TxIDIndex)

	apolloBE, err = apolloBE.
		SetChangeAddress(order.OrderTxInfo.ChangeAddress).
		AddCollateral(*collateralUtxo).
		AddLoadedUTxOs(userUtxos...).
		MintAssetsWithRedeemer(
			apollo.Unit{
				PolicyId: order.OrderTxInfo.StateTokenPolicyId,
				Name:     order.OrderInfo.OrderId,
				Quantity: int(1),
			},
			*constants.INDEX_ONE_MINT_REDEEMER,
		).
		AddReferenceInput(
			order.OrderTxInfo.StateTokenRefUtxo.TxID,
			order.OrderTxInfo.StateTokenRefUtxo.TxIDIndex,
		).
		PayToContract(
			order.OrderTxInfo.EscrowContractAddress,
			orderDatumMarshaled,
			int(makerCommittingAmount),
			true,
			apollo.Unit{
				PolicyId: order.OrderTxInfo.StateTokenPolicyId,
				Name:     order.OrderInfo.OrderId,
				Quantity: int(1),
			},
		).
		AddRequiredSigner(adminWallet.AdminPKH).
		AddRequiredSigner(serialization.PubKeyHash(order.OrderInfo.MakerAddress.PaymentPart)).
		SetTtl(int64(lastSlot) + 300).
		Complete()

	if err != nil {
		log.Println(err)
		return "", "", err
	}

	apolloBE, err = apolloBE.SignWithSkey(adminWallet.AdminVkey, adminWallet.AdminSkey)
	if err != nil {
		log.Println(err)
		return "", "", err
	}

	tx := apolloBE.GetTx()
	txHash, err := tx.TransactionBody.Hash()
	if err != nil {
		log.Println(err)
		return "", "", err
	}

	txByte, err := tx.Bytes()
	if err != nil {
		log.Println(err)
		return "", "", err
	}

	fmt.Println("TX EVAL: ", config.CHAIN_CTX.EvaluateTx(txByte))
	fmt.Println("CBOR: ", Utils.ToCbor(tx))

	// if len(config.CHAIN_CTX.EvaluateTx(txByte)) == 0 {
	// 	return "", "", fmt.Errorf("transaction evaluation failed")
	// }

	return Utils.ToCbor(tx), hex.EncodeToString(txHash), nil

}
