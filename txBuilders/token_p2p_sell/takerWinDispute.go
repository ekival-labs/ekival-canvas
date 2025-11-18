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

// makerFee := utility.CalculateFee(database.Precision, database.OrderAmount, database.OrderThreshold, database.MakerPct, database.MakerMinFee)
// takerFee := utility.CalculateFee(database.Precision, database.OrderAmount, database.OrderThreshold, database.TakerPct, database.TakerMinFee)
// collateralAmount := utility.CalculateFee(database.Precision, database.OrderAmount, database.OrderThreshold, database.CollateralPct, database.MinCollateral)
// disputeFee := utility.CalculateFee(database.Precision, database.OrderAmount, database.OrderThreshold, database.DisputePct, database.DisputeMinFee)

// orderInfo := &viewmodel.Order{
// 	OrderInfo: viewmodel.OrderInfo{
// 		TradeTokenName:     database.EKI_TokenName,
// 		TradeTokenPolicyId: database.EKI_PolicyId,
// 		OrderId:            database.OrderId,
// 		OrderAmount:        database.OrderAmount,
// 		MakerAddress:       ma,
// 		TakerAddress:       ta,
// 	},
// 	BrokerageInfo: viewmodel.BrokerageInfo{
// 		AdaCollateral:  database.AdaCollateral,
// 	},
// 	OrderTxInfo: viewmodel.OrderTxInfo{
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
// 		TakerFee:         takerFee,
// 		DisputeFee:       disputeFee,
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

// takerWinDisputeRedeemerUnmarshaled := viewmodel.TakerWinDisputeRedeemer{
// 	DisputeInfo: viewmodel.DisputeInfo{
// 		DisputePct:     database.DisputePct,
// 		DisputeMinFee:  database.DisputeMinFee,
// 		DisputePenalty: database.DisputePenalty,
// 	},
// }

// cborString, txHash, err := token_p2p_sell.TakerWinDispute(orderInfo, treasuryInfo, takerWinDisputeRedeemerUnmarshaled, config.GetEki_P2PSellAdminWallet())

func TakerWinDispute(order *model.Order, treasury *model.TreasuryInfo, disputeRedeemer *model.MakerWinDisputeRedeemer, adminWallet *config.Wallet) (string, string, error) {

	defer func() {
		if err := recover(); err != nil {
			log.Printf("Panic occurred: %v", err)
			return
		}
	}()

	apolloBE := apollo.New(&config.CHAIN_CTX)
	apolloBE = apolloBE.SetWalletFromBech32(order.OrderInfo.TakerAddress.String())

	takerWinDisputeRedeemerMarshaled, err := plutusEncoder.MarshalPlutus(*disputeRedeemer)
	if err != nil {
		log.Println(err)
		return "", "", err
	}

	takerWinDisputeRedeemer := Redeemer.Redeemer{
		Tag:   Redeemer.SPEND,
		Index: 0,
		Data:  *takerWinDisputeRedeemerMarshaled,
	}

	payToMakerAmount := order.OrderTxInfo.CollateralAmount + order.OrderTxInfo.MakerFee - disputeRedeemer.DisputeInfo.DisputePenalty - order.OrderTxInfo.DisputeFee
	payToTakerAmount := order.OrderTxInfo.CollateralAmount + order.OrderTxInfo.TakerFee + order.OrderInfo.OrderAmount + disputeRedeemer.DisputeInfo.DisputePenalty

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
		CollectFrom(*orderUTxO, takerWinDisputeRedeemer).
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
		PayToContract(
			treasury.Address,
			treasury.Datum,
			int(order.BrokerageInfo.AdaCollateral),
			true,
			apollo.Unit{
				PolicyId: order.OrderInfo.TradeTokenPolicyId,
				Name:     order.OrderInfo.TradeTokenName,
				Quantity: int(order.OrderTxInfo.DisputeFee),
			},
		).
		PayToAddress(
			order.OrderInfo.MakerAddress,
			int(order.BrokerageInfo.AdaCollateral),
			apollo.Unit{
				PolicyId: order.OrderInfo.TradeTokenPolicyId,
				Name:     order.OrderInfo.TradeTokenName,
				Quantity: int(payToMakerAmount),
			},
		).
		PayToAddress(
			order.OrderInfo.TakerAddress,
			int(order.BrokerageInfo.AdaCollateral*2),
			apollo.Unit{
				PolicyId: order.OrderInfo.TradeTokenPolicyId,
				Name:     order.OrderInfo.TradeTokenName,
				Quantity: int(payToTakerAmount),
			}).
		AddRequiredSigner(adminWallet.AdminPKH).
		AddRequiredSigner(serialization.PubKeyHash(order.OrderInfo.TakerAddress.PaymentPart)).
		SetTtl(int64(lastSlot) + 300).
		SetValidityStart(int64(lastSlot)).
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
