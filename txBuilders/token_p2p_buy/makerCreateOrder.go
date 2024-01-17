package token_p2p_buy

import (
	"ekival-canvas/config"
	"ekival-canvas/constants"
	"ekival-canvas/plutusEncoder"
	"ekival-canvas/utility"
	"ekival-canvas/viewmodel"
	"encoding/hex"
	"fmt"
	"log"

	"github.com/Salvionied/apollo"
	"github.com/Salvionied/apollo/serialization"
	"github.com/Salvionied/apollo/txBuilding/Utils"
)

func MakerCreateOrder(order *viewmodel.Order, adminWallet *config.AdminWallet) (string, string, error) {

	defer func() {
		if err := recover(); err != nil {
			log.Printf("Panic occurred: %v", err)
			return
		}
	}()

	apolloBE := apollo.New(&config.BFC)
	apolloBE = apolloBE.SetWalletFromBech32(order.OrderInfo.MakerAddress.String())

	makerFee := utility.CalculateFee(order.BrokerageInfo.Precision, order.OrderInfo.OrderAmount, order.BrokerageInfo.OrderThreshold, order.BrokerageInfo.MakerPct, order.BrokerageInfo.MakerMinFee)
	collateralAmount := utility.CalculateFee(order.BrokerageInfo.Precision, order.OrderInfo.OrderAmount, order.BrokerageInfo.OrderThreshold, order.BrokerageInfo.CollateralPct, order.BrokerageInfo.MinCollateral)

	makerCommittingAmount := order.OrderInfo.OrderAmount + makerFee + collateralAmount

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

	lastSlot := config.BFC.LastBlockSlot()

	collateralUtxo := config.BFC.GetUtxoFromRef(order.OrderTxInfo.CollateralUtxo.TxID, order.OrderTxInfo.CollateralUtxo.TxIDIndex)

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
			constants.INDEX_ONE_MINT_REDEEMER,
		).
		AddReferenceInput(
			order.OrderTxInfo.StateTokenRefUtxo.TxID,
			order.OrderTxInfo.StateTokenRefUtxo.TxIDIndex,
		).
		PayToContract(
			order.OrderTxInfo.EscrowContractAddress,
			orderDatumMarshaled,
			int(order.BrokerageInfo.AdaCollateral*2),
			true,
			apollo.Unit{
				PolicyId: order.OrderTxInfo.StateTokenPolicyId,
				Name:     order.OrderInfo.OrderId,
				Quantity: int(1),
			},
			apollo.Unit{
				PolicyId: order.OrderInfo.TradeTokenPolicyId,
				Name:     order.OrderInfo.TradeTokenName,
				Quantity: int(makerCommittingAmount),
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

	// fmt.Println("TX EVAL: ", config.BFC.EvaluateTx(txByte))
	// fmt.Println("CBOR: ", Utils.ToCbor(tx))
	
	if len(config.BFC.EvaluateTx(txByte)) == 0 {
		return "", "", fmt.Errorf("transaction evaluation failed")
	}

	return Utils.ToCbor(tx), hex.EncodeToString(txHash), nil

}
