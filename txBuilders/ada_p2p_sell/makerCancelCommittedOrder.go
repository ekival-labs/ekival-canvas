package ada_p2p_sell

import (
	"encoding/hex"
	"fmt"
	"runtime/debug"

	"ekival-canvas/config"
	"ekival-canvas/constants"
	"ekival-canvas/model"
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

// makerFee := utility.CalculateFee(database.Precision, database.OrderAmount, database.OrderThreshold, database.MakerPct, database.MakerMinFee)
// takerFee := utility.CalculateFee(database.Precision, database.OrderAmount, database.OrderThreshold, database.TakerPct, database.TakerMinFee)
// collateralAmount := utility.CalculateFee(database.Precision, database.OrderAmount, database.OrderThreshold, database.CollateralPct, database.MinCollateral)
// cancelFee := utility.CalculateFee(database.Precision, database.OrderAmount, database.OrderThreshold, database.CancelPct, database.CancelMinFee)

// orderInfo := &viewmodel.Order{
// 	OrderInfo: viewmodel.OrderInfo{
// 		OrderId:            database.OrderId,
// 		OrderAmount:        database.OrderAmount,
// 		MakerAddress:       ma,
// 		TakerAddress:       ta,
// 	},
// 	BrokerageInfo: viewmodel.BrokerageInfo{
// 		CancelPenalty:  database.CancelPenalty,
// 	},
// 	OrderTxInfo: viewmodel.OrderTxInfo{
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
// 		TakerFee:         takerFee,
// 		CollateralAmount: collateralAmount,
// 		CancelFee:        cancelFee,
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

// cborString, txHash, err := ada_p2p_sell.MakerCancelCommittedOrder(orderInfo, treasuryInfo, config.GetAda_P2PSellAdminWallet())
func MakerCancelCommittedOrder(order *model.Order, treasury *model.TreasuryInfo, adminWallet *config.Wallet) (string, string, error) {
	log.Trace().
		Str("function", "MakerCancelCommittedOrder").
		Str("order_id", order.OrderInfo.OrderId).
		Int64("order_amount", order.OrderInfo.OrderAmount).
		Str("maker_address", order.OrderInfo.MakerAddress.String()).
		Str("taker_address", order.OrderInfo.TakerAddress.String()).
		Str("treasury_address", treasury.Address.String()).
		Str("trade_state", order.TradeState).
		Msg("Function entry: Canceling committed maker order")

	defer func() {
		if r := recover(); r != nil {
			stack := debug.Stack()
			log.Panic().
				Interface("panic_value", r).
				Bytes("stack_trace", stack).
				Str("function", "MakerCancelCommittedOrder").
				Str("order_id", order.OrderInfo.OrderId).
				Msg("Panic occurred in MakerCancelCommittedOrder")
			panic(r)
		}
	}()

	log.Debug().
		Str("function", "MakerCancelCommittedOrder").
		Str("order_id", order.OrderInfo.OrderId).
		Msg("Initializing Apollo backend")

	apolloBE := apollo.New(&config.CHAIN_CTX)
	apolloBE = apolloBE.SetWalletFromBech32(order.OrderInfo.MakerAddress.String())

	log.Debug().
		Str("function", "MakerCancelCommittedOrder").
		Str("order_id", order.OrderInfo.OrderId).
		Str("maker_address", order.OrderInfo.MakerAddress.String()).
		Msg("Apollo backend initialized and wallet set from Bech32")

	payToMakerAmount := order.OrderTxInfo.MakerFee + order.OrderTxInfo.CollateralAmount - order.BrokerageInfo.CancelPenalty - order.OrderTxInfo.CancelFee
	payToTakerAmount := order.OrderInfo.OrderAmount + order.OrderTxInfo.TakerFee + order.OrderTxInfo.CollateralAmount + order.BrokerageInfo.CancelPenalty

	log.Trace().
		Str("function", "MakerCancelCommittedOrder").
		Str("order_id", order.OrderInfo.OrderId).
		Int64("maker_fee", order.OrderTxInfo.MakerFee).
		Int64("collateral_amount", order.OrderTxInfo.CollateralAmount).
		Int64("cancel_penalty", order.BrokerageInfo.CancelPenalty).
		Int64("cancel_fee", order.OrderTxInfo.CancelFee).
		Int64("pay_to_maker_amount", payToMakerAmount).
		Int64("taker_fee", order.OrderTxInfo.TakerFee).
		Int64("order_amount", order.OrderInfo.OrderAmount).
		Int64("pay_to_taker_amount", payToTakerAmount).
		Msg("Calculated payment amounts for committed order cancellation")

	log.Debug().
		Str("function", "MakerCancelCommittedOrder").
		Str("order_id", order.OrderInfo.OrderId).
		Int64("pay_to_maker_amount", payToMakerAmount).
		Int64("pay_to_taker_amount", payToTakerAmount).
		Msg("Amount calculations completed")

	log.Trace().
		Str("function", "MakerCancelCommittedOrder").
		Str("order_id", order.OrderInfo.OrderId).
		Str("collateral_tx_id", order.OrderTxInfo.CollateralUtxo.TxID).
		Int("collateral_tx_index", order.OrderTxInfo.CollateralUtxo.TxIDIndex).
		Msg("Retrieving collateral UTXO")

	collateralUtxo, err := config.CHAIN_CTX.GetUtxoFromRef(order.OrderTxInfo.CollateralUtxo.TxID, order.OrderTxInfo.CollateralUtxo.TxIDIndex)
	if err != nil {
		log.Error().
			Err(err).
			Str("function", "MakerCancelCommittedOrder").
			Str("operation", "get_collateral_utxo").
			Str("order_id", order.OrderInfo.OrderId).
			Str("collateral_tx_id", order.OrderTxInfo.CollateralUtxo.TxID).
			Int("collateral_tx_index", order.OrderTxInfo.CollateralUtxo.TxIDIndex).
			Msg("Failed to retrieve collateral UTXO")
		return "", "", err
	}

	log.Debug().
		Str("function", "MakerCancelCommittedOrder").
		Str("order_id", order.OrderInfo.OrderId).
		Msg("Collateral UTXO retrieved successfully")

	log.Trace().
		Str("function", "MakerCancelCommittedOrder").
		Str("order_id", order.OrderInfo.OrderId).
		Str("order_tx_id", order.OrderTxInfo.OrderUtxo.TxID).
		Int("order_tx_index", order.OrderTxInfo.OrderUtxo.TxIDIndex).
		Msg("Retrieving order UTXO")

	orderUTxO, err := config.CHAIN_CTX.GetUtxoFromRef(order.OrderTxInfo.OrderUtxo.TxID, order.OrderTxInfo.OrderUtxo.TxIDIndex)
	if err != nil {
		log.Error().
			Err(err).
			Str("function", "MakerCancelCommittedOrder").
			Str("operation", "get_order_utxo").
			Str("order_id", order.OrderInfo.OrderId).
			Str("order_tx_id", order.OrderTxInfo.OrderUtxo.TxID).
			Int("order_tx_index", order.OrderTxInfo.OrderUtxo.TxIDIndex).
			Msg("Failed to retrieve order UTXO")
		return "", "", err
	}

	log.Debug().
		Str("function", "MakerCancelCommittedOrder").
		Str("order_id", order.OrderInfo.OrderId).
		Msg("Order UTXO retrieved successfully")

	log.Trace().
		Str("function", "MakerCancelCommittedOrder").
		Str("order_id", order.OrderInfo.OrderId).
		Int("user_utxos_count", len(order.OrderTxInfo.UserUtxos)).
		Msg("Retrieving user UTXOs")

	userUtxos, err := utility.GetUserUTxOs(order.OrderTxInfo.UserUtxos)
	if err != nil {
		log.Error().
			Err(err).
			Str("function", "MakerCancelCommittedOrder").
			Str("operation", "get_user_utxos").
			Str("order_id", order.OrderInfo.OrderId).
			Int("user_utxos_count", len(order.OrderTxInfo.UserUtxos)).
			Msg("Failed to retrieve user UTXOs")
		return "", "", err
	}

	log.Debug().
		Str("function", "MakerCancelCommittedOrder").
		Str("order_id", order.OrderInfo.OrderId).
		Int("user_utxos_retrieved", len(userUtxos)).
		Msg("User UTXOs retrieved successfully")

	log.Trace().
		Str("function", "MakerCancelCommittedOrder").
		Str("order_id", order.OrderInfo.OrderId).
		Msg("Retrieving last block slot")

	lastSlot, err := config.CHAIN_CTX.LastBlockSlot()
	if err != nil {
		log.Error().
			Err(err).
			Str("function", "MakerCancelCommittedOrder").
			Str("operation", "get_last_block_slot").
			Str("order_id", order.OrderInfo.OrderId).
			Msg("Failed to retrieve last block slot")
		return "", "", err
	}

	log.Debug().
		Str("function", "MakerCancelCommittedOrder").
		Str("order_id", order.OrderInfo.OrderId).
		Int("last_slot", lastSlot).
		Int("ttl", lastSlot+300).
		Int("validity_start", lastSlot).
		Msg("Last block slot retrieved, TTL and validity start calculated")

	log.Info().
		Str("function", "MakerCancelCommittedOrder").
		Str("order_id", order.OrderInfo.OrderId).
		Msg("Building transaction with Apollo")

	log.Trace().
		Str("function", "MakerCancelCommittedOrder").
		Str("order_id", order.OrderInfo.OrderId).
		Str("change_address", order.OrderTxInfo.ChangeAddress.String()).
		Str("treasury_address", treasury.Address.String()).
		Int64("cancel_fee", order.OrderTxInfo.CancelFee).
		Str("maker_address", order.OrderInfo.MakerAddress.String()).
		Int64("pay_to_maker_amount", payToMakerAmount).
		Str("taker_address", order.OrderInfo.TakerAddress.String()).
		Int64("pay_to_taker_amount", payToTakerAmount).
		Str("state_token_policy_id", order.OrderTxInfo.StateTokenPolicyId).
		Int("ttl", lastSlot+300).
		Int("validity_start", lastSlot).
		Msg("Transaction building parameters")

	apolloBE, err = apolloBE.
		SetChangeAddress(order.OrderTxInfo.ChangeAddress).
		AddCollateral(*collateralUtxo).
		AddLoadedUTxOs(userUtxos...).
		CollectFrom(*orderUTxO, *constants.INDEX_THREE_SPEND_REDEEMER).
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
			treasury.Address, treasury.Datum, int(order.OrderTxInfo.CancelFee), true,
		).
		PayToAddress(order.OrderInfo.MakerAddress, int(payToMakerAmount)).
		PayToAddress(order.OrderInfo.TakerAddress, int(payToTakerAmount)).
		AddRequiredSigner(adminWallet.AdminPKH).
		AddRequiredSigner(serialization.PubKeyHash(order.OrderInfo.MakerAddress.PaymentPart)).
		SetTtl(int64(lastSlot) + 300).
		SetValidityStart(int64(lastSlot)).
		Complete()

	if err != nil {
		log.Error().
			Err(err).
			Str("function", "MakerCancelCommittedOrder").
			Str("operation", "build_transaction").
			Str("order_id", order.OrderInfo.OrderId).
			Msg("Failed to build transaction")
		return "", "", err
	}

	log.Info().
		Str("function", "MakerCancelCommittedOrder").
		Str("order_id", order.OrderInfo.OrderId).
		Msg("Transaction built successfully, signing with admin wallet")

	apolloBE, err = apolloBE.SignWithSkey(adminWallet.AdminVkey, adminWallet.AdminSkey)
	if err != nil {
		log.Error().
			Err(err).
			Str("function", "MakerCancelCommittedOrder").
			Str("operation", "sign_transaction").
			Str("order_id", order.OrderInfo.OrderId).
			Msg("Failed to sign transaction with admin wallet")
		return "", "", err
	}

	log.Debug().
		Str("function", "MakerCancelCommittedOrder").
		Str("order_id", order.OrderInfo.OrderId).
		Msg("Transaction signed successfully")

	tx := apolloBE.GetTx()
	txHash, err := tx.TransactionBody.Hash()
	if err != nil {
		log.Error().
			Err(err).
			Str("function", "MakerCancelCommittedOrder").
			Str("operation", "hash_transaction").
			Str("order_id", order.OrderInfo.OrderId).
			Msg("Failed to hash transaction body")
		return "", "", err
	}

	log.Trace().
		Str("function", "MakerCancelCommittedOrder").
		Str("order_id", order.OrderInfo.OrderId).
		Str("tx_hash", hex.EncodeToString(txHash)).
		Msg("Transaction hash calculated")

	txByte, err := tx.Bytes()
	if err != nil {
		log.Error().
			Err(err).
			Str("function", "MakerCancelCommittedOrder").
			Str("operation", "serialize_transaction").
			Str("order_id", order.OrderInfo.OrderId).
			Str("tx_hash", hex.EncodeToString(txHash)).
			Msg("Failed to serialize transaction to bytes")
		return "", "", err
	}

	log.Trace().
		Str("function", "MakerCancelCommittedOrder").
		Str("order_id", order.OrderInfo.OrderId).
		Str("tx_hash", hex.EncodeToString(txHash)).
		Int("tx_bytes_length", len(txByte)).
		Msg("Transaction serialized to bytes")

	log.Debug().
		Str("function", "MakerCancelCommittedOrder").
		Str("order_id", order.OrderInfo.OrderId).
		Str("tx_hash", hex.EncodeToString(txHash)).
		Msg("Evaluating transaction")

	evaluationResult, err := config.CHAIN_CTX.EvaluateTx(txByte)
	if err != nil {
		log.Error().
			Err(err).
			Str("function", "MakerCancelCommittedOrder").
			Str("operation", "evaluate_transaction").
			Str("order_id", order.OrderInfo.OrderId).
			Str("tx_hash", hex.EncodeToString(txHash)).
			Msg("Failed to evaluate transaction")
		return "", "", err
	}

	log.Debug().
		Str("function", "MakerCancelCommittedOrder").
		Str("order_id", order.OrderInfo.OrderId).
		Str("tx_hash", hex.EncodeToString(txHash)).
		Int("evaluation_result_length", len(evaluationResult)).
		Msg("Transaction evaluated successfully")

	if len(evaluationResult) == 0 {
		log.Error().
			Str("function", "MakerCancelCommittedOrder").
			Str("operation", "validate_evaluation").
			Str("order_id", order.OrderInfo.OrderId).
			Str("tx_hash", hex.EncodeToString(txHash)).
			Msg("Transaction evaluation returned empty result")
		return "", "", fmt.Errorf("transaction evaluation failed")
	}

	log.Trace().
		Str("function", "MakerCancelCommittedOrder").
		Str("order_id", order.OrderInfo.OrderId).
		Str("tx_hash", hex.EncodeToString(txHash)).
		Msg("Converting transaction to CBOR")

	cbor, err := Utils.ToCbor(tx)
	if err != nil {
		log.Error().
			Err(err).
			Str("function", "MakerCancelCommittedOrder").
			Str("operation", "convert_to_cbor").
			Str("order_id", order.OrderInfo.OrderId).
			Str("tx_hash", hex.EncodeToString(txHash)).
			Msg("Failed to convert transaction to CBOR")
		return "", "", err
	}

	log.Info().
		Str("function", "MakerCancelCommittedOrder").
		Str("order_id", order.OrderInfo.OrderId).
		Str("tx_id", hex.EncodeToString(txHash)).
		Msg("Maker committed order canceled successfully")

	return cbor, hex.EncodeToString(txHash), nil

}
