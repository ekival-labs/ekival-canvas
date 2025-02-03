package token_p2p_buy

import (
	"ekival-canvas/config"
	"ekival-canvas/constants"
	"ekival-canvas/model"
	"ekival-canvas/plutusEncoder"
	"ekival-canvas/utility"
	"encoding/hex"
	"fmt"

	"github.com/Salvionied/apollo"
	"github.com/Salvionied/apollo/serialization"
	"github.com/Salvionied/apollo/txBuilding/Utils"
	fiberLogger "github.com/gofiber/fiber/v2/log"
)

func MakerCreateOrder(order *model.Order, adminWallet *config.AdminWallet) (string, string, error) {

	defer func() {
		if err := recover(); err != nil {
			fiberLogger.Panic("Panic occurred: %v", err)
			return
		}
	}()

	apolloBE := apollo.New(&config.CHAIN_CTX)
	apolloBE = apolloBE.SetWalletFromBech32(order.OrderInfo.MakerAddress.String())

	makerCommittingAmount := order.OrderInfo.OrderAmount + order.OrderTxInfo.MakerFee + order.OrderTxInfo.CollateralAmount

	orderDatumMarshaled, err := plutusEncoder.MarshalPlutus(*order)
	if err != nil {
		fiberLogger.Error(err)
		return "", "", err
	}

	userUtxos, err := utility.GetUserUTxOs(order.OrderTxInfo.UserUtxos)
	if err != nil {
		fiberLogger.Error(err)
		return "", "", err
	}

	lastSlot, err := config.CHAIN_CTX.LastBlockSlot()
	if err != nil {
		fiberLogger.Error(err)
		return "", "", err
	}

	collateralUtxo, err := config.CHAIN_CTX.GetUtxoFromRef(order.OrderTxInfo.CollateralUtxo.TxID, order.OrderTxInfo.CollateralUtxo.TxIDIndex)
	if err != nil {
		fiberLogger.Error(err)
		return "", "", err
	}

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
		fiberLogger.Error(err)
		return "", "", err
	}

	apolloBE, err = apolloBE.SignWithSkey(adminWallet.AdminVkey, adminWallet.AdminSkey)
	if err != nil {
		fiberLogger.Error(err)
		return "", "", err
	}

	tx := apolloBE.GetTx()
	txHash, err := tx.TransactionBody.Hash()
	if err != nil {
		fiberLogger.Error(err)
		return "", "", err
	}

	txByte, err := tx.Bytes()
	if err != nil {
		fiberLogger.Error(err)
		return "", "", err
	}

	evalTx, err := config.CHAIN_CTX.EvaluateTx(txByte)
	if err != nil {
		fiberLogger.Error(err)
		return "", "", err
	}

	cbor, err := Utils.ToCbor(tx)
	if err != nil {
		fiberLogger.Error(err)
		return "", "", err
	}

	fiberLogger.Debug("TxID:", hex.EncodeToString(txHash))
	fiberLogger.Debug("Tx CBOR:", cbor)
	fiberLogger.Debug("EvaluateTx:", evalTx)

	if len(evalTx) == 0 {
		return "", "", fmt.Errorf("transaction evaluation failed")
	}

	return cbor, hex.EncodeToString(txHash), nil

}
