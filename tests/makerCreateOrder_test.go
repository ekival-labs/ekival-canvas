package tests

import (
	"encoding/hex"
	"os"
	"testing"
	"time"

	"ekival-canvas/config"
	"ekival-canvas/constants"
	"ekival-canvas/model"
	"ekival-canvas/txBuilders/ada_p2p_buy"
	"ekival-canvas/utility"

	"github.com/Salvionied/apollo"
	"github.com/Salvionied/apollo/serialization/Address"
	"github.com/Salvionied/apollo/txBuilding/Backend/BlockFrostChainContext"
)

func TestMakerCreateOrder(t *testing.T) {

	// Create the ./tmp directory if it doesn't exist
	err := os.MkdirAll("./tmp", os.ModePerm)
	if err != nil {
		t.Logf("ERROR: Failed to create ./tmp directory: %v", err)
		t.Error(err)
	}

	t.Log("=== Starting Maker Create Order Test ===")
	t.Logf("Test Constants: OrderID= %s, OrderAmount= %d lovelace, Precision= %d", ORDER_ID, ORDER_AMOUNT, PRECISION)
	t.Logf("Fee Configuration: MakerMinFee= %d, TakerMinFee= %d, CancelMinFee= %d", MAKER_MIN_FEE, TAKER_MIN_FEE, CANCEL_MIN_FEE)
	t.Logf("Collateral Config: MinCollateral= %d, AdaCollateral= %d", MIN_COLLATERAL, ADA_COLLATERAL)
	t.Logf("Thresholds: MinOrderAmount= %d, OrderThreshold= %d", MIN_ORDER_AMOUNT, ORDER_THRESHOLD)
	t.Logf("Penalties: CancelPenalty= %d, ExtendingPenalty= %d, DisputeMinFee= %d, DisputePenalty= %d", CANCEL_PENALTY, EXTENDING_PENALTY, DISPUTE_MIN_FEE, DISPUTE_PENALTY)
	t.Logf("Deadlines: MakerDeadline= %d seconds, TakerDeadline= %d seconds", MAKER_DEADLINE, TAKER_DEADLINE)
	t.Logf("Percentages: CollateralPct= %d, MakerPct= %d, TakerPct= %d, CancelPct= %d, DisputePct= %d", COLLATERAL_PCT, MAKER_PCT, TAKER_PCT, CANCEL_PCT, DISPUTE_PCT)

	be, err := BlockFrostChainContext.NewBlockfrostChainContext(
		BFC_API_URL,
		BFC_NETWORK_ID,
		BFC_API_KEY,
	)

	if err != nil {
		t.Logf("ERROR: Failed to create Blockfrost chain context: %v", err)
		t.Error(err)
	} else {
		t.Log("✓ Successfully created Blockfrost chain context")
	}

	// Initialize the global chain context used by MakerCreateOrder
	err = config.ChainCTXSetup()
	if err != nil {
		t.Logf("ERROR: Failed to set up global chain context: %v", err)
		t.Error(err)
	} else {
		t.Log("✓ Successfully set up global chain context")
	}

	ma, err := Address.DecodeAddress(MAKER_ADDRESS)
	if err != nil {
		t.Logf("ERROR: Failed to decode maker address %s: %v", MAKER_ADDRESS, err)
		t.Error(err)
	} else {
		t.Logf("✓ Decoded maker address: %s", ma.String())
		t.Logf("  Maker PaymentPart: %s", hex.EncodeToString(ma.PaymentPart))
	}

	ta, err := Address.DecodeAddress(TAKER_ADDRESS)
	if err != nil {
		t.Logf("ERROR: Failed to decode taker address %s: %v", TAKER_ADDRESS, err)
		t.Error(err)
	} else {
		t.Logf("✓ Decoded taker address: %s", ta.String())
		t.Logf("  Taker PaymentPart: %s", hex.EncodeToString(ta.PaymentPart))
	}

	escrowContractAddress, err := Address.DecodeAddress(ADA_P2P_BUY_ESCROW_ADDRESS)
	if err != nil {
		t.Logf("ERROR: Failed to decode escrow contract address %s: %v", ADA_P2P_BUY_ESCROW_ADDRESS, err)
		t.Error(err)
	} else {
		t.Logf("✓ Decoded escrow contract address: %s", escrowContractAddress.String())
	}

	currentTime := time.Now().UnixNano() / int64(time.Millisecond)
	var md int64 = currentTime + (MAKER_DEADLINE * 1000)
	var td int64 = currentTime + (TAKER_DEADLINE * 1000)

	t.Logf("Time Calculations: CurrentTime= %d ms, MakerDeadline= %d ms, TakerDeadline= %d ms", currentTime, md, td)
	t.Logf("Time Details: Current Unix Time= %d, Maker Deadline Offset= %d sec, Taker Deadline Offset= %d sec", currentTime, MAKER_DEADLINE, TAKER_DEADLINE)

	makerFee, err := utility.CalculateFee(PRECISION, ORDER_THRESHOLD, MIN_ORDER_AMOUNT, MAKER_MIN_FEE, MAKER_PCT, ORDER_AMOUNT)
	if err != nil {
		t.Logf("ERROR: Failed to calculate maker fee: %v", err)
		t.Error(err)
	}
	collateralAmount, err := utility.CalculateFee(PRECISION, ORDER_THRESHOLD, MIN_ORDER_AMOUNT, MIN_COLLATERAL, COLLATERAL_PCT, ORDER_AMOUNT)
	if err != nil {
		t.Logf("ERROR: Failed to calculate collateral amount: %v", err)
		t.Error(err)
	}
	makerCommittingAmount := ORDER_AMOUNT + makerFee + collateralAmount

	t.Logf("Fee Calculations: MakerFee= %d lovelace, CollateralAmount= %d lovelace", makerFee, collateralAmount)
	t.Logf("Fee Calculation Details: Precision= %d, OrderAmount= %d, OrderThreshold= %d, MakerPct= %d, CollateralPct= %d", PRECISION, ORDER_AMOUNT, ORDER_THRESHOLD, TAKER_PCT, COLLATERAL_PCT)

	// Get taker UTxOs for spending
	userUtxosForMaker, err := be.Utxos(ma)
	if err != nil {
		t.Logf("ERROR: Failed to fetch user UTxOs for taker address %s: %v", ta.String(), err)
		t.Error(err)
	} else {
		t.Logf("✓ Fetched %d user UTxOs for taker address", len(userUtxosForMaker))
		// for i := range userUtxosForMaker {
		// 	t.Logf("  UTxO %d: Available for transaction input", i+1)
		// }
	}

	if len(userUtxosForMaker) == 0 {
		t.Logf("ERROR: No UTxOs available for taker")
		t.Error("no UTxOs available for taker")
	}

	// Convert some taker UTxOs to model.EUTxO format for spending
	var makerUserUtxos []model.EUTxO
	var currentTakerAmount int64
	for _, utxo := range userUtxosForMaker {
		makerUserUtxos = append(makerUserUtxos, model.EUTxO{
			TxID:      hex.EncodeToString(utxo.Input.TransactionId),
			TxIDIndex: utxo.Input.Index,
		})
		currentTakerAmount += utxo.Output.GetAmount().GetCoin()
		if currentTakerAmount >= makerCommittingAmount {
			break // Enough UTxOs selected to cover the amount
		}
	}

	if currentTakerAmount < makerCommittingAmount || len(makerUserUtxos) == 0 {
		t.Logf("ERROR: Not enough UTxOs or insufficient amount available for taker committing amount. Needed: %d, Got: %d", makerCommittingAmount, currentTakerAmount)
		t.Error("not enough UTxOs or insufficient amount for taker committing")
	}
	t.Logf("✓ Selected %d taker UTxOs for spending", len(makerUserUtxos))

	t.Log("=== Retrieving Collateral UTxO ===")

	// Use a different UTxO for collateral (not the same as spending UTxOs)
	var collateralUtxo model.EUTxO
	foundCollateral := false
	for _, utxo := range userUtxosForMaker {
		if utxo.Output.GetAmount().GetCoin() > 5000000 {
			// Check if this UTxO is already in makerUserUtxos
			isTakerUserUtxo := false
			for _, takerUtxo := range makerUserUtxos {
				if takerUtxo.TxID == hex.EncodeToString(utxo.Input.TransactionId) && takerUtxo.TxIDIndex == utxo.Input.Index {
					isTakerUserUtxo = true
					break
				}
			}
			if !isTakerUserUtxo {
				collateralUtxo = model.EUTxO{
					TxID:      hex.EncodeToString(utxo.Input.TransactionId),
					TxIDIndex: utxo.Input.Index,
				}
				foundCollateral = true
				break
			}
		}
	}

	if !foundCollateral {
		t.Logf("ERROR: No suitable collateral UTxO found (lovelace amount > 5 and not in makerUserUtxos)")
		t.Error("no suitable collateral UTxO found")
	}
	t.Logf("✓ Selected collateral UTxO: %s#%d", collateralUtxo.TxID, collateralUtxo.TxIDIndex)

	t.Log("=== Constructing Order Object ===")

	order := &model.Order{
		OrderInfo: model.OrderInfo{
			OrderId:         ORDER_ID,
			OrderAmount:     ORDER_AMOUNT,
			MakerAddress:    ma,
			MakerRepAddress: model.Nothing{},
			TakerAddress:    ta,
			TakerRepAddress: model.Nothing{},
			MakerDeadline:   md,
			TakerDeadline:   td,
		},
		BrokerageInfo: model.BrokerageInfo{
			Precision:      PRECISION,
			CollateralPct:  COLLATERAL_PCT,
			TakerPct:       TAKER_PCT,
			MakerPct:       MAKER_PCT,
			CancelPct:      CANCEL_PCT,
			MinCollateral:  MIN_COLLATERAL,
			MakerMinFee:    MAKER_MIN_FEE,
			TakerMinFee:    TAKER_MIN_FEE,
			CancelMinFee:   CANCEL_MIN_FEE,
			MinOrderAmount: MIN_ORDER_AMOUNT,
			OrderThreshold: ORDER_THRESHOLD,
			CancelPenalty:  CANCEL_PENALTY,
		},
		TradeState: constants.UNCOMMITTED_ORDER_STATUS,

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
			UserUtxos:        makerUserUtxos,
			CollateralUtxo:   collateralUtxo,
		},
	}

	t.Logf("✓ Order constructed successfully")
	t.Logf("  OrderInfo: ID= %s, Amount= %d, Maker= %s, Taker= %s", order.OrderInfo.OrderId, order.OrderInfo.OrderAmount, order.OrderInfo.MakerAddress.String(), order.OrderInfo.TakerAddress.String())
	t.Logf("  Deadlines: Maker= %d, Taker= %d", order.OrderInfo.MakerDeadline, order.OrderInfo.TakerDeadline)
	t.Logf("  BrokerageInfo: Precision= %d, CollateralPct= %d, MakerPct= %d, CancelPct= %d", order.BrokerageInfo.Precision, order.BrokerageInfo.CollateralPct, order.BrokerageInfo.MakerPct, order.BrokerageInfo.CancelPct)
	t.Logf("  TradeState: %s", order.TradeState)
	t.Logf("  OrderTxInfo: EscrowAddr= %s, StateTokenPolicy= %s", order.OrderTxInfo.EscrowContractAddress.String(), order.OrderTxInfo.StateTokenPolicyId)
	t.Logf("  Reference UTxOs: EscrowRef= %s#%d, StateTokenRef= %s#%d", order.OrderTxInfo.EscrowContractRefUtxo.TxID, order.OrderTxInfo.EscrowContractRefUtxo.TxIDIndex, order.OrderTxInfo.StateTokenRefUtxo.TxID, order.OrderTxInfo.StateTokenRefUtxo.TxIDIndex)
	t.Logf("  Fees & Amounts: MakerFee= %d, CollateralAmount= %d", order.OrderTxInfo.MakerFee, order.OrderTxInfo.CollateralAmount)

	t.Run("Successful Order Creation via MakerCreateOrder Function", func(t *testing.T) {

		t.Log("=== Calling MakerCreateOrder Function ===")

		cbor, txHash, err := ada_p2p_buy.MakerCreateOrder(order, adminWallet)
		if err != nil {
			t.Logf("ERROR: MakerCreateOrder function failed: %v", err)
			t.Error(err)
		} else {
			t.Log("✓ MakerCreateOrder function executed successfully")
			t.Logf("  Transaction Hash: %s", txHash)
			t.Logf("  CBOR Length: %d characters", len(cbor))
		}

		t.Log("=== Validating Transaction Hash ===")

		if txHash == "" {
			t.Logf("ERROR: Transaction hash is empty")
			t.Errorf("transaction hash is empty")
		} else {
			t.Logf("✓ Transaction hash is valid: %s", txHash)

			// Verify hash is valid hex
			if _, err := hex.DecodeString(txHash); err != nil {
				t.Logf("ERROR: Transaction hash is not valid hex: %v", err)
				t.Error(err)
			} else {
				t.Log("✓ Transaction hash is valid hex")
			}
		}

		t.Log("=== Validating CBOR ===")

		if cbor == "" {
			t.Logf("ERROR: CBOR is empty")
			t.Errorf("cbor is empty")
		} else {
			t.Logf("✓ CBOR is not empty")
			t.Logf("  CBOR (first 100 chars): %s...", cbor[:min(100, len(cbor))])
		}

		t.Log("=== Re-evaluating Transaction for Validation ===")

		// Convert CBOR back to transaction bytes for evaluation
		txBytes, err := hex.DecodeString(cbor)
		if err != nil {
			t.Logf("ERROR: Failed to decode CBOR to bytes: %v", err)
			t.Error(err)
		} else {
			t.Logf("✓ CBOR decoded to bytes: %d bytes", len(txBytes))
		}

		// Evaluate the transaction
		evalTx, err := be.EvaluateTx(txBytes)
		if err != nil {
			t.Logf("ERROR: Transaction evaluation failed: %v", err)
			t.Error(err)
		} else {
			t.Logf("✓ Transaction evaluated successfully")
			t.Logf("  Evaluation Results: %d execution units", len(evalTx))
			t.Logf("  Full evaluation details: %+v", evalTx)
		}

		if len(evalTx) == 0 {
			t.Logf("ERROR: Transaction evaluation returned no execution units - validation failed")
			t.Errorf("transaction evaluation failed")
		} else {
			t.Log("✓ Transaction validation passed")
		}

		t.Log("=== Submitting Transaction ===")
		// Initialize apolloBE for submission
		apolloBE := apollo.New(&config.CHAIN_CTX)
		apolloBE, err = apolloBE.LoadTxCbor(cbor)
		if err != nil {
			t.Logf("ERROR: Failed to load Tx CBOR for submission: %v", err)
			t.Error(err)
		}

		apolloBE, err = apolloBE.SignWithSkey(makerWallet.Vkey, makerWallet.Skey)
		if err != nil {
			t.Logf("ERROR: Failed to sign with maker key: %v", err)
			t.Error(err)
		} else {
			t.Log("✓ Signed transaction with maker key")
		}

		txHashSubmit, err := apolloBE.Submit()
		if err != nil {
			t.Logf("ERROR: Failed to submit transaction: %v", err)
			t.Error(err)
		} else {
			t.Logf("✓ Transaction submitted successfully")
			t.Logf("  Transaction ID: %s", hex.EncodeToString(txHashSubmit.Payload))
			t.Logf("  TxID Details: Hash length= %d bytes, Full payload=%x", len(txHashSubmit.Payload), txHashSubmit.Payload)
			t.Logf("  Deadlines: Maker= %d, Taker= %d", order.OrderInfo.MakerDeadline, order.OrderInfo.TakerDeadline)
		}

	})

	t.Log("=== Test Completed Successfully ===")
}

// Helper function for min
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
