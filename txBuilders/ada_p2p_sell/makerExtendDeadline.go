package ada_p2p_sell

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
	"github.com/Salvionied/apollo/serialization/Redeemer"
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

// 	escrowContractAddress, err := Address.DecodeAddress(database.Ada_P2PSellEscrow.Address)
// 	if err != nil {
// 		log.Println(err)
// 		return err
// 	}

// makerFee := utility.CalculateFee(database.Precision, database.OrderAmount, database.OrderThreshold, database.MakerPct, database.MakerMinFee)
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
// // TradeState msu be same as the before state
// 	TradeState: database.TradeState,
// 	OrderTxInfo: viewmodel.OrderTxInfo{
// 		EscrowContractAddress: escrowContractAddress,
// 		EscrowContractRefUtxo: model.EUTxO{
// 			TxID:      database.Ada_P2PSellEscrow.RefTxID,
// 			TxIDIndex: database.Ada_P2PSellEscrow.RefTxIDIndex,
// 		},
// 		StateTokenPolicyId: database.Ada_SST.PolicyID,
// 		StateTokenRefUtxo: model.EUTxO{
// 			TxID:      database.Ada_SST.RefTxID,
// 			TxIDIndex: database.Ada_SST.RefTxIDIndex,
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
// Datum:   utility.CreateSimpleDatum(constants.INDEX_ONE, constants.ADA_P2P_SELL_FEE_TYPE),
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

