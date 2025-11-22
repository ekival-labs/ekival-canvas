package token_p2p_sell

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
	"github.com/Salvionied/apollo/serialization/Redeemer"
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

// 	escrowContractAddress, err := Address.DecodeAddress(database.Eki_P2PSellEscrow.Address)
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
// 		AdaCollateral:  database.AdaCollateral,
// 	},
// // TradeState msu be same as the before state
// 	TradeState: database.TradeState,
// 	OrderTxInfo: viewmodel.OrderTxInfo{
// 		EscrowContractAddress: escrowContractAddress,
// 		EscrowContractRefUtxo: model.EUTxO{
// 			TxID:      database.Eki_P2PSellEscrow.RefTxID,
// 			TxIDIndex: database.Eki_P2PSellEscrow.RefTxIDIndex,
// 		},
// 		StateTokenPolicyId: database.Eki_SST.PolicyID,
// 		StateTokenRefUtxo: model.EUTxO{
// 			TxID:      database.Eki_SST.RefTxID,
// 			TxIDIndex: database.Eki_SST.RefTxIDIndex,
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

// treasuryAddress, err := Address.DecodeAddress(database.treasuryAddress)
// if err != nil {
// 	utility.ErrorWriter(w, 400, err.Error())
// 	return
// }
// treasuryInfo := &viewmodel.TreasuryInfo{
// 	Address: treasuryAddress,
// Datum:   utility.CreateSimpleDatum(constants.INDEX_ONE, constants.EKI_P2P_SELL_FEE_TYPE),
// }

// if userShouldPayForExtendingDeadline == true {
// 	extendPenalty = database.extendPenalty // Or 5_000_000 as fee
// } else {
// 	extendPenalty = 0
// }

// if order.TradeState == constants.REMIT_CONFIRMED_STATUS || order.TradeState == constants.COMMITTED_ORDER_STATUS {
// 	takerFee := utility.CalculateFee(database.Precision, database.OrderAmount, database.OrderThreshold, database.TakerPct, database.TakerMinFee)
// 	order.OrderTxInfo.TakerFee = takerFee
// }

// cborString, txHash, err := token_p2p_sell.TakerExtendDeadline(orderInfo, treasuryInfo, extendPenalty, config.GetEki_P2PSellAdminWallet())
func TakerExtendDeadline(order *model.Order, treasury *model.TreasuryInfo, extend_penalty int64, adminWallet *config.Wallet) (string, string, error) {

	defer func() {
		if err := recover(); err != nil {
			log.Printf("Panic occurred: %v", err)
			return
		}
	}()

	apolloBE := apollo.New(&config.CHAIN_CTX)
	apolloBE = apolloBE.SetWalletFromBech32(order.OrderInfo.TakerAddress.String())

	orderDatumMarshaled, err := plutusEncoder.MarshalPlutus(*order)
	if err != nil {
		log.Println(err)
		return "", "", err
	}

	var amount int64 = 0

	switch order.TradeState {
	case constants.UNCOMMITTED_ORDER_STATUS:
		amount = order.OrderTxInfo.MakerFee + order.OrderTxInfo.CollateralAmount
	default:
		makerCommittedAmount := order.OrderTxInfo.MakerFee + order.OrderTxInfo.CollateralAmount
		takerCommittedAmount := order.OrderInfo.OrderAmount + order.OrderTxInfo.TakerFee + order.OrderTxInfo.CollateralAmount
		amount = makerCommittedAmount + takerCommittedAmount
	}

	var extendDeadlinesRedeemer Redeemer.Redeemer = Redeemer.Redeemer{}

	switch extend_penalty {
	case 0:
		extendDeadlinesRedeemer = *constants.EXTEND_DEADLINES_WITH_NO_PENALTY_REDEEMER

	default:
		extendDeadlinesWithPenaltyRedeemerUnmarshaled := model.ExtendDeadlineWithPenaltyRedeemer{
			ExtendingPenalty: model.ExtendingPenalty{
				Penalty: extend_penalty,
			},
		}

		extendDeadlinesWithPenaltyRedeemerMarshaled, err := plutusEncoder.MarshalPlutus(extendDeadlinesWithPenaltyRedeemerUnmarshaled)
		if err != nil {
			log.Println(err)
			return "", "", err
		}

		extendDeadlinesRedeemer = Redeemer.Redeemer{
			Tag:   Redeemer.SPEND,
			Index: 0,
			Data:  *extendDeadlinesWithPenaltyRedeemerMarshaled,
		}

		apolloBE = apolloBE.PayToContract(
			treasury.Address, treasury.Datum, int(extend_penalty), true,
		)
	}

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
		CollectFrom(*orderUTxO, extendDeadlinesRedeemer).
		AddReferenceInputV3(
			order.OrderTxInfo.EscrowContractRefUtxo.TxID,
			order.OrderTxInfo.EscrowContractRefUtxo.TxIDIndex,
		).
		PayToContract(
			order.OrderTxInfo.EscrowContractAddress,
			orderDatumMarshaled,
			int(order.BrokerageInfo.AdaCollateral*4),
			true,
			apollo.Unit{
				PolicyId: order.OrderTxInfo.StateTokenPolicyId,
				Name:     order.OrderInfo.OrderId,
				Quantity: int(1),
			},
			apollo.Unit{
				PolicyId: order.OrderInfo.TradeTokenPolicyId,
				Name:     order.OrderInfo.TradeTokenName,
				Quantity: int(amount),
			},
		).
		AddRequiredSigner(adminWallet.PKH).
		AddRequiredSigner(serialization.PubKeyHash(order.OrderInfo.TakerAddress.PaymentPart)).
		SetTtl(int64(lastSlot) + 300).
		Complete()

	if err != nil {
		log.Println(err)
		return "", "", err
	}

	apolloBE, err = apolloBE.SignWithSkey(adminWallet.Vkey, adminWallet.Skey)
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
