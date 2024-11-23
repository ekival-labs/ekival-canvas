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
	fiberLogger "github.com/gofiber/fiber/v2/log"
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

// makerFee := utility.CalculateFee(database.Precision, database.OrderAmount, database.OrderThreshold, database.MakerPct, database.MakerMinFee)
// takerFee := utility.CalculateFee(database.Precision, database.OrderAmount, database.OrderThreshold, database.TakerPct, database.TakerMinFee)
// collateralAmount := utility.CalculateFee(database.Precision, database.OrderAmount, database.OrderThreshold, database.CollateralPct, database.MinCollateral)

// orderInfo := &viewmodel.Order{
// 	OrderInfo: viewmodel.OrderInfo{
// 		OrderId:            database.OrderId,
// 		OrderAmount:        body.OrderAmount,
// 		MakerAddress:       ma,
// 		TakerAddress:       ta,
// 		MakerDeadline:      database.MakerDeadline,
// 		TakerDeadline:      database.TakerDeadline,
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
// 	TradeState: constants.COMMITTED_ORDER_STATUS,
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
// 		TakerFee:         takerFee,
// 		CollateralAmount: collateralAmount,
// 		ChangeAddress:    changeAddress,
// 		UserUtxos:        body.UserUTxOs,
// 		CollateralUtxo:   body.CollateralUTxO,
// 		OrderUtxo : model.EUTxO{
// 			TxID:      database.OrderUtxo.TxID,
// 			TxIDIndex: database.OrderUtxo.TxIDIndex,
// 		}
// 	},
// }

// cborString, txHash, err := ada_p2p_buy.TakerCommitToOrder(orderInfo, config.GetAda_P2PBuyAdminWallet())
func TakerCommitToOrder(order *model.Order, adminWallet *config.AdminWallet) (string, string, error) {

	// defer func() {
	// 	if err := recover(); err != nil {
	// 		log.Printf("Panic occurred: %v", err)
	// 		return
	// 	}
	// }()

	apolloBE := apollo.New(&config.BFC)
	apolloBE = apolloBE.SetWalletFromBech32(order.OrderInfo.TakerAddress.String())

	makerCommittedAmount := order.OrderInfo.OrderAmount + order.OrderTxInfo.MakerFee + order.OrderTxInfo.CollateralAmount
	takerCommittingAmount := order.OrderTxInfo.TakerFee + order.OrderTxInfo.CollateralAmount
	totalAmount := makerCommittedAmount + takerCommittingAmount

	orderDatumMarshaled, err := plutusEncoder.MarshalPlutus(*order)
	if err != nil {
		log.Println(err)
		return "", "", err
	}

	collateralUtxo := config.BFC.GetUtxoFromRef(order.OrderTxInfo.CollateralUtxo.TxID, order.OrderTxInfo.CollateralUtxo.TxIDIndex)
	orderUTxO := config.BFC.GetUtxoFromRef(order.OrderTxInfo.OrderUtxo.TxID, order.OrderTxInfo.OrderUtxo.TxIDIndex)
	userUtxos, err := utility.GetUserUTxOs(order.OrderTxInfo.UserUtxos)
	if err != nil {
		log.Println(err)
		return "", "", err
	}

	lastSlot := config.BFC.LastBlockSlot()

	apolloBE, err = apolloBE.
		SetChangeAddress(order.OrderTxInfo.ChangeAddress).
		AddCollateral(*collateralUtxo).
		AddLoadedUTxOs(userUtxos...).
		CollectFrom(*orderUTxO, *constants.INDEX_TWO_SPEND_REDEEMER).
		AddReferenceInput(
			order.OrderTxInfo.EscrowContractRefUtxo.TxID,
			order.OrderTxInfo.EscrowContractRefUtxo.TxIDIndex,
		).
		PayToContract(
			order.OrderTxInfo.EscrowContractAddress, orderDatumMarshaled, int(totalAmount), true, apollo.Unit{
				PolicyId: order.OrderTxInfo.StateTokenPolicyId,
				Name:     order.OrderInfo.OrderId,
				Quantity: int(1),
			},
		).
		AddRequiredSigner(adminWallet.AdminPKH).
		AddRequiredSigner(serialization.PubKeyHash(order.OrderInfo.TakerAddress.PaymentPart)).
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

	fiberLogger.Debug("TxID:", hex.EncodeToString(txHash))
	fiberLogger.Debug("Tx CBOR:", Utils.ToCbor(tx))
	fiberLogger.Debug("EvaluateTx:", config.BFC.EvaluateTx(txByte))

	if len(config.BFC.EvaluateTx(txByte)) == 0 {
		return "", "", fmt.Errorf("transaction evaluation failed")
	}

	return Utils.ToCbor(tx), hex.EncodeToString(txHash), nil

}
