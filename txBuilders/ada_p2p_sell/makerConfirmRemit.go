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
// // The TradeState only changes to REMIT_CONFIRMED_STATUS when the maker confirms the remit.
// 	TradeState: constants.REMIT_CONFIRMED_STATUS,
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

// cborString, txHash, err := ada_p2p_sell.MakerConfirmRemit(orderInfo, config.GetAda_P2PSellAdminWallet())
func MakerConfirmRemit(order *model.Order, adminWallet *config.Wallet) (string, string, error) {
	log.Trace().
		Str("function", "MakerConfirmRemit").
		Str("order_id", order.OrderInfo.OrderId).
		Int64("order_amount", order.OrderInfo.OrderAmount).
		Str("maker_address", order.OrderInfo.MakerAddress.String()).
		Str("taker_address", order.OrderInfo.TakerAddress.String()).
		Str("trade_state", order.TradeState).
		Msg("Function entry: Maker confirming remit")

	defer func() {
		if r := recover(); r != nil {
			stack := debug.Stack()
			log.Panic().
				Interface("panic_value", r).
				Bytes("stack_trace", stack).
				Str("function", "MakerConfirmRemit").
				Str("order_id", order.OrderInfo.OrderId).
				Msg("Panic occurred in MakerConfirmRemit")
			panic(r)
		}
	}()

	log.Debug().
		Str("function", "MakerConfirmRemit").
		Str("order_id", order.OrderInfo.OrderId).
		Msg("Initializing Apollo backend")

	apolloBE := apollo.New(&config.CHAIN_CTX)
	apolloBE = apolloBE.SetWalletFromBech32(order.OrderInfo.MakerAddress.String())

	log.Debug().
		Str("function", "MakerConfirmRemit").
		Str("order_id", order.OrderInfo.OrderId).
		Str("maker_address", order.OrderInfo.MakerAddress.String()).
		Msg("Apollo backend initialized and wallet set from Bech32")

	makerCommittedAmount := order.OrderTxInfo.MakerFee + order.OrderTxInfo.CollateralAmount
	takerCommittedAmount := order.OrderInfo.OrderAmount + order.OrderTxInfo.TakerFee + order.OrderTxInfo.CollateralAmount
	totalAmount := makerCommittedAmount + takerCommittedAmount

	log.Trace().
		Str("function", "MakerConfirmRemit").
		Str("order_id", order.OrderInfo.OrderId).
		Int64("maker_fee", order.OrderTxInfo.MakerFee).
		Int64("collateral_amount", order.OrderTxInfo.CollateralAmount).
		Int64("maker_committed_amount", makerCommittedAmount).
		Int64("order_amount", order.OrderInfo.OrderAmount).
		Int64("taker_fee", order.OrderTxInfo.TakerFee).
		Int64("taker_committed_amount", takerCommittedAmount).
		Int64("total_amount", totalAmount).
		Msg("Calculated amounts for maker confirm remit")

	log.Debug().
		Str("function", "MakerConfirmRemit").
		Str("order_id", order.OrderInfo.OrderId).
		Int64("total_amount", totalAmount).
		Msg("Amount calculations completed")

	log.Trace().
		Str("function", "MakerConfirmRemit").
		Str("order_id", order.OrderInfo.OrderId).
		Msg("Marshaling order datum to Plutus")

	orderDatumMarshaled, err := plutusEncoder.MarshalPlutus(*order)
	if err != nil {
		log.Error().
			Err(err).
			Str("function", "MakerConfirmRemit").
			Str("operation", "marshal_order_datum").
			Str("order_id", order.OrderInfo.OrderId).
			Msg("Failed to marshal order datum to Plutus")
		return "", "", err
	}

	log.Debug().
		Str("function", "MakerConfirmRemit").
		Str("order_id", order.OrderInfo.OrderId).
		Msg("Order datum marshaled successfully")

	log.Trace().
		Str("function", "MakerConfirmRemit").
		Str("order_id", order.OrderInfo.OrderId).
		Str("collateral_tx_id", order.OrderTxInfo.CollateralUtxo.TxID).
		Int("collateral_tx_index", order.OrderTxInfo.CollateralUtxo.TxIDIndex).
		Msg("Retrieving collateral UTXO")

	collateralUtxo, err := config.CHAIN_CTX.GetUtxoFromRef(order.OrderTxInfo.CollateralUtxo.TxID, order.OrderTxInfo.CollateralUtxo.TxIDIndex)
	if err != nil {
		log.Error().
			Err(err).
			Str("function", "MakerConfirmRemit").
			Str("operation", "get_collateral_utxo").
			Str("order_id", order.OrderInfo.OrderId).
			Str("collateral_tx_id", order.OrderTxInfo.CollateralUtxo.TxID).
			Int("collateral_tx_index", order.OrderTxInfo.CollateralUtxo.TxIDIndex).
			Msg("Failed to retrieve collateral UTXO")
		return "", "", err
	}

	log.Debug().
		Str("function", "MakerConfirmRemit").
		Str("order_id", order.OrderInfo.OrderId).
		Msg("Collateral UTXO retrieved successfully")

	log.Trace().
		Str("function", "MakerConfirmRemit").
		Str("order_id", order.OrderInfo.OrderId).
		Str("order_tx_id", order.OrderTxInfo.OrderUtxo.TxID).
		Int("order_tx_index", order.OrderTxInfo.OrderUtxo.TxIDIndex).
		Msg("Retrieving order UTXO")

	orderUTxO, err := config.CHAIN_CTX.GetUtxoFromRef(order.OrderTxInfo.OrderUtxo.TxID, order.OrderTxInfo.OrderUtxo.TxIDIndex)
	if err != nil {
		log.Error().
			Err(err).
			Str("function", "MakerConfirmRemit").
			Str("operation", "get_order_utxo").
			Str("order_id", order.OrderInfo.OrderId).
			Str("order_tx_id", order.OrderTxInfo.OrderUtxo.TxID).
			Int("order_tx_index", order.OrderTxInfo.OrderUtxo.TxIDIndex).
			Msg("Failed to retrieve order UTXO")
		return "", "", err
	}

	log.Debug().
		Str("function", "MakerConfirmRemit").
		Str("order_id", order.OrderInfo.OrderId).
		Msg("Order UTXO retrieved successfully")

	log.Trace().
		Str("function", "MakerConfirmRemit").
		Str("order_id", order.OrderInfo.OrderId).
		Int("user_utxos_count", len(order.OrderTxInfo.UserUtxos)).
		Msg("Retrieving user UTXOs")

	userUtxos, err := utility.GetUserUTxOs(order.OrderTxInfo.UserUtxos)
	if err != nil {
		log.Error().
			Err(err).
			Str("function", "MakerConfirmRemit").
			Str("operation", "get_user_utxos").
			Str("order_id", order.OrderInfo.OrderId).
			Int("user_utxos_count", len(order.OrderTxInfo.UserUtxos)).
			Msg("Failed to retrieve user UTXOs")
		return "", "", err
	}

	log.Debug().
		Str("function", "MakerConfirmRemit").
		Str("order_id", order.OrderInfo.OrderId).
		Int("user_utxos_retrieved", len(userUtxos)).
		Msg("User UTXOs retrieved successfully")

	log.Trace().
		Str("function", "MakerConfirmRemit").
		Str("order_id", order.OrderInfo.OrderId).
		Msg("Retrieving last block slot")

	lastSlot, err := config.CHAIN_CTX.LastBlockSlot()
	if err != nil {
		log.Error().
			Err(err).
			Str("function", "MakerConfirmRemit").
			Str("operation", "get_last_block_slot").
			Str("order_id", order.OrderInfo.OrderId).
			Msg("Failed to retrieve last block slot")
		return "", "", err
	}

	log.Debug().
		Str("function", "MakerConfirmRemit").
		Str("order_id", order.OrderInfo.OrderId).
		Int("last_slot", lastSlot).
		Int("ttl", lastSlot+300).
		Msg("Last block slot retrieved, TTL calculated")

	log.Info().
		Str("function", "MakerConfirmRemit").
		Str("order_id", order.OrderInfo.OrderId).
		Msg("Building transaction with Apollo")

	log.Trace().
		Str("function", "MakerConfirmRemit").
		Str("order_id", order.OrderInfo.OrderId).
		Str("change_address", order.OrderTxInfo.ChangeAddress.String()).
		Str("escrow_contract_address", order.OrderTxInfo.EscrowContractAddress.String()).
		Int64("total_amount", totalAmount).
		Str("state_token_policy_id", order.OrderTxInfo.StateTokenPolicyId).
		Int("ttl", lastSlot+300).
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
		AddRequiredSigner(serialization.PubKeyHash(order.OrderInfo.MakerAddress.PaymentPart)).
		SetTtl(int64(lastSlot) + 300).
		Complete()

	if err != nil {
		log.Error().
			Err(err).
			Str("function", "MakerConfirmRemit").
			Str("operation", "build_transaction").
			Str("order_id", order.OrderInfo.OrderId).
			Msg("Failed to build transaction")
		return "", "", err
	}

	log.Info().
		Str("function", "MakerConfirmRemit").
		Str("order_id", order.OrderInfo.OrderId).
		Msg("Transaction built successfully, signing with admin wallet")

	apolloBE, err = apolloBE.SignWithSkey(adminWallet.Vkey, adminWallet.Skey)
	if err != nil {
		log.Error().
			Err(err).
			Str("function", "MakerConfirmRemit").
			Str("operation", "sign_transaction").
			Str("order_id", order.OrderInfo.OrderId).
			Msg("Failed to sign transaction with admin wallet")
		return "", "", err
	}

	log.Debug().
		Str("function", "MakerConfirmRemit").
		Str("order_id", order.OrderInfo.OrderId).
		Msg("Transaction signed successfully")

	tx := apolloBE.GetTx()
	txHash, err := tx.TransactionBody.Hash()
	if err != nil {
		log.Error().
			Err(err).
			Str("function", "MakerConfirmRemit").
			Str("operation", "hash_transaction").
			Str("order_id", order.OrderInfo.OrderId).
			Msg("Failed to hash transaction body")
		return "", "", err
	}

	log.Trace().
		Str("function", "MakerConfirmRemit").
		Str("order_id", order.OrderInfo.OrderId).
		Str("tx_hash", hex.EncodeToString(txHash)).
		Msg("Transaction hash calculated")

	txByte, err := tx.Bytes()
	if err != nil {
		log.Error().
			Err(err).
			Str("function", "MakerConfirmRemit").
			Str("operation", "serialize_transaction").
			Str("order_id", order.OrderInfo.OrderId).
			Str("tx_hash", hex.EncodeToString(txHash)).
			Msg("Failed to serialize transaction to bytes")
		return "", "", err
	}

	log.Trace().
		Str("function", "MakerConfirmRemit").
		Str("order_id", order.OrderInfo.OrderId).
		Str("tx_hash", hex.EncodeToString(txHash)).
		Int("tx_bytes_length", len(txByte)).
		Msg("Transaction serialized to bytes")

	log.Debug().
		Str("function", "MakerConfirmRemit").
		Str("order_id", order.OrderInfo.OrderId).
		Str("tx_hash", hex.EncodeToString(txHash)).
		Msg("Evaluating transaction")

	evaluationResult, err := config.CHAIN_CTX.EvaluateTx(txByte)
	if err != nil {
		log.Error().
			Err(err).
			Str("function", "MakerConfirmRemit").
			Str("operation", "evaluate_transaction").
			Str("order_id", order.OrderInfo.OrderId).
			Str("tx_hash", hex.EncodeToString(txHash)).
			Msg("Failed to evaluate transaction")
		return "", "", err
	}

	log.Debug().
		Str("function", "MakerConfirmRemit").
		Str("order_id", order.OrderInfo.OrderId).
		Str("tx_hash", hex.EncodeToString(txHash)).
		Int("evaluation_result_length", len(evaluationResult)).
		Msg("Transaction evaluated successfully")

	if len(evaluationResult) == 0 {
		log.Error().
			Str("function", "MakerConfirmRemit").
			Str("operation", "validate_evaluation").
			Str("order_id", order.OrderInfo.OrderId).
			Str("tx_hash", hex.EncodeToString(txHash)).
			Msg("Transaction evaluation returned empty result")
		return "", "", fmt.Errorf("transaction evaluation failed")
	}

	log.Trace().
		Str("function", "MakerConfirmRemit").
		Str("order_id", order.OrderInfo.OrderId).
		Str("tx_hash", hex.EncodeToString(txHash)).
		Msg("Converting transaction to CBOR")

	cbor, err := Utils.ToCbor(tx)
	if err != nil {
		log.Error().
			Err(err).
			Str("function", "MakerConfirmRemit").
			Str("operation", "convert_to_cbor").
			Str("order_id", order.OrderInfo.OrderId).
			Str("tx_hash", hex.EncodeToString(txHash)).
			Msg("Failed to convert transaction to CBOR")
		return "", "", err
	}

	log.Info().
		Str("function", "MakerConfirmRemit").
		Str("order_id", order.OrderInfo.OrderId).
		Str("tx_id", hex.EncodeToString(txHash)).
		Msg("Maker confirmed remit successfully")

	return cbor, hex.EncodeToString(txHash), nil

}
