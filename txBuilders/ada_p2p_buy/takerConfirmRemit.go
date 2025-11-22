package ada_p2p_buy

import (
	"encoding/hex"
	"fmt"
	"runtime/debug"

	"ekival-canvas/config"
	"ekival-canvas/constants"
	"ekival-canvas/model"
	"ekival-canvas/plutusEncoder"
	"ekival-canvas/utility"

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

// 	escrowContractAddress, err := Address.DecodeAddress(database.Ada_P2PBuyEscrow.Address)
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
// 	},
// // The TradeState only changes to REMIT_CONFIRMED_STATUS when the taker confirms the remit.
// 	TradeState: constants.REMIT_CONFIRMED_STATUS,
// 	OrderTxInfo: viewmodel.OrderTxInfo{
// 		EscrowContractAddress: escrowContractAddress,
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

// cborString, txHash, err := ada_p2p_buy.TakerConfirmRemit(orderInfo, config.GetAda_P2PBuyAdminWallet())
func TakerConfirmRemit(order *model.Order, adminWallet *config.Wallet) (string, string, error) {
	log.Trace().
		Str("function", "TakerConfirmRemit").
		Str("order_id", order.OrderInfo.OrderId).
		Int64("order_amount", order.OrderInfo.OrderAmount).
		Str("maker_address", order.OrderInfo.MakerAddress.String()).
		Str("taker_address", order.OrderInfo.TakerAddress.String()).
		Str("trade_state", order.TradeState).
		Msg("Function entry: Taker confirming remit")

	defer func() {
		if r := recover(); r != nil {
			stack := debug.Stack()
			log.Panic().
				Interface("panic_value", r).
				Bytes("stack_trace", stack).
				Str("function", "TakerConfirmRemit").
				Str("order_id", func() string {
					if order != nil {
						return order.OrderInfo.OrderId
					}
					return "unknown"
				}()).
				Msg("Panic occurred in TakerConfirmRemit")
			panic(r)
		}
	}()

	log.Trace().
		Str("function", "TakerConfirmRemit").
		Msg("Validating input parameters")

	if order == nil {
		log.Error().
			Str("function", "TakerConfirmRemit").
			Str("operation", "validate_input").
			Msg("Order parameter is nil")
		return "", "", fmt.Errorf("order parameter is nil")
	}
	if adminWallet == nil {
		log.Error().
			Str("function", "TakerConfirmRemit").
			Str("operation", "validate_input").
			Str("order_id", order.OrderInfo.OrderId).
			Msg("Admin wallet parameter is nil")
		return "", "", fmt.Errorf("adminWallet parameter is nil")
	}

	log.Debug().
		Str("function", "TakerConfirmRemit").
		Str("order_id", order.OrderInfo.OrderId).
		Int64("order_amount", order.OrderInfo.OrderAmount).
		Str("maker_address", order.OrderInfo.MakerAddress.String()).
		Str("taker_address", order.OrderInfo.TakerAddress.String()).
		Int64("maker_fee", order.OrderTxInfo.MakerFee).
		Int64("taker_fee", order.OrderTxInfo.TakerFee).
		Int64("collateral_amount", order.OrderTxInfo.CollateralAmount).
		Str("change_address", order.OrderTxInfo.ChangeAddress.String()).
		Str("escrow_contract_address", order.OrderTxInfo.EscrowContractAddress.String()).
		Str("state_token_policy_id", order.OrderTxInfo.StateTokenPolicyId).
		Str("trade_state", order.TradeState).
		Interface("admin_pkh", adminWallet.PKH).
		Msg("Input parameters validated and logged")

	log.Debug().
		Str("function", "TakerConfirmRemit").
		Str("order_id", order.OrderInfo.OrderId).
		Msg("Initializing Apollo backend")

	apolloBE := apollo.New(&config.CHAIN_CTX)

	log.Debug().
		Str("function", "TakerConfirmRemit").
		Str("order_id", order.OrderInfo.OrderId).
		Str("taker_address", order.OrderInfo.TakerAddress.String()).
		Msg("Setting wallet from Bech32 address")

	apolloBE = apolloBE.SetWalletFromBech32(order.OrderInfo.TakerAddress.String())

	log.Debug().
		Str("function", "TakerConfirmRemit").
		Str("order_id", order.OrderInfo.OrderId).
		Str("taker_address", order.OrderInfo.TakerAddress.String()).
		Msg("Apollo backend initialized and wallet set from Bech32")

	log.Trace().
		Str("function", "TakerConfirmRemit").
		Str("order_id", order.OrderInfo.OrderId).
		Str("trade_state", order.TradeState).
		Int64("precision", order.BrokerageInfo.Precision).
		Int64("collateral_pct", order.BrokerageInfo.CollateralPct).
		Int64("maker_pct", order.BrokerageInfo.MakerPct).
		Int64("taker_pct", order.BrokerageInfo.TakerPct).
		Msg("Marshaling order datum to Plutus")

	orderDatumMarshaled, err := plutusEncoder.MarshalPlutus(*order)
	if err != nil {
		log.Error().
			Err(err).
			Str("function", "TakerConfirmRemit").
			Str("operation", "marshal_order_datum").
			Str("order_id", order.OrderInfo.OrderId).
			Str("trade_state", order.TradeState).
			Msg("Failed to marshal order datum to Plutus")
		return "", "", fmt.Errorf("failed to marshal order datum: %w", err)
	}

	log.Debug().
		Str("function", "TakerConfirmRemit").
		Str("order_id", order.OrderInfo.OrderId).
		Str("trade_state", order.TradeState).
		Msg("Order datum marshaled successfully, state transition: COMMITTED_ORDER → REMIT_CONFIRMED")

	log.Trace().
		Str("function", "TakerConfirmRemit").
		Str("order_id", order.OrderInfo.OrderId).
		Str("collateral_tx_id", order.OrderTxInfo.CollateralUtxo.TxID).
		Int("collateral_tx_index", order.OrderTxInfo.CollateralUtxo.TxIDIndex).
		Msg("Retrieving collateral UTXO")

	collateralUtxo, err := config.CHAIN_CTX.GetUtxoFromRef(order.OrderTxInfo.CollateralUtxo.TxID, order.OrderTxInfo.CollateralUtxo.TxIDIndex)
	if err != nil {
		log.Error().
			Err(err).
			Str("function", "TakerConfirmRemit").
			Str("operation", "get_collateral_utxo").
			Str("order_id", order.OrderInfo.OrderId).
			Str("collateral_tx_id", order.OrderTxInfo.CollateralUtxo.TxID).
			Int("collateral_tx_index", order.OrderTxInfo.CollateralUtxo.TxIDIndex).
			Msg("Failed to retrieve collateral UTXO")
		return "", "", fmt.Errorf("failed to get Collateral UTxO: %w", err)
	}

	log.Debug().
		Str("function", "TakerConfirmRemit").
		Str("order_id", order.OrderInfo.OrderId).
		Interface("collateral_utxo", collateralUtxo).
		Msg("Collateral UTXO retrieved successfully")

	log.Trace().
		Str("function", "TakerConfirmRemit").
		Str("order_id", order.OrderInfo.OrderId).
		Str("order_tx_id", order.OrderTxInfo.OrderUtxo.TxID).
		Int("order_tx_index", order.OrderTxInfo.OrderUtxo.TxIDIndex).
		Msg("Retrieving order UTXO")

	orderUTxO, err := config.CHAIN_CTX.GetUtxoFromRef(order.OrderTxInfo.OrderUtxo.TxID, order.OrderTxInfo.OrderUtxo.TxIDIndex)
	if err != nil {
		log.Error().
			Err(err).
			Str("function", "TakerConfirmRemit").
			Str("operation", "get_order_utxo").
			Str("order_id", order.OrderInfo.OrderId).
			Str("order_tx_id", order.OrderTxInfo.OrderUtxo.TxID).
			Int("order_tx_index", order.OrderTxInfo.OrderUtxo.TxIDIndex).
			Msg("Failed to retrieve order UTXO")
		return "", "", fmt.Errorf("failed to get Order UTxO: %w", err)
	}

	log.Debug().
		Str("function", "TakerConfirmRemit").
		Str("order_id", order.OrderInfo.OrderId).
		Interface("order_utxo", orderUTxO).
		Msg("Order UTXO retrieved successfully, contains committed funds from previous transaction")

	totalAmount := int64(orderUTxO.Output.GetAmount().GetCoin())

	log.Trace().
		Str("function", "TakerConfirmRemit").
		Str("order_id", order.OrderInfo.OrderId).
		Int64("order_utxo_amount", totalAmount).
		Msg("Calculated total amount from order UTXO")

	log.Debug().
		Str("function", "TakerConfirmRemit").
		Str("order_id", order.OrderInfo.OrderId).
		Int64("total_amount", totalAmount).
		Msg("Amount breakdown: Order UTXO contains OrderAmount + MakerFee + MakerCollateral + TakerFee + TakerCollateral, this transaction only changes state")

	log.Trace().
		Str("function", "TakerConfirmRemit").
		Str("order_id", order.OrderInfo.OrderId).
		Int("user_utxos_count", len(order.OrderTxInfo.UserUtxos)).
		Interface("user_utxos", order.OrderTxInfo.UserUtxos).
		Msg("Retrieving user UTXOs")

	userUtxos, err := utility.GetUserUTxOs(order.OrderTxInfo.UserUtxos)
	if err != nil {
		log.Error().
			Err(err).
			Str("function", "TakerConfirmRemit").
			Str("operation", "get_user_utxos").
			Str("order_id", order.OrderInfo.OrderId).
			Int("user_utxos_count", len(order.OrderTxInfo.UserUtxos)).
			Interface("user_utxos_input", order.OrderTxInfo.UserUtxos).
			Msg("Failed to retrieve user UTXOs")
		return "", "", fmt.Errorf("failed to get User UTxOs: %w", err)
	}

	log.Debug().
		Str("function", "TakerConfirmRemit").
		Str("order_id", order.OrderInfo.OrderId).
		Int("user_utxos_retrieved", len(userUtxos)).
		Msg("User UTXOs retrieved successfully")

	log.Trace().
		Str("function", "TakerConfirmRemit").
		Str("order_id", order.OrderInfo.OrderId).
		Msg("Retrieving last block slot")

	lastSlot, err := config.CHAIN_CTX.LastBlockSlot()
	if err != nil {
		log.Error().
			Err(err).
			Str("function", "TakerConfirmRemit").
			Str("operation", "get_last_block_slot").
			Str("order_id", order.OrderInfo.OrderId).
			Msg("Failed to retrieve last block slot")
		return "", "", fmt.Errorf("failed to get last block slot: %w", err)
	}

	ttl := int64(lastSlot) + 300

	log.Debug().
		Str("function", "TakerConfirmRemit").
		Str("order_id", order.OrderInfo.OrderId).
		Int("last_slot", lastSlot).
		Int64("ttl", ttl).
		Msg("Last block slot retrieved, TTL calculated")

	log.Info().
		Str("function", "TakerConfirmRemit").
		Str("order_id", order.OrderInfo.OrderId).
		Msg("Building transaction with Apollo")

	log.Trace().
		Str("function", "TakerConfirmRemit").
		Str("order_id", order.OrderInfo.OrderId).
		Str("change_address", order.OrderTxInfo.ChangeAddress.String()).
		Str("escrow_contract_address", order.OrderTxInfo.EscrowContractAddress.String()).
		Int64("total_amount", totalAmount).
		Str("state_token_policy_id", order.OrderTxInfo.StateTokenPolicyId).
		Str("order_token_name", order.OrderInfo.OrderId).
		Int("user_utxos_count", len(userUtxos)).
		Int64("ttl", ttl).
		Str("redeemer", "INDEX_FIVE_SPEND_REDEEMER").
		Msg("Transaction building parameters")

	apolloBE, err = apolloBE.
		SetChangeAddress(order.OrderTxInfo.ChangeAddress).
		AddCollateral(*collateralUtxo).
		AddLoadedUTxOs(userUtxos...).
		CollectFrom(*orderUTxO, *constants.INDEX_FIVE_SPEND_REDEEMER).
		AddReferenceInputV3(
			order.OrderTxInfo.EscrowContractRefUtxo.TxID,
			order.OrderTxInfo.EscrowContractRefUtxo.TxIDIndex,
		).
		PayToContract(
			order.OrderTxInfo.EscrowContractAddress, orderDatumMarshaled, int(totalAmount), true, apollo.Unit{
				PolicyId: order.OrderTxInfo.StateTokenPolicyId,
				Name:     order.OrderInfo.OrderId,
				Quantity: int(1),
			},
		).
		AddRequiredSigner(adminWallet.PKH).
		AddRequiredSigner(serialization.PubKeyHash(order.OrderInfo.TakerAddress.PaymentPart)).
		SetTtl(ttl).
		SetFeePadding(apolloBE.Fee + 100000).
		Complete()

	if err != nil {
		log.Error().
			Err(err).
			Str("function", "TakerConfirmRemit").
			Str("operation", "build_transaction").
			Str("order_id", order.OrderInfo.OrderId).
			Str("order_utxo", fmt.Sprintf("%s#%d", order.OrderTxInfo.OrderUtxo.TxID, order.OrderTxInfo.OrderUtxo.TxIDIndex)).
			Str("escrow_contract", order.OrderTxInfo.EscrowContractAddress.String()).
			Int64("total_amount", totalAmount).
			Str("trade_state", order.TradeState).
			Str("redeemer", "INDEX_FIVE_SPEND_REDEEMER").
			Int("user_utxos_count", len(userUtxos)).
			Str("collateral_utxo", fmt.Sprintf("%s#%d", order.OrderTxInfo.CollateralUtxo.TxID, order.OrderTxInfo.CollateralUtxo.TxIDIndex)).
			Msg("Failed to build transaction - 'could not estimate ExUnits' typically means: datum structure mismatch, invalid redeemer, transaction structure validation failed, or contract logic validation failed")
		return "", "", fmt.Errorf("transaction building failed: %w", err)
	}

	log.Info().
		Str("function", "TakerConfirmRemit").
		Str("order_id", order.OrderInfo.OrderId).
		Msg("Transaction built successfully")

	log.Info().
		Str("function", "TakerConfirmRemit").
		Str("order_id", order.OrderInfo.OrderId).
		Msg("Transaction built successfully, signing with admin wallet")

	apolloBE, err = apolloBE.SignWithSkey(adminWallet.Vkey, adminWallet.Skey)
	if err != nil {
		log.Error().
			Err(err).
			Str("function", "TakerConfirmRemit").
			Str("operation", "sign_transaction").
			Str("order_id", order.OrderInfo.OrderId).
			Msg("Failed to sign transaction with admin wallet")
		return "", "", fmt.Errorf("transaction signing failed: %w", err)
	}

	log.Debug().
		Str("function", "TakerConfirmRemit").
		Str("order_id", order.OrderInfo.OrderId).
		Msg("Transaction signed successfully")

	tx := apolloBE.GetTx()
	if tx == nil {
		log.Error().
			Str("function", "TakerConfirmRemit").
			Str("operation", "get_transaction").
			Str("order_id", order.OrderInfo.OrderId).
			Msg("Transaction is nil after GetTx()")
		return "", "", fmt.Errorf("transaction is nil after GetTx()")
	}

	log.Trace().
		Str("function", "TakerConfirmRemit").
		Str("order_id", order.OrderInfo.OrderId).
		Msg("Computing transaction hash")

	txHash, err := tx.TransactionBody.Hash()
	if err != nil {
		log.Error().
			Err(err).
			Str("function", "TakerConfirmRemit").
			Str("operation", "hash_transaction").
			Str("order_id", order.OrderInfo.OrderId).
			Msg("Failed to compute transaction hash")
		return "", "", fmt.Errorf("failed to compute transaction hash: %w", err)
	}

	txHashHex := hex.EncodeToString(txHash)

	log.Trace().
		Str("function", "TakerConfirmRemit").
		Str("order_id", order.OrderInfo.OrderId).
		Str("tx_hash", txHashHex).
		Msg("Transaction hash calculated")

	log.Trace().
		Str("function", "TakerConfirmRemit").
		Str("order_id", order.OrderInfo.OrderId).
		Str("tx_hash", txHashHex).
		Msg("Serializing transaction to bytes")

	txByte, err := tx.Bytes()
	if err != nil {
		log.Error().
			Err(err).
			Str("function", "TakerConfirmRemit").
			Str("operation", "serialize_transaction").
			Str("order_id", order.OrderInfo.OrderId).
			Str("tx_hash", txHashHex).
			Msg("Failed to serialize transaction to bytes")
		return "", "", fmt.Errorf("failed to serialize transaction to bytes: %w", err)
	}

	log.Trace().
		Str("function", "TakerConfirmRemit").
		Str("order_id", order.OrderInfo.OrderId).
		Str("tx_hash", txHashHex).
		Int("tx_bytes_length", len(txByte)).
		Msg("Transaction serialized to bytes")

	log.Debug().
		Str("function", "TakerConfirmRemit").
		Str("order_id", order.OrderInfo.OrderId).
		Str("tx_hash", txHashHex).
		Msg("Evaluating transaction")

	evaluationResult, err := config.CHAIN_CTX.EvaluateTx(txByte)
	if err != nil {
		log.Error().
			Err(err).
			Str("function", "TakerConfirmRemit").
			Str("operation", "evaluate_transaction").
			Str("order_id", order.OrderInfo.OrderId).
			Str("tx_hash", txHashHex).
			Msg("Failed to evaluate transaction")
		return "", "", fmt.Errorf("transaction evaluation failed: %w", err)
	}

	log.Debug().
		Str("function", "TakerConfirmRemit").
		Str("order_id", order.OrderInfo.OrderId).
		Str("tx_hash", txHashHex).
		Int("evaluation_result_length", len(evaluationResult)).
		Msg("Transaction evaluated successfully")

	if len(evaluationResult) == 0 {
		log.Error().
			Str("function", "TakerConfirmRemit").
			Str("operation", "validate_evaluation").
			Str("order_id", order.OrderInfo.OrderId).
			Str("tx_hash", txHashHex).
			Msg("Transaction evaluation returned empty result")
		return "", "", fmt.Errorf("transaction evaluation failed: empty result")
	}

	log.Trace().
		Str("function", "TakerConfirmRemit").
		Str("order_id", order.OrderInfo.OrderId).
		Str("tx_hash", txHashHex).
		Msg("Converting transaction to CBOR")

	cbor, err := Utils.ToCbor(tx)
	if err != nil {
		log.Error().
			Err(err).
			Str("function", "TakerConfirmRemit").
			Str("operation", "convert_to_cbor").
			Str("order_id", order.OrderInfo.OrderId).
			Str("tx_hash", txHashHex).
			Msg("Failed to convert transaction to CBOR")
		return "", "", fmt.Errorf("failed to convert transaction to CBOR: %w", err)
	}

	log.Info().
		Str("function", "TakerConfirmRemit").
		Str("order_id", order.OrderInfo.OrderId).
		Str("tx_id", txHashHex).
		Str("tx_cbor", cbor).
		Int("cbor_length", len(cbor)).
		Msg("Taker remit confirmed successfully")

	return cbor, txHashHex, nil

}
