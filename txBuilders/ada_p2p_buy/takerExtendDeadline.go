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

	"github.com/rs/zerolog/log"
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

// 	escrowContractAddress, err := Address.DecodeAddress(database.Ada_P2PBuyEscrow.Address)
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
// 			TxID:      database.Ada_P2PBuyEscrow.RefTxID,
// 			TxIDIndex: database.Ada_P2PBuyEscrow.RefTxIDIndex,
// 		},
// 		StateTokenPolicyId: database.Ada_BST.PolicyID,
// 		StateTokenRefUtxo: model.EUTxO{
// 			TxID:      database.Ada_BST.RefTxID,
// 			TxIDIndex: database.Ada_BST.RefTxIDIndex,
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
// Datum:   utility.CreateSimpleDatum(constants.INDEX_ONE, constants.ADA_P2P_BUY_FEE_TYPE),
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

// cborString, txHash, err := ada_p2p_buy.TakerExtendDeadline(orderInfo, treasuryInfo, extendPenalty, config.GetAda_P2PBuyAdminWallet())
func TakerExtendDeadline(order *model.Order, treasury *model.TreasuryInfo, extend_penalty int64, adminWallet *config.Wallet) (string, string, error) {
	log.Trace().
		Str("function", "TakerExtendDeadline").
		Str("order_id", order.OrderInfo.OrderId).
		Int64("order_amount", order.OrderInfo.OrderAmount).
		Str("taker_address", order.OrderInfo.TakerAddress.String()).
		Str("trade_state", order.TradeState).
		Int64("extend_penalty", extend_penalty).
		Int64("maker_deadline", order.OrderInfo.MakerDeadline).
		Int64("taker_deadline", order.OrderInfo.TakerDeadline).
		Msg("Function entry: Extending deadline for taker order")

	defer func() {
		if r := recover(); r != nil {
			stack := debug.Stack()
			log.Panic().
				Interface("panic_value", r).
				Bytes("stack_trace", stack).
				Str("function", "TakerExtendDeadline").
				Str("order_id", order.OrderInfo.OrderId).
				Msg("Panic occurred in TakerExtendDeadline")
			panic(r)
		}
	}()

	log.Debug().
		Str("function", "TakerExtendDeadline").
		Str("order_id", order.OrderInfo.OrderId).
		Msg("Initializing Apollo backend")

	apolloBE := apollo.New(&config.CHAIN_CTX)
	apolloBE = apolloBE.SetWalletFromBech32(order.OrderInfo.TakerAddress.String())

	log.Debug().
		Str("function", "TakerExtendDeadline").
		Str("order_id", order.OrderInfo.OrderId).
		Str("taker_address", order.OrderInfo.TakerAddress.String()).
		Msg("Apollo backend initialized and wallet set from Bech32")

	log.Trace().
		Str("function", "TakerExtendDeadline").
		Str("order_id", order.OrderInfo.OrderId).
		Msg("Marshaling order datum to Plutus")

	orderDatumMarshaled, err := plutusEncoder.MarshalPlutus(*order)
	if err != nil {
		log.Error().
			Err(err).
			Str("function", "TakerExtendDeadline").
			Str("operation", "marshal_order_datum").
			Str("order_id", order.OrderInfo.OrderId).
			Msg("Failed to marshal order datum to Plutus")
		return "", "", err
	}

	log.Debug().
		Str("function", "TakerExtendDeadline").
		Str("order_id", order.OrderInfo.OrderId).
		Msg("Order datum marshaled successfully")

	var amount int64 = 0

	log.Trace().
		Str("function", "TakerExtendDeadline").
		Str("order_id", order.OrderInfo.OrderId).
		Str("trade_state", order.TradeState).
		Msg("Calculating amount based on trade state")

	switch order.TradeState {
	case constants.UNCOMMITTED_ORDER_STATUS:
		amount = order.OrderInfo.OrderAmount + order.OrderTxInfo.MakerFee + order.OrderTxInfo.CollateralAmount
		log.Debug().
			Str("function", "TakerExtendDeadline").
			Str("order_id", order.OrderInfo.OrderId).
			Str("trade_state", constants.UNCOMMITTED_ORDER_STATUS).
			Int64("amount", amount).
			Msg("Using uncommitted order amount calculation")
	default:
		makerCommittedAmount := order.OrderInfo.OrderAmount + order.OrderTxInfo.MakerFee + order.OrderTxInfo.CollateralAmount
		takerCommittingAmount := order.OrderTxInfo.TakerFee + order.OrderTxInfo.CollateralAmount
		amount = makerCommittedAmount + takerCommittingAmount
		log.Debug().
			Str("function", "TakerExtendDeadline").
			Str("order_id", order.OrderInfo.OrderId).
			Str("trade_state", order.TradeState).
			Int64("maker_committed_amount", makerCommittedAmount).
			Int64("taker_committing_amount", takerCommittingAmount).
			Int64("total_amount", amount).
			Msg("Using committed order amount calculation")
	}

	var extendDeadlinesRedeemer Redeemer.Redeemer = Redeemer.Redeemer{}

	log.Trace().
		Str("function", "TakerExtendDeadline").
		Str("order_id", order.OrderInfo.OrderId).
		Int64("extend_penalty", extend_penalty).
		Msg("Selecting redeemer based on extend penalty")

	switch extend_penalty {
	case 0:
		extendDeadlinesRedeemer = *constants.EXTEND_DEADLINES_WITH_NO_PENALTY_REDEEMER
		log.Debug().
			Str("function", "TakerExtendDeadline").
			Str("order_id", order.OrderInfo.OrderId).
			Msg("Using no-penalty redeemer for deadline extension")

	default:
		log.Debug().
			Str("function", "TakerExtendDeadline").
			Str("order_id", order.OrderInfo.OrderId).
			Int64("extend_penalty", extend_penalty).
			Msg("Creating penalty redeemer for deadline extension")

		extendDeadlinesWithPenaltyRedeemerUnmarshaled := model.ExtendDeadlineWithPenaltyRedeemer{
			ExtendingPenalty: model.ExtendingPenalty{
				Penalty: extend_penalty,
			},
		}

		extendDeadlinesWithPenaltyRedeemerMarshaled, err := plutusEncoder.MarshalPlutus(extendDeadlinesWithPenaltyRedeemerUnmarshaled)
		if err != nil {
			log.Error().
				Err(err).
				Str("function", "TakerExtendDeadline").
				Str("operation", "marshal_penalty_redeemer").
				Str("order_id", order.OrderInfo.OrderId).
				Int64("extend_penalty", extend_penalty).
				Msg("Failed to marshal penalty redeemer")
			return "", "", err
		}

		extendDeadlinesRedeemer = Redeemer.Redeemer{
			Tag:   Redeemer.SPEND,
			Index: 0,
			Data:  *extendDeadlinesWithPenaltyRedeemerMarshaled,
		}

		log.Debug().
			Str("function", "TakerExtendDeadline").
			Str("order_id", order.OrderInfo.OrderId).
			Int64("extend_penalty", extend_penalty).
			Str("treasury_address", treasury.Address.String()).
			Msg("Adding penalty payment to treasury")

		apolloBE = apolloBE.PayToContract(
			treasury.Address, treasury.Datum, int(extend_penalty), true,
		)
	}

	log.Trace().
		Str("function", "TakerExtendDeadline").
		Str("order_id", order.OrderInfo.OrderId).
		Str("collateral_tx_id", order.OrderTxInfo.CollateralUtxo.TxID).
		Int("collateral_tx_index", order.OrderTxInfo.CollateralUtxo.TxIDIndex).
		Msg("Retrieving collateral UTXO")

	collateralUtxo, err := config.CHAIN_CTX.GetUtxoFromRef(order.OrderTxInfo.CollateralUtxo.TxID, order.OrderTxInfo.CollateralUtxo.TxIDIndex)
	if err != nil {
		log.Error().
			Err(err).
			Str("function", "TakerExtendDeadline").
			Str("operation", "get_collateral_utxo").
			Str("order_id", order.OrderInfo.OrderId).
			Str("collateral_tx_id", order.OrderTxInfo.CollateralUtxo.TxID).
			Int("collateral_tx_index", order.OrderTxInfo.CollateralUtxo.TxIDIndex).
			Msg("Failed to retrieve collateral UTXO")
		return "", "", err
	}

	log.Debug().
		Str("function", "TakerExtendDeadline").
		Str("order_id", order.OrderInfo.OrderId).
		Msg("Collateral UTXO retrieved successfully")

	log.Trace().
		Str("function", "TakerExtendDeadline").
		Str("order_id", order.OrderInfo.OrderId).
		Str("order_tx_id", order.OrderTxInfo.OrderUtxo.TxID).
		Int("order_tx_index", order.OrderTxInfo.OrderUtxo.TxIDIndex).
		Msg("Retrieving order UTXO")

	orderUTxO, err := config.CHAIN_CTX.GetUtxoFromRef(order.OrderTxInfo.OrderUtxo.TxID, order.OrderTxInfo.OrderUtxo.TxIDIndex)
	if err != nil {
		log.Error().
			Err(err).
			Str("function", "TakerExtendDeadline").
			Str("operation", "get_order_utxo").
			Str("order_id", order.OrderInfo.OrderId).
			Str("order_tx_id", order.OrderTxInfo.OrderUtxo.TxID).
			Int("order_tx_index", order.OrderTxInfo.OrderUtxo.TxIDIndex).
			Msg("Failed to retrieve order UTXO")
		return "", "", err
	}

	log.Debug().
		Str("function", "TakerExtendDeadline").
		Str("order_id", order.OrderInfo.OrderId).
		Msg("Order UTXO retrieved successfully")

	log.Trace().
		Str("function", "TakerExtendDeadline").
		Str("order_id", order.OrderInfo.OrderId).
		Int("user_utxos_count", len(order.OrderTxInfo.UserUtxos)).
		Msg("Retrieving user UTXOs")

	userUtxos, err := utility.GetUserUTxOs(order.OrderTxInfo.UserUtxos)
	if err != nil {
		log.Error().
			Err(err).
			Str("function", "TakerExtendDeadline").
			Str("operation", "get_user_utxos").
			Str("order_id", order.OrderInfo.OrderId).
			Int("user_utxos_count", len(order.OrderTxInfo.UserUtxos)).
			Msg("Failed to retrieve user UTXOs")
		return "", "", err
	}

	log.Debug().
		Str("function", "TakerExtendDeadline").
		Str("order_id", order.OrderInfo.OrderId).
		Int("user_utxos_retrieved", len(userUtxos)).
		Msg("User UTXOs retrieved successfully")

	log.Trace().
		Str("function", "TakerExtendDeadline").
		Str("order_id", order.OrderInfo.OrderId).
		Msg("Retrieving last block slot")

	lastSlot, err := config.CHAIN_CTX.LastBlockSlot()
	if err != nil {
		log.Error().
			Err(err).
			Str("function", "TakerExtendDeadline").
			Str("operation", "get_last_block_slot").
			Str("order_id", order.OrderInfo.OrderId).
			Msg("Failed to retrieve last block slot")
		return "", "", err
	}

	log.Debug().
		Str("function", "TakerExtendDeadline").
		Str("order_id", order.OrderInfo.OrderId).
		Int("last_slot", lastSlot).
		Int("ttl", lastSlot+300).
		Msg("Last block slot retrieved, TTL calculated")

	log.Info().
		Str("function", "TakerExtendDeadline").
		Str("order_id", order.OrderInfo.OrderId).
		Msg("Building deadline extension transaction with Apollo")

	log.Trace().
		Str("function", "TakerExtendDeadline").
		Str("order_id", order.OrderInfo.OrderId).
		Str("change_address", order.OrderTxInfo.ChangeAddress.String()).
		Str("escrow_contract_address", order.OrderTxInfo.EscrowContractAddress.String()).
		Int64("amount", amount).
		Str("state_token_policy_id", order.OrderTxInfo.StateTokenPolicyId).
		Int64("extend_penalty", extend_penalty).
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
		AddRequiredSigner(adminWallet.AdminPKH).
		AddRequiredSigner(serialization.PubKeyHash(order.OrderInfo.TakerAddress.PaymentPart)).
		SetTtl(int64(lastSlot) + 300).
		Complete()

	if err != nil {
		log.Error().
			Err(err).
			Str("function", "TakerExtendDeadline").
			Str("operation", "build_transaction").
			Str("order_id", order.OrderInfo.OrderId).
			Msg("Failed to build deadline extension transaction")
		return "", "", err
	}

	log.Info().
		Str("function", "TakerExtendDeadline").
		Str("order_id", order.OrderInfo.OrderId).
		Msg("Transaction built successfully, signing with admin wallet")

	apolloBE, err = apolloBE.SignWithSkey(adminWallet.AdminVkey, adminWallet.AdminSkey)
	if err != nil {
		log.Error().
			Err(err).
			Str("function", "TakerExtendDeadline").
			Str("operation", "sign_transaction").
			Str("order_id", order.OrderInfo.OrderId).
			Msg("Failed to sign transaction with admin wallet")
		return "", "", err
	}

	log.Debug().
		Str("function", "TakerExtendDeadline").
		Str("order_id", order.OrderInfo.OrderId).
		Msg("Transaction signed successfully")

	tx := apolloBE.GetTx()
	txHash, err := tx.TransactionBody.Hash()
	if err != nil {
		log.Error().
			Err(err).
			Str("function", "TakerExtendDeadline").
			Str("operation", "hash_transaction").
			Str("order_id", order.OrderInfo.OrderId).
			Msg("Failed to hash transaction body")
		return "", "", err
	}

	log.Trace().
		Str("function", "TakerExtendDeadline").
		Str("order_id", order.OrderInfo.OrderId).
		Str("tx_hash", hex.EncodeToString(txHash)).
		Msg("Transaction hash calculated")

	txByte, err := tx.Bytes()
	if err != nil {
		log.Error().
			Err(err).
			Str("function", "TakerExtendDeadline").
			Str("operation", "serialize_transaction").
			Str("order_id", order.OrderInfo.OrderId).
			Str("tx_hash", hex.EncodeToString(txHash)).
			Msg("Failed to serialize transaction to bytes")
		return "", "", err
	}

	log.Trace().
		Str("function", "TakerExtendDeadline").
		Str("order_id", order.OrderInfo.OrderId).
		Str("tx_hash", hex.EncodeToString(txHash)).
		Int("tx_bytes_length", len(txByte)).
		Msg("Transaction serialized to bytes")

	log.Debug().
		Str("function", "TakerExtendDeadline").
		Str("order_id", order.OrderInfo.OrderId).
		Str("tx_hash", hex.EncodeToString(txHash)).
		Msg("Evaluating transaction")

	evaluationResult, err := config.CHAIN_CTX.EvaluateTx(txByte)
	if err != nil {
		log.Error().
			Err(err).
			Str("function", "TakerExtendDeadline").
			Str("operation", "evaluate_transaction").
			Str("order_id", order.OrderInfo.OrderId).
			Str("tx_hash", hex.EncodeToString(txHash)).
			Msg("Failed to evaluate transaction")
		return "", "", err
	}

	log.Debug().
		Str("function", "TakerExtendDeadline").
		Str("order_id", order.OrderInfo.OrderId).
		Str("tx_hash", hex.EncodeToString(txHash)).
		Int("evaluation_result_length", len(evaluationResult)).
		Msg("Transaction evaluated successfully")

	if len(evaluationResult) == 0 {
		log.Error().
			Str("function", "TakerExtendDeadline").
			Str("operation", "validate_evaluation").
			Str("order_id", order.OrderInfo.OrderId).
			Str("tx_hash", hex.EncodeToString(txHash)).
			Msg("Transaction evaluation returned empty result")
		return "", "", fmt.Errorf("transaction evaluation failed")
	}

	log.Trace().
		Str("function", "TakerExtendDeadline").
		Str("order_id", order.OrderInfo.OrderId).
		Str("tx_hash", hex.EncodeToString(txHash)).
		Msg("Converting transaction to CBOR")

	cbor, err := Utils.ToCbor(tx)
	if err != nil {
		log.Error().
			Err(err).
			Str("function", "TakerExtendDeadline").
			Str("operation", "convert_to_cbor").
			Str("order_id", order.OrderInfo.OrderId).
			Str("tx_hash", hex.EncodeToString(txHash)).
			Msg("Failed to convert transaction to CBOR")
		return "", "", err
	}

	log.Info().
		Str("function", "TakerExtendDeadline").
		Str("order_id", order.OrderInfo.OrderId).
		Str("tx_id", hex.EncodeToString(txHash)).
		Int64("extend_penalty", extend_penalty).
		Msg("Deadline extended successfully by taker")

	return cbor, hex.EncodeToString(txHash), nil

}
