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
	"github.com/Salvionied/apollo/txBuilding/Utils"
)

// cfg := config.GetGlobalConfig()

// 	ma, err := Address.DecodeAddress(database.MakerAddress)
// 	if err != nil {
// 		fiberLogger.Error(err)
// 		return err
// 	}

// 	ta, err := Address.DecodeAddress(database.TakerAddress)
// 	if err != nil {
// 		fiberLogger.Error(err)
// 		return err
// 	}

// 	changeAddress, err := Address.DecodeAddress(body.ChangeAddress)
// 	if err != nil {
// 		fiberLogger.Error(err)
// 		return err
// 	}

// 	escrowContractAddress, err := Address.DecodeAddress(database.Ada_P2PBuyEscrow.Address)
// 	if err != nil {
// 		fiberLogger.Error(err)
// 		return err
// 	}

// 	currentTime := time.Now().UnixNano() / int64(time.Millisecond)
// // MakerDeadline is whole Order Time in Hours for example 5 days = 5 * 24 * 60 * 60 * 1000
// 	var md int64 = currentTime + (MakerDeadline * 1000)
// // TakerDeadline less than MakerDeadline in Hours for example 4 day = 4 * 24 * 60 * 60 * 1000
// 	var td int64 = currentTime + (TakerDeadline * 1000)

// makerFee := utility.CalculateFee(database.Precision, database.OrderAmount, database.OrderThreshold, database.MakerPct, database.MakerMinFee)
// collateralAmount := utility.CalculateFee(database.Precision, database.OrderAmount, database.OrderThreshold, database.CollateralPct, database.MinCollateral)

// orderInfo := &viewmodel.Order{
// 	OrderInfo: viewmodel.OrderInfo{
// 		OrderId:            database.OrderId,
// 		OrderAmount:        body.OrderAmount,
// 		MakerAddress:       ma,
// 		TakerAddress:       ta,
// 		MakerDeadline:      md,
// 		TakerDeadline:      td,
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
// 	TradeState: constants.UNCOMMITTED_ORDER_STATUS,
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
// 	},
// }

