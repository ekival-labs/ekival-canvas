package token_p2p_buy

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

// 	changeAddress, err := Address.DecodeAddress(body.ChangeAddress)
// 	if err != nil {
// 		log.Println(err)
// 		return err
// 	}

// makerFee := utility.CalculateFee(database.Precision, database.OrderAmount, database.OrderThreshold, database.MakerPct, database.MakerMinFee)
// collateralAmount := utility.CalculateFee(database.Precision, database.OrderAmount, database.OrderThreshold, database.CollateralPct, database.MinCollateral)

// orderInfo := &viewmodel.Order{
// 	OrderInfo: viewmodel.OrderInfo{
// 		TradeTokenName:     database.EKI_TokenName,
// 		TradeTokenPolicyId: database.EKI_PolicyId,
// 		OrderId:            database.OrderId,
// 		OrderAmount:        database.OrderAmount,
// 		MakerAddress:       ma,
// 	},
// 	BrokerageInfo: viewmodel.BrokerageInfo{
// 		AdaCollateral:  database.AdaCollateral,
// 	},
// 	OrderTxInfo: viewmodel.OrderTxInfo{
// 		EscrowContractRefUtxo: model.EUTxO{
// 			TxID:      database.Eki_P2PBuyEscrow.RefTxID,
// 			TxIDIndex: database.Eki_P2PBuyEscrow.RefTxIDIndex,
// 		},
// 		StateTokenPolicyId: database.Eki_BST.PolicyID,
// 		StateTokenRefUtxo: model.EUTxO{
// 			TxID:      database.Eki_BST.RefTxID,
// 			TxIDIndex: database.Eki_BST.RefTxIDIndex,
// 		},
// 		MakerFee:         makerFee,
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

// cborString, txHash, err := token_p2p_buy.MakerCancelUncommittedOrder(orderInfo, config.GetEki_P2PBuyAdminWallet())

func MakerCancelUncommittedOrder(order *model.Order, adminWallet *config.Wallet) (string, string, error) {

	defer func() {
		if err := recover(); err != nil {
			log.Printf("Panic occurred: %v", err)
			return
		}
	}()

	apolloBE := apollo.New(&config.CHAIN_CTX)
	apolloBE = apolloBE.SetWalletFromBech32(order.OrderInfo.MakerAddress.String())

	makerCommittedAmount := order.OrderInfo.OrderAmount + order.OrderTxInfo.MakerFee + order.OrderTxInfo.CollateralAmount

	collateralUtxo, err := config.CHAIN_CTX.GetUtxoFromRef(order.OrderTxInfo.CollateralUtxo.TxID, order.OrderTxInfo.CollateralUtxo.TxIDIndex)
	if err != nil {
		log.Println(err)
		return "", "", err
	}
	orderUTxO, err := config.CHAIN_CTX.GetUtxoFromRef(order.OrderTxInfo.OrderUtxo.TxID, order.OrderTxInfo.OrderUtxo.TxIDIndex)
	if err != nil {
		log.Println(err)
		return "", "", err
	}
	userUtxos, err := utility.GetUserUTxOs(order.OrderTxInfo.UserUtxos)
	if err != nil {
		log.Println(err)
		return "", "", err
	}

	lastSlot, err := config.CHAIN_CTX.LastBlockSlot()
	if err != nil {
		log.Println(err)
		return "", "", err
	}

	apolloBE, err = apolloBE.
		SetChangeAddress(order.OrderTxInfo.ChangeAddress).
		AddCollateral(*collateralUtxo).
		AddLoadedUTxOs(userUtxos...).
		CollectFrom(*orderUTxO, *constants.INDEX_ONE_SPEND_REDEEMER).
		MintAssetsWithRedeemer(
			apollo.Unit{
				PolicyId: order.OrderTxInfo.StateTokenPolicyId,
				Name:     order.OrderInfo.OrderId,
				Quantity: int(-1),
			},
			*constants.INDEX_TWO_MINT_REDEEMER,
		).
		AddReferenceInputV3(
			order.OrderTxInfo.StateTokenRefUtxo.TxID,
			order.OrderTxInfo.StateTokenRefUtxo.TxIDIndex,
		).
		AddReferenceInputV3(
			order.OrderTxInfo.EscrowContractRefUtxo.TxID,
			order.OrderTxInfo.EscrowContractRefUtxo.TxIDIndex,
		).
		PayToAddress(
			order.OrderInfo.MakerAddress,
			int(order.BrokerageInfo.AdaCollateral*2),
			apollo.Unit{
				PolicyId: order.OrderInfo.TradeTokenPolicyId,
				Name:     order.OrderInfo.TradeTokenName,
				Quantity: int(makerCommittedAmount),
			}).
		AddRequiredSigner(serialization.PubKeyHash(adminWallet.AdminPKH)).
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

	evaluationResult, err := config.CHAIN_CTX.EvaluateTx(txByte)
	if err != nil {
		log.Println(err)
		return "", "", err
	}
	if len(evaluationResult) == 0 {
		return "", "", fmt.Errorf("transaction evaluation failed")
	}

	cbor, err := Utils.ToCbor(tx)
	if err != nil {
		log.Println(err)
		return "", "", err
	}
	return cbor, hex.EncodeToString(txHash), nil

}
