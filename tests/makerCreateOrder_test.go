package tests

import (
	"encoding/hex"
	"testing"
	"time"

	"ekival-canvas/constants"
	"ekival-canvas/model"
	"ekival-canvas/plutusEncoder"
	"ekival-canvas/utility"

	"github.com/Salvionied/apollo"
	"github.com/Salvionied/apollo/serialization"
	"github.com/Salvionied/apollo/serialization/Address"
	"github.com/Salvionied/apollo/txBuilding/Backend/MaestroChainContext"
	"github.com/Salvionied/apollo/txBuilding/Utils"
)

func TestMain(m *testing.M) {
	m.Run()
}

func TestMakerCreateOrder(t *testing.T) {
	// be := OgmiosChainContext.NewOgmiosChainContext(*ogmigo.New(ogmigo.WithEndpoint(BLINKLABS_OGMIOS_ENDPOINT)), *kugo.New(kugo.WithEndpoint(BLINKLABS_KUPO_ENDPOINT)))

	be, err := MaestroChainContext.NewMaestroChainContext(
		MAESTRO_NETWORK_ID,
		MAESTRO_API_KEY,
	)

	// be, err := BlockFrostChainContext.NewBlockfrostChainContext(
	// 	BFC_API_URL,
	// 	BFC_NETWORK_ID,
	// 	BFC_API_KEY,
	// )

	// if err != nil {
	// 	t.Error(err)
	// }

	ma, err := Address.DecodeAddress(MAKER_ADDRESS)
	if err != nil {
		t.Error(err)
	}

	ta, err := Address.DecodeAddress(TAKER_ADDRESS)
	if err != nil {
		t.Error(err)
	}

	escrowContractAddress, err := Address.DecodeAddress(ADA_P2P_BUY_ESCROW_ADDRESS)
	if err != nil {
		t.Error(err)
	}

	currentTime := time.Now().UnixNano() / int64(time.Millisecond)
	var md int64 = currentTime + (makerDeadline * 1000)
	var td int64 = currentTime + (takerDeadline * 1000)

	makerFee := utility.CalculateFee(precision, orderAmount, orderThreshold, makerPct, makerMinFee)
	collateralAmount := utility.CalculateFee(precision, orderAmount, orderThreshold, collateralPct, minCollateral)

	mockOrder := &model.Order{
		OrderInfo: model.OrderInfo{
			OrderId:         orderId,
			OrderAmount:     orderAmount,
			MakerAddress:    ma,
			MakerRepAddress: model.Nothing{},
			TakerAddress:    ta,
			TakerRepAddress: model.Nothing{},
			MakerDeadline:   md,
			TakerDeadline:   td,
		},
		BrokerageInfo: model.BrokerageInfo{
			Precision:      precision,
			CollateralPct:  collateralPct,
			MakerPct:       takerPct,
			CancelPct:      cancelPct,
			MinCollateral:  minCollateral,
			MakerMinFee:    makerMinFee,
			TakerMinFee:    takerMinFee,
			CancelMinFee:   cancelMinFee,
			MinOrderAmount: minOrderAmount,
			OrderThreshold: orderThreshold,
			CancelPenalty:  cancelPenalty,
		},
		TradeState: UNCOMMITTED_ORDER_STATUS,

		OrderTxInfo: model.OrderTxInfo{
			EscrowContractAddress: escrowContractAddress,
			EscrowContractRefUtxo: model.EUTxO{
				TxID:      ADA_P2P_BUY_ESCROW_SCRIPT_REF_UTXO_TXID,
				TxIDIndex: ADA_P2P_BUY_ESCROW_SCRIPT_REF_UTXO_TXID_INDEX,
			},
			StateTokenPolicyId: APBST_POLICY_ID,
			StateTokenRefUtxo: model.EUTxO{
				TxID:      APBST_SCRIPT_REF_UTXO_TXID,
				TxIDIndex: APBST_SCRIPT_REF_UTXO_TXID_INDEX,
			},
			MakerFee:         makerFee,
			CollateralAmount: collateralAmount,
			ChangeAddress:    ma,
			UserUtxos:        []model.EUTxO{},
			CollateralUtxo:   model.EUTxO{},
		},
	}

	t.Run("Successful Order Creation", func(t *testing.T) {

		apolloBE := apollo.New(&be)
		apolloBE = apolloBE.SetWalletFromBech32(mockOrder.OrderInfo.MakerAddress.String())

		makerCommittingAmount := mockOrder.OrderInfo.OrderAmount + mockOrder.OrderTxInfo.MakerFee + mockOrder.OrderTxInfo.CollateralAmount

		orderDatumMarshaled, err := plutusEncoder.MarshalPlutus(*mockOrder)
		if err != nil {
			t.Error(err)
		}

		userUtxos, err := be.Utxos(ma)
		if err != nil {
			t.Error(err)
		}

		lastSlot, err := be.LastBlockSlot()
		if err != nil {
			t.Error(err)
		}

		apolloBE, err = apolloBE.
			SetChangeAddress(mockOrder.OrderTxInfo.ChangeAddress).
			AddLoadedUTxOs(userUtxos...).
			MintAssetsWithRedeemer(
				apollo.Unit{
					PolicyId: mockOrder.OrderTxInfo.StateTokenPolicyId,
					Name:     mockOrder.OrderInfo.OrderId,
					Quantity: int(1),
				},
				*constants.INDEX_ONE_MINT_REDEEMER,
			).
			AddReferenceInputV3(
				mockOrder.OrderTxInfo.StateTokenRefUtxo.TxID,
				mockOrder.OrderTxInfo.StateTokenRefUtxo.TxIDIndex,
			).
			PayToContract(
				mockOrder.OrderTxInfo.EscrowContractAddress,
				orderDatumMarshaled,
				int(makerCommittingAmount),
				true,
				apollo.Unit{
					PolicyId: mockOrder.OrderTxInfo.StateTokenPolicyId,
					Name:     mockOrder.OrderInfo.OrderId,
					Quantity: int(1),
				},
			).
			AddRequiredSigner(adminWallet.AdminPKH).
			AddRequiredSigner(serialization.PubKeyHash(mockOrder.OrderInfo.MakerAddress.PaymentPart)).
			SetTtl(int64(lastSlot) + 300).
			Complete()

		if err != nil {
			t.Error(err)
		}

		apolloBE, err = apolloBE.SignWithSkey(adminWallet.AdminVkey, adminWallet.AdminSkey)
		if err != nil {
			t.Error(err)
		}

		apolloBE, err = apolloBE.SignWithSkey(makerWallet.AdminVkey, makerWallet.AdminSkey)
		if err != nil {
			t.Error(err)
		}

		tx := apolloBE.GetTx()

		txByte, err := tx.Bytes()
		if err != nil {
			t.Error(err)
		}

		evalTx, err := be.EvaluateTx(txByte)
		if err != nil {
			t.Error(err)
		}

		cbor, err := Utils.ToCbor(tx)
		if err != nil {
			t.Error(err)
		}

		t.Log("Tx CBOR:", cbor)
		t.Log("EvaluateTx:", evalTx)

		if len(evalTx) == 0 {
			t.Errorf("transaction evaluation failed")
		}

		txHash, err := apolloBE.Submit()
		if err != nil {
			t.Error(err)
		}

		t.Log("TxID:", hex.EncodeToString(txHash.Payload))

	})

}