// cborString, txHash, err := ada_p2p_buy.MakerCreateOrder(orderInfo, config.GetAda_BSTAdminWallet())
func MakerCreateOrder(order *model.Order, adminWallet *config.Wallet) (string, string, error) {
	log.Trace().
		Str("function", "MakerCreateOrder").
		Str("order_id", order.OrderInfo.OrderId).
		Int64("order_amount", order.OrderInfo.OrderAmount).
		Str("maker_address", order.OrderInfo.MakerAddress.String()).
		Str("taker_address", order.OrderInfo.TakerAddress.String()).
		Str("trade_state", order.TradeState).
		Msg("Function entry: Creating maker order")

	defer func() {
		if r := recover(); r != nil {
			stack := debug.Stack()
			log.Panic().
				Interface("panic_value", r).
				Bytes("stack_trace", stack).
				Str("function", "MakerCreateOrder").
				Str("order_id", order.OrderInfo.OrderId).
				Msg("Panic occurred in MakerCreateOrder")
			panic(r)
		}
	}()

	log.Debug().
		Str("function", "MakerCreateOrder").
		Str("order_id", order.OrderInfo.OrderId).
		Msg("Initializing Apollo backend")

	apolloBE := apollo.New(&config.CHAIN_CTX)
	apolloBE = apolloBE.SetWalletFromBech32(order.OrderInfo.MakerAddress.String())

	log.Debug().
		Str("function", "MakerCreateOrder").
		Str("order_id", order.OrderInfo.OrderId).
		Str("maker_address", order.OrderInfo.MakerAddress.String()).
		Msg("Apollo backend initialized and wallet set from Bech32")

	makerCommittingAmount := order.OrderInfo.OrderAmount + order.OrderTxInfo.MakerFee + order.OrderTxInfo.CollateralAmount

	log.Trace().
		Str("function", "MakerCreateOrder").
		Str("order_id", order.OrderInfo.OrderId).
		Int64("order_amount", order.OrderInfo.OrderAmount).
		Int64("maker_fee", order.OrderTxInfo.MakerFee).
		Int64("collateral_amount", order.OrderTxInfo.CollateralAmount).
		Int64("maker_committing_amount", makerCommittingAmount).
		Msg("Calculated maker committing amount")

	log.Debug().
		Str("function", "MakerCreateOrder").
		Str("order_id", order.OrderInfo.OrderId).
		Interface("order", order).
		Int64("order_amount", order.OrderInfo.OrderAmount).
		Int64("maker_committing_amount", makerCommittingAmount).
		Msg("Order details and calculated amounts")

	log.Trace().
		Str("function", "MakerCreateOrder").
		Str("order_id", order.OrderInfo.OrderId).
		Msg("Marshaling order datum to Plutus")

	orderDatumMarshaled, err := plutusEncoder.MarshalPlutus(*order)
	if err != nil {
		log.Error().
			Err(err).
			Str("function", "MakerCreateOrder").
			Str("operation", "marshal_order_datum").
			Str("order_id", order.OrderInfo.OrderId).
			Msg("Failed to marshal order datum to Plutus")
		return "", "", err
	}

	log.Debug().
		Str("function", "MakerCreateOrder").
		Str("order_id", order.OrderInfo.OrderId).
		Msg("Order datum marshaled successfully")

	log.Trace().
		Str("function", "MakerCreateOrder").
		Str("order_id", order.OrderInfo.OrderId).
		Int("user_utxos_count", len(order.OrderTxInfo.UserUtxos)).
		Msg("Retrieving user UTXOs")

	userUtxos, err := utility.GetUserUTxOs(order.OrderTxInfo.UserUtxos)
	if err != nil {
		log.Error().
			Err(err).
			Str("function", "MakerCreateOrder").
			Str("operation", "get_user_utxos").
			Str("order_id", order.OrderInfo.OrderId).
			Int("user_utxos_count", len(order.OrderTxInfo.UserUtxos)).
			Msg("Failed to retrieve user UTXOs")
		return "", "", err
	}

	log.Debug().
		Str("function", "MakerCreateOrder").
		Str("order_id", order.OrderInfo.OrderId).
		Int("user_utxos_retrieved", len(userUtxos)).
		Msg("User UTXOs retrieved successfully")

	log.Trace().
		Str("function", "MakerCreateOrder").
		Str("order_id", order.OrderInfo.OrderId).
		Msg("Retrieving last block slot")

	lastSlot, err := config.CHAIN_CTX.LastBlockSlot()
	if err != nil {
		log.Error().
			Err(err).
			Str("function", "MakerCreateOrder").
			Str("operation", "get_last_block_slot").
			Str("order_id", order.OrderInfo.OrderId).
			Msg("Failed to retrieve last block slot")
		return "", "", err
	}

	log.Debug().
		Str("function", "MakerCreateOrder").
		Str("order_id", order.OrderInfo.OrderId).
		Int("last_slot", lastSlot).
		Int("ttl", lastSlot+600).
		Msg("Last block slot retrieved, TTL calculated")

	log.Trace().
		Str("function", "MakerCreateOrder").
		Str("order_id", order.OrderInfo.OrderId).
		Str("collateral_tx_id", order.OrderTxInfo.CollateralUtxo.TxID).
		Int("collateral_tx_index", order.OrderTxInfo.CollateralUtxo.TxIDIndex).
		Msg("Retrieving collateral UTXO")

	collateralUtxo, err := config.CHAIN_CTX.GetUtxoFromRef(order.OrderTxInfo.CollateralUtxo.TxID, order.OrderTxInfo.CollateralUtxo.TxIDIndex)
	if err != nil {
		log.Error().
			Err(err).
			Str("function", "MakerCreateOrder").
			Str("operation", "get_collateral_utxo").
			Str("order_id", order.OrderInfo.OrderId).
			Str("collateral_tx_id", order.OrderTxInfo.CollateralUtxo.TxID).
			Int("collateral_tx_index", order.OrderTxInfo.CollateralUtxo.TxIDIndex).
			Msg("Failed to retrieve collateral UTXO")
		return "", "", err
	}

	log.Debug().
		Str("function", "MakerCreateOrder").
		Str("order_id", order.OrderInfo.OrderId).
		Str("collateral_tx_id", order.OrderTxInfo.CollateralUtxo.TxID).
		Int("collateral_tx_index", order.OrderTxInfo.CollateralUtxo.TxIDIndex).
		Msg("Collateral UTXO retrieved successfully")

	log.Info().
		Str("function", "MakerCreateOrder").
		Str("order_id", order.OrderInfo.OrderId).
		Msg("Building transaction with Apollo")

	log.Trace().
		Str("function", "MakerCreateOrder").
		Str("order_id", order.OrderInfo.OrderId).
		Str("change_address", order.OrderTxInfo.ChangeAddress.String()).
		Str("escrow_contract_address", order.OrderTxInfo.EscrowContractAddress.String()).
		Int64("maker_committing_amount", makerCommittingAmount).
		Str("state_token_policy_id", order.OrderTxInfo.StateTokenPolicyId).
		Int("ttl", lastSlot+600).
		Msg("Transaction building parameters")

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
		AddReferenceInputV3(
			order.OrderTxInfo.StateTokenRefUtxo.TxID,
			order.OrderTxInfo.StateTokenRefUtxo.TxIDIndex,
		).
		PayToContract(
			order.OrderTxInfo.EscrowContractAddress,
			orderDatumMarshaled,
			int(makerCommittingAmount),
			true,
			apollo.Unit{
				PolicyId: order.OrderTxInfo.StateTokenPolicyId,
				Name:     order.OrderInfo.OrderId,
				Quantity: int(1),
			},
		).
		AddRequiredSigner(adminWallet.AdminPKH).
		AddRequiredSigner(serialization.PubKeyHash(order.OrderInfo.MakerAddress.PaymentPart)).
		SetTtl(int64(lastSlot) + 600).
		Complete()

	if err != nil {
		log.Error().
			Err(err).
			Str("function", "MakerCreateOrder").
			Str("operation", "build_transaction").
			Str("order_id", order.OrderInfo.OrderId).
			Msg("Failed to build transaction")
		return "", "", err
	}

	log.Info().
		Str("function", "MakerCreateOrder").
		Str("order_id", order.OrderInfo.OrderId).
		Msg("Transaction built successfully, signing with admin wallet")

	apolloBE, err = apolloBE.SignWithSkey(adminWallet.AdminVkey, adminWallet.AdminSkey)
	if err != nil {
		log.Error().
			Err(err).
			Str("function", "MakerCreateOrder").
			Str("operation", "sign_transaction").
			Str("order_id", order.OrderInfo.OrderId).
			Msg("Failed to sign transaction with admin wallet")
		return "", "", err
	}

	log.Debug().
		Str("function", "MakerCreateOrder").
		Str("order_id", order.OrderInfo.OrderId).
		Msg("Transaction signed successfully")

	tx := apolloBE.GetTx()
	txHash, err := tx.TransactionBody.Hash()
	if err != nil {
		log.Error().
			Err(err).
			Str("function", "MakerCreateOrder").
			Str("operation", "hash_transaction").
			Str("order_id", order.OrderInfo.OrderId).
			Msg("Failed to hash transaction body")
		return "", "", err
	}

	log.Trace().
		Str("function", "MakerCreateOrder").
		Str("order_id", order.OrderInfo.OrderId).
		Str("tx_hash", hex.EncodeToString(txHash)).
		Msg("Transaction hash calculated")

	txByte, err := tx.Bytes()
	if err != nil {
		log.Error().
			Err(err).
			Str("function", "MakerCreateOrder").
			Str("operation", "serialize_transaction").
			Str("order_id", order.OrderInfo.OrderId).
			Str("tx_hash", hex.EncodeToString(txHash)).
			Msg("Failed to serialize transaction to bytes")
		return "", "", err
	}

	log.Trace().
		Str("function", "MakerCreateOrder").
		Str("order_id", order.OrderInfo.OrderId).
		Str("tx_hash", hex.EncodeToString(txHash)).
		Int("tx_bytes_length", len(txByte)).
		Msg("Transaction serialized to bytes")

	log.Debug().
		Str("function", "MakerCreateOrder").
		Str("order_id", order.OrderInfo.OrderId).
		Str("tx_hash", hex.EncodeToString(txHash)).
		Msg("Evaluating transaction")

	evalTx, err := config.CHAIN_CTX.EvaluateTx(txByte)
	if err != nil {
		log.Error().
			Err(err).
			Str("function", "MakerCreateOrder").
			Str("operation", "evaluate_transaction").
			Str("order_id", order.OrderInfo.OrderId).
			Str("tx_hash", hex.EncodeToString(txHash)).
			Msg("Failed to evaluate transaction")
		return "", "", err
	}

	log.Debug().
		Str("function", "MakerCreateOrder").
		Str("order_id", order.OrderInfo.OrderId).
		Str("tx_hash", hex.EncodeToString(txHash)).
		Int("evaluation_result_length", len(evalTx)).
		Msg("Transaction evaluated successfully")

	cbor, err := Utils.ToCbor(tx)
	if err != nil {
		log.Error().
			Err(err).
			Str("function", "MakerCreateOrder").
			Str("operation", "convert_to_cbor").
			Str("order_id", order.OrderInfo.OrderId).
			Str("tx_hash", hex.EncodeToString(txHash)).
			Msg("Failed to convert transaction to CBOR")
		return "", "", err
	}

	log.Debug().
		Str("function", "MakerCreateOrder").
		Str("order_id", order.OrderInfo.OrderId).
		Str("tx_id", hex.EncodeToString(txHash)).
		Str("tx_cbor", cbor).
		Interface("evaluate_tx", evalTx).
		Msg("Transaction created and converted to CBOR")

	if len(evalTx) == 0 {
		log.Error().
			Str("function", "MakerCreateOrder").
			Str("operation", "validate_evaluation").
			Str("order_id", order.OrderInfo.OrderId).
			Str("tx_hash", hex.EncodeToString(txHash)).
			Msg("Transaction evaluation returned empty result")
		return "", "", fmt.Errorf("transaction evaluation failed")
	}

	log.Info().
		Str("function", "MakerCreateOrder").
		Str("order_id", order.OrderInfo.OrderId).
		Str("tx_id", hex.EncodeToString(txHash)).
		Msg("Maker order created successfully")

	return cbor, hex.EncodeToString(txHash), nil

}
