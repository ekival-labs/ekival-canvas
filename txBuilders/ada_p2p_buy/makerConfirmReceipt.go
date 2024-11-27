package ada_p2p_buy

import (
	"ekival-canvas/config"
	"ekival-canvas/constants"
	"ekival-canvas/model"
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

// makerFee := utility.CalculateFee(database.Precision, database.OrderAmount, database.OrderThreshold, database.MakerPct, database.MakerMinFee)
// takerFee := utility.CalculateFee(database.Precision, database.OrderAmount, database.OrderThreshold, database.TakerPct, database.TakerMinFee)
// collateralAmount := utility.CalculateFee(database.Precision, database.OrderAmount, database.OrderThreshold, database.CollateralPct, database.MinCollateral)

// orderInfo := &viewmodel.Order{
// 	OrderInfo: viewmodel.OrderInfo{
// 		OrderId:            database.OrderId,
// 		OrderAmount:        database.OrderAmount,
// 		MakerAddress:       ma,
// 		TakerAddress:       ta,
// 	},
// 	OrderTxInfo: viewmodel.OrderTxInfo{
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

// cborString, txHash, err := ada_p2p_buy.MakerConfirmReceipt(orderInfo, config.GetAda_P2PBuyAdminWallet())

func MakerConfirmReceipt(order *model.Order, treasury *model.TreasuryInfo, adminWallet *config.AdminWallet) (string, string, error) {

	// defer func() {
	// 	if err := recover(); err != nil {
	// 		log.Printf("Panic occurred: %v", err)
	// 		return
	// 	}
	// }()

	apolloBE := apollo.New(&config.CHAIN_CTX)
	apolloBE = apolloBE.SetWalletFromBech32(order.OrderInfo.MakerAddress.String())

	payToMakerAmount := order.OrderTxInfo.CollateralAmount
	payToTakerAmount := order.OrderInfo.OrderAmount + order.OrderTxInfo.CollateralAmount
	ekivalFee := order.OrderTxInfo.MakerFee + order.OrderTxInfo.TakerFee

	collateralUtxo := config.CHAIN_CTX.GetUtxoFromRef(order.OrderTxInfo.CollateralUtxo.TxID, order.OrderTxInfo.CollateralUtxo.TxIDIndex)
	orderUTxO := config.CHAIN_CTX.GetUtxoFromRef(order.OrderTxInfo.OrderUtxo.TxID, order.OrderTxInfo.OrderUtxo.TxIDIndex)
	userUtxos, err := utility.GetUserUTxOs(order.OrderTxInfo.UserUtxos)
	if err != nil {
		log.Println(err)
		return "", "", err
	}

	lastSlot := config.CHAIN_CTX.LastBlockSlot()

	apolloBE, err = apolloBE.
		SetChangeAddress(order.OrderTxInfo.ChangeAddress).
		AddCollateral(*collateralUtxo).
		AddLoadedUTxOs(userUtxos...).
		CollectFrom(*orderUTxO, *constants.INDEX_THREE_SPEND_REDEEMER).
		MintAssetsWithRedeemer(
			apollo.Unit{
				PolicyId: order.OrderTxInfo.StateTokenPolicyId,
				Name:     order.OrderInfo.OrderId,
				Quantity: int(-1),
			},
			*constants.INDEX_TWO_MINT_REDEEMER,
		).
		AddReferenceInput(
			order.OrderTxInfo.StateTokenRefUtxo.TxID,
			order.OrderTxInfo.StateTokenRefUtxo.TxIDIndex,
		).
		AddReferenceInput(
			order.OrderTxInfo.EscrowContractRefUtxo.TxID,
			order.OrderTxInfo.EscrowContractRefUtxo.TxIDIndex,
		).
		PayToContract(
			treasury.Address, treasury.Datum, int(ekivalFee), true,
		).
		PayToAddress(order.OrderInfo.MakerAddress, int(payToMakerAmount)).
		PayToAddress(order.OrderInfo.TakerAddress, int(payToTakerAmount)).
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

	if len(config.CHAIN_CTX.EvaluateTx(txByte)) == 0 {
		return "", "", fmt.Errorf("transaction evaluation failed")
	}

	return Utils.ToCbor(tx), hex.EncodeToString(txHash), nil

}
