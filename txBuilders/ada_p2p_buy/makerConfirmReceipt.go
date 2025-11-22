package ada_p2p_buy

import (
	"ekival-canvas/config"
	"ekival-canvas/constants"
	"ekival-canvas/model"
	"ekival-canvas/utility"
	"encoding/hex"
	"fmt"
	"runtime/debug"

	"github.com/Salvionied/apollo"
	"github.com/Salvionied/apollo/serialization"
	"github.com/Salvionied/apollo/txBuilding/Utils"
	"github.com/rs/zerolog/log"
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
func MakerConfirmReceipt(order *model.Order, treasury *model.TreasuryInfo, adminWallet *config.Wallet) (string, string, error) {
	log.Trace().
		Str("function", "MakerConfirmReceipt").
		Str("order_id", order.OrderInfo.OrderId).
		Msg("Function entry: Confirming receipt for maker order")

	defer func() {
		if r := recover(); r != nil {
			stack := debug.Stack()
			log.Panic().
				Interface("panic_value", r).
				Bytes("stack_trace", stack).
				Str("function", "MakerConfirmReceipt").
				Str("order_id", order.OrderInfo.OrderId).
				Msg("Panic occurred in MakerConfirmReceipt")
			panic(r)
		}
	}()

	// Log input parameters
	log.Debug().
		Str("function", "MakerConfirmReceipt").
		Msg("Validating input parameters")
	if order == nil {
		log.Error().
			Str("function", "MakerConfirmReceipt").
			Msg("ERROR: order parameter is nil")
		return "", "", fmt.Errorf("order parameter is nil")
	}
	if treasury == nil {
		log.Error().
			Str("function", "MakerConfirmReceipt").
			Str("order_id", order.OrderInfo.OrderId).
			Msg("ERROR: treasury parameter is nil")
		return "", "", fmt.Errorf("treasury parameter is nil")
	}
	if adminWallet == nil {
		log.Error().
			Str("function", "MakerConfirmReceipt").
			Str("order_id", order.OrderInfo.OrderId).
			Msg("ERROR: adminWallet parameter is nil")
		return "", "", fmt.Errorf("adminWallet parameter is nil")
	}

	log.Debug().
		Str("function", "MakerConfirmReceipt").
		Str("order_id", order.OrderInfo.OrderId).
		Int64("order_amount", order.OrderInfo.OrderAmount).
		Str("maker_address", order.OrderInfo.MakerAddress.String()).
		Str("taker_address", order.OrderInfo.TakerAddress.String()).
		Int64("maker_fee", order.OrderTxInfo.MakerFee).
		Int64("taker_fee", order.OrderTxInfo.TakerFee).
		Int64("collateral_amount", order.OrderTxInfo.CollateralAmount).
		Str("change_address", order.OrderTxInfo.ChangeAddress.String()).
		Str("state_token_policy_id", order.OrderTxInfo.StateTokenPolicyId).
		Str("treasury_address", treasury.Address.String()).
		Interface("admin_pkh", adminWallet.PKH).
		Msg("Input parameters validated and logged")

	// Initialize Apollo backend
	log.Debug().
		Str("function", "MakerConfirmReceipt").
		Str("order_id", order.OrderInfo.OrderId).
		Msg("Initializing Apollo backend")
	apolloBE := apollo.New(&config.CHAIN_CTX)
	log.Debug().
		Str("function", "MakerConfirmReceipt").
		Str("order_id", order.OrderInfo.OrderId).
		Msg("Apollo backend created successfully")

	log.Debug().
		Str("function", "MakerConfirmReceipt").
		Str("order_id", order.OrderInfo.OrderId).
		Str("maker_address", order.OrderInfo.MakerAddress.String()).
		Msg("Setting wallet from Bech32 address")
	apolloBE = apolloBE.SetWalletFromBech32(order.OrderInfo.MakerAddress.String())
	log.Debug().
		Str("function", "MakerConfirmReceipt").
		Str("order_id", order.OrderInfo.OrderId).
		Str("maker_address", order.OrderInfo.MakerAddress.String()).
		Msg("Wallet set from Bech32")

	// Calculate payment amounts
	log.Debug().
		Str("function", "MakerConfirmReceipt").
		Str("order_id", order.OrderInfo.OrderId).
		Msg("Calculating payment amounts")
	payToMakerAmount := order.OrderTxInfo.CollateralAmount
	payToTakerAmount := order.OrderInfo.OrderAmount + order.OrderTxInfo.CollateralAmount
	ekivalFee := order.OrderTxInfo.MakerFee + order.OrderTxInfo.TakerFee
	log.Debug().
		Str("function", "MakerConfirmReceipt").
		Str("order_id", order.OrderInfo.OrderId).
		Int64("pay_to_maker_amount", payToMakerAmount).
		Int64("pay_to_taker_amount", payToTakerAmount).
		Int64("ekival_fee", ekivalFee).
		Msg("Payment amounts calculated")

	// Get Collateral UTxO
	log.Debug().
		Str("function", "MakerConfirmReceipt").
		Str("order_id", order.OrderInfo.OrderId).
		Str("collateral_utxo_txid", order.OrderTxInfo.CollateralUtxo.TxID).
		Int("collateral_utxo_index", order.OrderTxInfo.CollateralUtxo.TxIDIndex).
		Msg("Fetching Collateral UTxO")
	collateralUtxo, err := config.CHAIN_CTX.GetUtxoFromRef(order.OrderTxInfo.CollateralUtxo.TxID, order.OrderTxInfo.CollateralUtxo.TxIDIndex)
	if err != nil {
		log.Error().
			Err(err).
			Str("function", "MakerConfirmReceipt").
			Str("order_id", order.OrderInfo.OrderId).
			Str("collateral_utxo_txid", order.OrderTxInfo.CollateralUtxo.TxID).
			Int("collateral_utxo_index", order.OrderTxInfo.CollateralUtxo.TxIDIndex).
			Msg("ERROR: Failed to get Collateral UTxO")
		return "", "", fmt.Errorf("failed to get Collateral UTxO: %w", err)
	}
	log.Debug().
		Str("function", "MakerConfirmReceipt").
		Str("order_id", order.OrderInfo.OrderId).
		Interface("collateral_utxo", collateralUtxo).
		Msg("Collateral UTxO fetched successfully")

	// Get Order UTxO
	log.Debug().
		Str("function", "MakerConfirmReceipt").
		Str("order_id", order.OrderInfo.OrderId).
		Str("order_utxo_txid", order.OrderTxInfo.OrderUtxo.TxID).
		Int("order_utxo_index", order.OrderTxInfo.OrderUtxo.TxIDIndex).
		Msg("Fetching Order UTxO")
	orderUTxO, err := config.CHAIN_CTX.GetUtxoFromRef(order.OrderTxInfo.OrderUtxo.TxID, order.OrderTxInfo.OrderUtxo.TxIDIndex)
	if err != nil {
		log.Error().
			Err(err).
			Str("function", "MakerConfirmReceipt").
			Str("order_id", order.OrderInfo.OrderId).
			Str("order_utxo_txid", order.OrderTxInfo.OrderUtxo.TxID).
			Int("order_utxo_index", order.OrderTxInfo.OrderUtxo.TxIDIndex).
			Msg("ERROR: Failed to get Order UTxO")
		return "", "", fmt.Errorf("failed to get Order UTxO: %w", err)
	}
	log.Debug().
		Str("function", "MakerConfirmReceipt").
		Str("order_id", order.OrderInfo.OrderId).
		Interface("order_utxo", orderUTxO).
		Msg("Order UTxO fetched successfully")

	// Get User UTxOs
	log.Debug().
		Str("function", "MakerConfirmReceipt").
		Str("order_id", order.OrderInfo.OrderId).
		Int("user_utxos_count_input", len(order.OrderTxInfo.UserUtxos)).
		Msg("Fetching User UTxOs")
	userUtxos, err := utility.GetUserUTxOs(order.OrderTxInfo.UserUtxos)
	if err != nil {
		log.Error().
			Err(err).
			Str("function", "MakerConfirmReceipt").
			Str("order_id", order.OrderInfo.OrderId).
			Interface("user_utxos_input", order.OrderTxInfo.UserUtxos).
			Msg("ERROR: Failed to get User UTxOs")
		return "", "", fmt.Errorf("failed to get User UTxOs: %w", err)
	}
	log.Debug().
		Str("function", "MakerConfirmReceipt").
		Str("order_id", order.OrderInfo.OrderId).
		Int("user_utxos_count_processed", len(userUtxos)).
		Msg("User UTxOs fetched successfully")
	for i, utxo := range userUtxos {
		log.Trace().
			Str("function", "MakerConfirmReceipt").
			Str("order_id", order.OrderInfo.OrderId).
			Int("index", i).
			Interface("utxo", utxo).
			Msg("Individual User UTxO")
	}

	// Get last slot
	log.Debug().
		Str("function", "MakerConfirmReceipt").
		Str("order_id", order.OrderInfo.OrderId).
		Msg("Fetching last block slot")
	lastSlot, err := config.CHAIN_CTX.LastBlockSlot()
	if err != nil {
		log.Error().
			Err(err).
			Str("function", "MakerConfirmReceipt").
			Str("order_id", order.OrderInfo.OrderId).
			Msg("ERROR: Failed to get last block slot")
		return "", "", fmt.Errorf("failed to get last block slot: %w", err)
	}
	log.Debug().
		Str("function", "MakerConfirmReceipt").
		Str("order_id", order.OrderInfo.OrderId).
		Int("last_slot", lastSlot).
		Msg("Last block slot fetched successfully")
	ttl := int64(lastSlot) + 300
	log.Debug().
		Str("function", "MakerConfirmReceipt").
		Str("order_id", order.OrderInfo.OrderId).
		Int64("ttl", ttl).
		Msg("Transaction TTL calculated")

	// Build transaction
	log.Info().
		Str("function", "MakerConfirmReceipt").
		Str("order_id", order.OrderInfo.OrderId).
		Msg("Building transaction with Apollo")

	log.Trace().
		Str("function", "MakerConfirmReceipt").
		Str("order_id", order.OrderInfo.OrderId).
		Str("change_address", order.OrderTxInfo.ChangeAddress.String()).
		Str("collateral_utxo", order.OrderTxInfo.CollateralUtxo.TxID).
		Int("collateral_utxo_index", order.OrderTxInfo.CollateralUtxo.TxIDIndex).
		Int("user_utxos_count", len(userUtxos)).
		Str("order_utxo_txid", order.OrderTxInfo.OrderUtxo.TxID).
		Int("order_utxo_index", order.OrderTxInfo.OrderUtxo.TxIDIndex).
		Str("state_token_policy_id", order.OrderTxInfo.StateTokenPolicyId).
		Str("order_id_mint", order.OrderInfo.OrderId).
		Str("state_token_ref_utxo_txid", order.OrderTxInfo.StateTokenRefUtxo.TxID).
		Int("state_token_ref_utxo_index", order.OrderTxInfo.StateTokenRefUtxo.TxIDIndex).
		Str("escrow_contract_ref_utxo_txid", order.OrderTxInfo.EscrowContractRefUtxo.TxID).
		Int("escrow_contract_ref_utxo_index", order.OrderTxInfo.EscrowContractRefUtxo.TxIDIndex).
		Str("treasury_address", treasury.Address.String()).
		Int64("ekival_fee", ekivalFee).
		Str("maker_address", order.OrderInfo.MakerAddress.String()).
		Int64("pay_to_maker_amount", payToMakerAmount).
		Str("taker_address", order.OrderInfo.TakerAddress.String()).
		Int64("pay_to_taker_amount", payToTakerAmount).
		Interface("admin_pkh", adminWallet.PKH).
		Interface("maker_pkh", serialization.PubKeyHash(order.OrderInfo.MakerAddress.PaymentPart)).
		Int64("ttl", ttl).
		Msg("Transaction building parameters")

	apolloBE, err = apolloBE.
		SetChangeAddress(order.OrderTxInfo.ChangeAddress).
		AddCollateral(*collateralUtxo).
		AddLoadedUTxOs(userUtxos...).
		CollectFrom(*orderUTxO, *constants.INDEX_SIX_SPEND_REDEEMER).
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
			treasury.Address, treasury.Datum, int(ekivalFee), true,
		).
		PayToAddress(order.OrderInfo.MakerAddress, int(payToMakerAmount)).
		PayToAddress(order.OrderInfo.TakerAddress, int(payToTakerAmount)).
		AddRequiredSigner(adminWallet.PKH).
		AddRequiredSigner(serialization.PubKeyHash(order.OrderInfo.MakerAddress.PaymentPart)).
		SetTtl(ttl).
		SetFeePadding(apolloBE.Fee + 100000).
		Complete()

	if err != nil {
		log.Error().
			Err(err).
			Str("function", "MakerConfirmReceipt").
			Str("operation", "build_transaction").
			Str("order_id", order.OrderInfo.OrderId).
			Msg("ERROR: Transaction building failed")
		return "", "", fmt.Errorf("transaction building failed: %w", err)
	}
	log.Info().
		Str("function", "MakerConfirmReceipt").
		Str("order_id", order.OrderInfo.OrderId).
		Msg("Transaction built successfully, signing with admin key")

	// Sign transaction
	log.Debug().
		Str("function", "MakerConfirmReceipt").
		Str("order_id", order.OrderInfo.OrderId).
		Interface("admin_vkey", adminWallet.Vkey).
		Msg("Signing transaction with admin key")
	apolloBE, err = apolloBE.SignWithSkey(adminWallet.Vkey, adminWallet.Skey)
	if err != nil {
		log.Error().
			Err(err).
			Str("function", "MakerConfirmReceipt").
			Str("operation", "sign_transaction").
			Str("order_id", order.OrderInfo.OrderId).
			Msg("ERROR: Transaction signing failed")
		return "", "", fmt.Errorf("transaction signing failed: %w", err)
	}
	log.Debug().
		Str("function", "MakerConfirmReceipt").
		Str("order_id", order.OrderInfo.OrderId).
		Msg("Transaction signed successfully")

	// Get transaction and hash
	log.Debug().
		Str("function", "MakerConfirmReceipt").
		Str("order_id", order.OrderInfo.OrderId).
		Msg("Getting transaction and computing hash")
	tx := apolloBE.GetTx()
	if tx == nil {
		log.Error().
			Str("function", "MakerConfirmReceipt").
			Str("order_id", order.OrderInfo.OrderId).
			Msg("ERROR: Transaction is nil after GetTx()")
		return "", "", fmt.Errorf("transaction is nil after GetTx()")
	}
	log.Trace().
		Str("function", "MakerConfirmReceipt").
		Str("order_id", order.OrderInfo.OrderId).
		Interface("transaction_body", tx.TransactionBody). // Added log
		Msg("Transaction retrieved successfully, logging transaction body")

	txHash, err := tx.TransactionBody.Hash()
	if err != nil {
		log.Error().
			Err(err).
			Str("function", "MakerConfirmReceipt").
			Str("operation", "compute_tx_hash").
			Str("order_id", order.OrderInfo.OrderId).
			Msg("ERROR: Failed to compute transaction hash")
		return "", "", fmt.Errorf("failed to compute transaction hash: %w", err)
	}
	txHashHex := hex.EncodeToString(txHash)
	log.Debug().
		Str("function", "MakerConfirmReceipt").
		Str("order_id", order.OrderInfo.OrderId).
		Str("tx_hash_hex", txHashHex).
		Msg("Transaction hash computed successfully")

	// Serialize transaction to bytes
	log.Debug().
		Str("function", "MakerConfirmReceipt").
		Str("order_id", order.OrderInfo.OrderId).
		Str("tx_hash_hex", txHashHex).
		Msg("Serializing transaction to bytes")
	txByte, err := tx.Bytes()
	if err != nil {
		log.Error().
			Err(err).
			Str("function", "MakerConfirmReceipt").
			Str("operation", "serialize_tx_to_bytes").
			Str("order_id", order.OrderInfo.OrderId).
			Str("tx_hash_hex", txHashHex).
			Msg("ERROR: Failed to serialize transaction to bytes")
		return "", "", fmt.Errorf("failed to serialize transaction to bytes: %w", err)
	}
	log.Trace().
		Str("function", "MakerConfirmReceipt").
		Str("order_id", order.OrderInfo.OrderId).
		Str("tx_hash_hex", txHashHex).
		Int("tx_byte_length", len(txByte)).
		Msg("Transaction serialized successfully")

	// Convert to CBOR
	log.Debug().
		Str("function", "MakerConfirmReceipt").
		Str("order_id", order.OrderInfo.OrderId).
		Str("tx_hash_hex", txHashHex).
		Msg("Converting transaction to CBOR")
	cbor, err := Utils.ToCbor(tx)
	if err != nil {
		log.Error().
			Err(err).
			Str("function", "MakerConfirmReceipt").
			Str("operation", "convert_to_cbor").
			Str("order_id", order.OrderInfo.OrderId).
			Str("tx_hash_hex", txHashHex).
			Msg("ERROR: Failed to convert transaction to CBOR")
		return "", "", fmt.Errorf("failed to convert transaction to CBOR: %w", err)
	}
	log.Debug().
		Str("function", "MakerConfirmReceipt").
		Str("order_id", order.OrderInfo.OrderId).
		Str("tx_hash_hex", txHashHex).
		Str("tx_cbor", cbor).
		Int("cbor_length", len(cbor)).
		Msg("Transaction converted to CBOR successfully")

	log.Info().
		Str("function", "MakerConfirmReceipt").
		Str("order_id", order.OrderInfo.OrderId).
		Str("tx_hash_hex", txHashHex).
		Msg("Function success: Returning CBOR and transaction hash")

		// Evaluate transaction
	log.Debug().
		Str("function", "MakerConfirmReceipt").
		Str("order_id", order.OrderInfo.OrderId).
		Str("tx_hash_hex", txHashHex).
		Msg("Evaluating transaction")
	evaluationResult, err := config.CHAIN_CTX.EvaluateTx(txByte)
	if err != nil {
		log.Error().
			Err(err).
			Str("function", "MakerConfirmReceipt").
			Str("operation", "evaluate_transaction").
			Str("order_id", order.OrderInfo.OrderId).
			Str("tx_hash_hex", txHashHex).
			Msg("ERROR: Transaction evaluation failed")
		return "", "", fmt.Errorf("transaction evaluation failed: %w", err)
	}
	if len(evaluationResult) == 0 {
		log.Error().
			Str("function", "MakerConfirmReceipt").
			Str("operation", "evaluate_transaction_empty_result").
			Str("order_id", order.OrderInfo.OrderId).
			Str("tx_hash_hex", txHashHex).
			Msg("ERROR: Transaction evaluation returned empty result")
		return "", "", fmt.Errorf("transaction evaluation failed: empty result")
	}
	log.Debug().
		Str("function", "MakerConfirmReceipt").
		Str("order_id", order.OrderInfo.OrderId).
		Str("tx_hash_hex", txHashHex).
		Int("evaluation_result_length", len(evaluationResult)).
		Msg("Transaction evaluation successful")
	return cbor, txHashHex, nil
}