// cborString, txHash, err := ada_p2p_sell.MakerExtendDeadline(orderInfo, treasuryInfo, extendPenalty, config.GetAda_P2PSellAdminWallet())
func MakerExtendDeadline(order *model.Order, treasury *model.TreasuryInfo, extend_penalty int64, adminWallet *config.Wallet) (string, string, error) {
	log.Trace().
		Str("function", "MakerExtendDeadline").
		Str("order_id", order.OrderInfo.OrderId).
		Int64("order_amount", order.OrderInfo.OrderAmount).
		Str("maker_address", order.OrderInfo.MakerAddress.String()).
		Str("trade_state", order.TradeState).
		Int64("extend_penalty", extend_penalty).
		Int64("maker_deadline", order.OrderInfo.MakerDeadline).
		Int64("taker_deadline", order.OrderInfo.TakerDeadline).
		Msg("Function entry: Extending deadline for maker order")

	defer func() {
		if r := recover(); r != nil {
			stack := debug.Stack()
			log.Panic().
				Interface("panic_value", r).
				Bytes("stack_trace", stack).
				Str("function", "MakerExtendDeadline").
				Str("order_id", order.OrderInfo.OrderId).
				Msg("Panic occurred in MakerExtendDeadline")
			panic(r)
		}
	}()

	log.Debug().
		Str("function", "MakerExtendDeadline").
		Str("order_id", order.OrderInfo.OrderId).
		Msg("Initializing Apollo backend")

	apolloBE := apollo.New(&config.CHAIN_CTX)
	apolloBE = apolloBE.SetWalletFromBech32(order.OrderInfo.MakerAddress.String())

	log.Debug().
		Str("function", "MakerExtendDeadline").
		Str("order_id", order.OrderInfo.OrderId).
		Str("maker_address", order.OrderInfo.MakerAddress.String()).
		Msg("Apollo backend initialized and wallet set from Bech32")

	log.Trace().
		Str("function", "MakerExtendDeadline").
		Str("order_id", order.OrderInfo.OrderId).
		Msg("Marshaling order datum to Plutus")

	orderDatumMarshaled, err := plutusEncoder.MarshalPlutus(*order)
	if err != nil {
		log.Error().
			Err(err).
			Str("function", "MakerExtendDeadline").
			Str("operation", "marshal_order_datum").
			Str("order_id", order.OrderInfo.OrderId).
			Msg("Failed to marshal order datum to Plutus")
		return "", "", err
	}

	log.Debug().
		Str("function", "MakerExtendDeadline").
		Str("order_id", order.OrderInfo.OrderId).
		Msg("Order datum marshaled successfully")

	var amount int64 = 0

	switch order.TradeState {
	case constants.UNCOMMITTED_ORDER_STATUS:
		amount = order.OrderTxInfo.MakerFee + order.OrderTxInfo.CollateralAmount
		log.Trace().
			Str("function", "MakerExtendDeadline").
			Str("order_id", order.OrderInfo.OrderId).
			Str("trade_state", order.TradeState).
			Int64("maker_fee", order.OrderTxInfo.MakerFee).
			Int64("collateral_amount", order.OrderTxInfo.CollateralAmount).
			Int64("amount", amount).
			Msg("Calculated amount for uncommitted order")
	default:
		makerCommittedAmount := order.OrderTxInfo.MakerFee + order.OrderTxInfo.CollateralAmount
		takerCommittedAmount := order.OrderInfo.OrderAmount + order.OrderTxInfo.TakerFee + order.OrderTxInfo.CollateralAmount
		amount = makerCommittedAmount + takerCommittedAmount
		log.Trace().
			Str("function", "MakerExtendDeadline").
			Str("order_id", order.OrderInfo.OrderId).
			Str("trade_state", order.TradeState).
			Int64("maker_committed_amount", makerCommittedAmount).
			Int64("taker_committed_amount", takerCommittedAmount).
			Int64("total_amount", amount).
			Msg("Calculated amount for committed order")
	}

	log.Debug().
		Str("function", "MakerExtendDeadline").
		Str("order_id", order.OrderInfo.OrderId).
		Int64("amount", amount).
		Msg("Amount calculations completed")

	var extendDeadlinesRedeemer Redeemer.Redeemer = Redeemer.Redeemer{}

	switch extend_penalty {
	case 0:
		log.Trace().
			Str("function", "MakerExtendDeadline").
			Str("order_id", order.OrderInfo.OrderId).
			Msg("Using extend deadlines with no penalty redeemer")
		extendDeadlinesRedeemer = *constants.EXTEND_DEADLINES_WITH_NO_PENALTY_REDEEMER

	default:
		log.Trace().
			Str("function", "MakerExtendDeadline").
			Str("order_id", order.OrderInfo.OrderId).
			Int64("extend_penalty", extend_penalty).
			Msg("Creating extend deadlines with penalty redeemer")

		extendDeadlinesWithPenaltyRedeemerUnmarshaled := model.ExtendDeadlineWithPenaltyRedeemer{
			ExtendingPenalty: model.ExtendingPenalty{
				Penalty: extend_penalty,
			},
		}

		extendDeadlinesWithPenaltyRedeemerMarshaled, err := plutusEncoder.MarshalPlutus(extendDeadlinesWithPenaltyRedeemerUnmarshaled)
		if err != nil {
			log.Error().
				Err(err).
				Str("function", "MakerExtendDeadline").
				Str("operation", "marshal_extend_redeemer").
				Str("order_id", order.OrderInfo.OrderId).
				Int64("extend_penalty", extend_penalty).
				Msg("Failed to marshal extend deadline with penalty redeemer to Plutus")
			return "", "", err
		}

		extendDeadlinesRedeemer = Redeemer.Redeemer{
			Tag:   Redeemer.SPEND,
			Index: 0,
			Data:  *extendDeadlinesWithPenaltyRedeemerMarshaled,
		}

		log.Debug().
			Str("function", "MakerExtendDeadline").
			Str("order_id", order.OrderInfo.OrderId).
			Int64("extend_penalty", extend_penalty).
			Msg("Adding penalty payment to treasury")

		apolloBE = apolloBE.PayToContract(
			treasury.Address, treasury.Datum, int(extend_penalty), true,
		)
	}

	log.Trace().
		Str("function", "MakerExtendDeadline").
		Str("order_id", order.OrderInfo.OrderId).
		Str("collateral_tx_id", order.OrderTxInfo.CollateralUtxo.TxID).
		Int("collateral_tx_index", order.OrderTxInfo.CollateralUtxo.TxIDIndex).
		Msg("Retrieving collateral UTXO")

	collateralUtxo, err := config.CHAIN_CTX.GetUtxoFromRef(order.OrderTxInfo.CollateralUtxo.TxID, order.OrderTxInfo.CollateralUtxo.TxIDIndex)
	if err != nil {
		log.Error().
			Err(err).
			Str("function", "MakerExtendDeadline").
			Str("operation", "get_collateral_utxo").
			Str("order_id", order.OrderInfo.OrderId).
			Str("collateral_tx_id", order.OrderTxInfo.CollateralUtxo.TxID).
			Int("collateral_tx_index", order.OrderTxInfo.CollateralUtxo.TxIDIndex).
			Msg("Failed to retrieve collateral UTXO")
		return "", "", err
	}

	log.Debug().
		Str("function", "MakerExtendDeadline").
		Str("order_id", order.OrderInfo.OrderId).
		Msg("Collateral UTXO retrieved successfully")

	log.Trace().
		Str("function", "MakerExtendDeadline").
		Str("order_id", order.OrderInfo.OrderId).
		Str("order_tx_id", order.OrderTxInfo.OrderUtxo.TxID).
		Int("order_tx_index", order.OrderTxInfo.OrderUtxo.TxIDIndex).
		Msg("Retrieving order UTXO")

	orderUTxO, err := config.CHAIN_CTX.GetUtxoFromRef(order.OrderTxInfo.OrderUtxo.TxID, order.OrderTxInfo.OrderUtxo.TxIDIndex)
	if err != nil {
		log.Error().
			Err(err).
			Str("function", "MakerExtendDeadline").
			Str("operation", "get_order_utxo").
			Str("order_id", order.OrderInfo.OrderId).
			Str("order_tx_id", order.OrderTxInfo.OrderUtxo.TxID).
			Int("order_tx_index", order.OrderTxInfo.OrderUtxo.TxIDIndex).
			Msg("Failed to retrieve order UTXO")
		return "", "", err
	}

	log.Debug().
		Str("function", "MakerExtendDeadline").
		Str("order_id", order.OrderInfo.OrderId).
		Msg("Order UTXO retrieved successfully")

	log.Trace().
		Str("function", "MakerExtendDeadline").
		Str("order_id", order.OrderInfo.OrderId).
		Int("user_utxos_count", len(order.OrderTxInfo.UserUtxos)).
		Msg("Retrieving user UTXOs")

	userUtxos, err := utility.GetUserUTxOs(order.OrderTxInfo.UserUtxos)
	if err != nil {
		log.Error().
			Err(err).
			Str("function", "MakerExtendDeadline").
			Str("operation", "get_user_utxos").
			Str("order_id", order.OrderInfo.OrderId).
			Int("user_utxos_count", len(order.OrderTxInfo.UserUtxos)).
			Msg("Failed to retrieve user UTXOs")
		return "", "", err
	}

	log.Debug().
		Str("function", "MakerExtendDeadline").
		Str("order_id", order.OrderInfo.OrderId).
		Int("user_utxos_retrieved", len(userUtxos)).
		Msg("User UTXOs retrieved successfully")

	log.Trace().
		Str("function", "MakerExtendDeadline").
		Str("order_id", order.OrderInfo.OrderId).
		Msg("Retrieving last block slot")

	lastSlot, err := config.CHAIN_CTX.LastBlockSlot()
	if err != nil {
		log.Error().
			Err(err).
			Str("function", "MakerExtendDeadline").
			Str("operation", "get_last_block_slot").
			Str("order_id", order.OrderInfo.OrderId).
			Msg("Failed to retrieve last block slot")
		return "", "", err
	}

	log.Debug().
		Str("function", "MakerExtendDeadline").
		Str("order_id", order.OrderInfo.OrderId).
		Int("last_slot", lastSlot).
		Int("ttl", lastSlot+300).
		Msg("Last block slot retrieved, TTL calculated")

	log.Info().
		Str("function", "MakerExtendDeadline").
		Str("order_id", order.OrderInfo.OrderId).
		Msg("Building transaction with Apollo")

	log.Trace().
		Str("function", "MakerExtendDeadline").
		Str("order_id", order.OrderInfo.OrderId).
		Str("change_address", order.OrderTxInfo.ChangeAddress.String()).
		Str("escrow_contract_address", order.OrderTxInfo.EscrowContractAddress.String()).
		Int64("amount", amount).
		Str("state_token_policy_id", order.OrderTxInfo.StateTokenPolicyId).
		Int("ttl", lastSlot+300).
		Msg("Transaction building parameters")

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
			order.OrderTxInfo.EscrowContractAddress, orderDatumMarshaled, int(amount), true, apollo.Unit{
				PolicyId: order.OrderTxInfo.StateTokenPolicyId,
				Name:     order.OrderInfo.OrderId,
				Quantity: int(1),
			},
		).
		AddRequiredSigner(adminWallet.PKH).
		AddRequiredSigner(serialization.PubKeyHash(order.OrderInfo.MakerAddress.PaymentPart)).
		SetTtl(int64(lastSlot) + 300).
		Complete()

	if err != nil {
		log.Error().
			Err(err).
			Str("function", "MakerExtendDeadline").
			Str("operation", "build_transaction").
			Str("order_id", order.OrderInfo.OrderId).
			Msg("Failed to build transaction")
		return "", "", err
	}

	log.Info().
		Str("function", "MakerExtendDeadline").
		Str("order_id", order.OrderInfo.OrderId).
		Msg("Transaction built successfully, signing with admin wallet")

	apolloBE, err = apolloBE.SignWithSkey(adminWallet.Vkey, adminWallet.Skey)
	if err != nil {
		log.Error().
			Err(err).
			Str("function", "MakerExtendDeadline").
			Str("operation", "sign_transaction").
			Str("order_id", order.OrderInfo.OrderId).
			Msg("Failed to sign transaction with admin wallet")
		return "", "", err
	}

	log.Debug().
		Str("function", "MakerExtendDeadline").
		Str("order_id", order.OrderInfo.OrderId).
		Msg("Transaction signed successfully")

	tx := apolloBE.GetTx()
	txHash, err := tx.TransactionBody.Hash()
	if err != nil {
		log.Error().
			Err(err).
			Str("function", "MakerExtendDeadline").
			Str("operation", "hash_transaction").
			Str("order_id", order.OrderInfo.OrderId).
			Msg("Failed to hash transaction body")
		return "", "", err
	}

	log.Trace().
		Str("function", "MakerExtendDeadline").
		Str("order_id", order.OrderInfo.OrderId).
		Str("tx_hash", hex.EncodeToString(txHash)).
		Msg("Transaction hash calculated")

	txByte, err := tx.Bytes()
	if err != nil {
		log.Error().
			Err(err).
			Str("function", "MakerExtendDeadline").
			Str("operation", "serialize_transaction").
			Str("order_id", order.OrderInfo.OrderId).
			Str("tx_hash", hex.EncodeToString(txHash)).
			Msg("Failed to serialize transaction to bytes")
		return "", "", err
	}

	log.Trace().
		Str("function", "MakerExtendDeadline").
		Str("order_id", order.OrderInfo.OrderId).
		Str("tx_hash", hex.EncodeToString(txHash)).
		Int("tx_bytes_length", len(txByte)).
		Msg("Transaction serialized to bytes")

	log.Debug().
		Str("function", "MakerExtendDeadline").
		Str("order_id", order.OrderInfo.OrderId).
		Str("tx_hash", hex.EncodeToString(txHash)).
		Msg("Evaluating transaction")

	evaluationResult, err := config.CHAIN_CTX.EvaluateTx(txByte)
	if err != nil {
		log.Error().
			Err(err).
			Str("function", "MakerExtendDeadline").
			Str("operation", "evaluate_transaction").
			Str("order_id", order.OrderInfo.OrderId).
			Str("tx_hash", hex.EncodeToString(txHash)).
			Msg("Failed to evaluate transaction")
		return "", "", err
	}

	log.Debug().
		Str("function", "MakerExtendDeadline").
		Str("order_id", order.OrderInfo.OrderId).
		Str("tx_hash", hex.EncodeToString(txHash)).
		Int("evaluation_result_length", len(evaluationResult)).
		Msg("Transaction evaluated successfully")

	if len(evaluationResult) == 0 {
		log.Error().
			Str("function", "MakerExtendDeadline").
			Str("operation", "validate_evaluation").
			Str("order_id", order.OrderInfo.OrderId).
			Str("tx_hash", hex.EncodeToString(txHash)).
			Msg("Transaction evaluation returned empty result")
		return "", "", fmt.Errorf("transaction evaluation failed")
	}

	log.Trace().
		Str("function", "MakerExtendDeadline").
		Str("order_id", order.OrderInfo.OrderId).
		Str("tx_hash", hex.EncodeToString(txHash)).
		Msg("Converting transaction to CBOR")

	cbor, err := Utils.ToCbor(tx)
	if err != nil {
		log.Error().
			Err(err).
			Str("function", "MakerExtendDeadline").
			Str("operation", "convert_to_cbor").
			Str("order_id", order.OrderInfo.OrderId).
			Str("tx_hash", hex.EncodeToString(txHash)).
			Msg("Failed to convert transaction to CBOR")
		return "", "", err
	}

	log.Info().
		Str("function", "MakerExtendDeadline").
		Str("order_id", order.OrderInfo.OrderId).
		Str("tx_id", hex.EncodeToString(txHash)).
		Int64("new_maker_deadline", order.OrderInfo.MakerDeadline).
		Int64("new_taker_deadline", order.OrderInfo.TakerDeadline).
		Msg("Maker deadline extended successfully")

	return cbor, hex.EncodeToString(txHash), nil

}
