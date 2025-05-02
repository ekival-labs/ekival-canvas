package tests

import (
	"encoding/hex"
	"testing"
	"time"

	"github.com/Salvionied/apollo"
	"github.com/Salvionied/apollo/serialization"
	"github.com/Salvionied/apollo/serialization/Address"
	"github.com/Salvionied/apollo/txBuilding/Backend/MaestroChainContext"
	"github.com/Salvionied/apollo/txBuilding/Utils"
	"github.com/ekival-labs/ekival-canvas/config"
	"github.com/ekival-labs/ekival-canvas/constants"
	"github.com/ekival-labs/ekival-canvas/model"
	"github.com/ekival-labs/ekival-canvas/plutusEncoder"
	"github.com/ekival-labs/ekival-canvas/utility"
)

func TestMain(m *testing.M) {
	m.Run()
}

const (
	MAKER_ADDRESS    string = "addr_test1qr43kyyys0sg0d8khhjq4y348zuc4mnzln836hxfrqjasx4rqghlrj99l5vrdmyrtg6mhkyxa88kwq5yf225a4m9pkes2a4xz3"
	TAKER_ADDRESS    string = "addr_test1qqt85kcauy3uktlfldmhqn4dnn5vxgv8s97t729akje92f3z5uathwrlk6dwzy9j89lmsy2qp7evugtxmhf9pycw44rqj76x7l"
	DEBUGGER_ADDRESS string = "addr_test1qqd9lnkkk4uj96ftyqcgune8jq7g04wccmzcd3hehrddl8at3u4s70lrcaw39t522nq3y9953k2w47t4mgtj678xpp7q0dx95k"
	ADMIN_ADDRESS    string = "addr_test1qp0ccf96g40e45wtc98g4nvdxrgxwdwy0c4ymlpq9gprts5eqx40n59saqk2jlwc49uw30mfmlm7nckja2sg6wlrclmqaymgpe"

	UNCOMMITTED_ORDER_STATUS string = "UNCOMMITTED_ORDER"
	COMMITTED_ORDER_STATUS   string = "COMMITTED_ORDER"
	REMIT_CONFIRMED_STATUS   string = "REMIT_CONFIRMED"

	TEKI_POLICY_ID  string = "fe691a72d5a591a2ec8602a9904b30c005138d20af67b36087b27221"
	TEKI_TOKEN_NAME string = "tEKI"

	FEE_TREASURY_ADDRESS string = "addr_test1xpcpdkxprshunrpvwfrq24swu8tpkl9f33v7rq45welsppm7axsw0am6sjwdpkpejhfpzmwe2k50cckcselxt5hlmnqqz444tl"

	ADA_P2P_BUY_ESCROW_ADDRESS                    string = "addr_test1xrgysjt2g0t4l7h54px984wehx9uhtautarly2nqq83w5622m3nsefutuywd8uqzd2n4u7w0vmymjhmmuuc46hdxy85qzlh2dc"
	ADA_P2P_BUY_ESCROW_SCRIPT_REF_UTXO_TXID       string = "87ce9218b91f09156b15bdb6fa5c3a8aa78d846184f2e1b5375df0114abd8b6d"
	ADA_P2P_BUY_ESCROW_SCRIPT_REF_UTXO_TXID_INDEX int    = 0
	APBST_POLICY_ID                               string = "aefdb5f954ea897ec536de425e62632728fa78b4638882654d0f4075"
	APBST_SCRIPT_REF_UTXO_TXID                    string = "fbf2cc43233490bc7f748957711bda85a530484113c8b65381b0613c30c13ac2"
	APBST_SCRIPT_REF_UTXO_TXID_INDEX              int    = 0

	ADA_P2P_BUY_FEE_TYPE int = 1

	NETWORK            string = "preprod"
	BFC_NETWORK_ID     int    = 0
	MAESTRO_NETWORK_ID int    = 3 // preprod
	BFC_API_URL        string = "https://cardano-preprod.blockfrost.io/api"
	BFC_API_KEY        string = "preprod9zzl4g8Xa3faU50a1OVDZdPeQ92ZsdcT"
	MAESTRO_API_KEY    string = "so4a45BCnj80EdcFa9OwLr8pK8um4bWE"
	OGMIGO_ENDPOINT    string = "ws://localhost:1337"
	KUGO_ENDPOINT      string = "http://localhost:1442"

	EKIVAL_FEE int = 25_000_000
)

var (
	orderUTxOTxId      string = "b8d70feea6bca4bc51e4478490753a5020d3e2aecc1d83b1b307543a47b6a6a6"
	orderUTxOTxIdIndex int    = 0

	orderId    string = "AwesomeID_TT_RR_1"
	toLovelace int64  = 1000000
	// orderAmount int64  = 719854678691
	orderAmount      int64 = 197 * toLovelace
	precision        int64 = 10
	collateralPct    int64 = utility.ToFraction(10)
	makerPct         int64 = utility.ToFraction(0.25)
	takerPct         int64 = utility.ToFraction(0.75)
	cancelPct        int64 = utility.ToFraction(1)
	minCollateral    int64 = 25 * toLovelace
	makerMinFee      int64 = 1250000
	takerMinFee      int64 = 3750000
	cancelMinFee     int64 = 3 * toLovelace
	minOrderAmount   int64 = 10 * toLovelace
	orderThreshold   int64 = 500 * toLovelace
	cancelPenalty    int64 = 2 * cancelMinFee
	adaCollateral    int64 = 3 * toLovelace
	extendingPenalty int64 = 5 * toLovelace
	disputePct       int64 = utility.ToFraction(2)
	disputeMinFee    int64 = 5 * toLovelace
	disputePenalty   int64 = 2 * disputeMinFee

	makerDeadline int64 = 7200
	takerDeadline int64 = 3600

	adminWallet = config.SetWallet("garment shadow into lab truck quiz file warm sheriff marriage voice icon thunder iron early jungle bird dash material strike desk mango deer letter")
	makerWallet = config.SetWallet("bone miracle mother grocery rabbit decorate rain moment print empty harbor cinnamon resource desert roof attitude suspect cupboard allow hunt inhale praise sausage mom")
)

func TestMakerCreateOrder(t *testing.T) {
	// be := OgmiosChainContext.NewOgmiosChainContext(*ogmigo.New(ogmigo.WithEndpoint(OGMIGO_ENDPOINT)), *kugo.New(kugo.WithEndpoint(KUGO_ENDPOINT)))

	be, err := MaestroChainContext.NewMaestroChainContext(
		MAESTRO_NETWORK_ID,
		MAESTRO_API_KEY,
	)

	// be, err := BlockFrostChainContext.NewBlockfrostChainContext(
	// 	BFC_API_URL,
	// 	BFC_NETWORK_ID,
	// 	BFC_API_KEY,
	// )

	if err != nil {
		t.Error(err)
	}

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
			AddReferenceInput(
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
