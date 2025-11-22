package ada_buy_handlers

import (
	"encoding/json"
	"io"
	"log"

	"ekival-canvas/config"
	"ekival-canvas/constants"
	"ekival-canvas/errors"
	"ekival-canvas/model"
	"ekival-canvas/txBuilders/ada_p2p_buy"
	"ekival-canvas/utility"
	"ekival-canvas/viewmodel"

	"github.com/Salvionied/apollo/serialization/Address"
	"github.com/gofiber/fiber/v2"
)

// TakerWinAdaBuyDisputeHandler handles the case when taker wins a dispute on ADA buy order.
func TakerWinAdaBuyDisputeHandler(c *fiber.Ctx) error {
	var u *viewmodel.UserTxInfo

	enc := json.NewEncoder(c.Response().BodyWriter())
	enc.SetIndent("", "    ")

	c.Response().Header.Set("Content-Type", "application/json")
	c.Response().Header.Set("Access-Control-Allow-Origin", "*")
	c.Response().Header.Set("Access-Control-Allow-Methods", "POST")

	err := c.BodyParser(&u)
	if err != nil {
		log.Println(err)
		return errors.BadRequestErrorHandler(c, err)
	}

	switch {
	case err == io.EOF:
		err := enc.Encode(errors.GeneralError("Nothing Was Sent"))
		if err != nil {
			log.Println(err)
			return errors.ServerErrorHandler(c, err)
		}
		return err
	}

	if u.IsValid() != nil {
		err := enc.Encode(errors.FieldError("Address", "ChangeAddress", "UserUTxOs"))
		if err != nil {
			log.Println(err)
			return errors.ServerErrorHandler(c, err)
		}
		return err
	}

	cfg := config.GetGlobalConfig()

	ma, err := Address.DecodeAddress(u.MakerAddress)
	if err != nil {
		log.Println(err)
		return err
	}

	ta, err := Address.DecodeAddress(u.TakerAddress)
	if err != nil {
		log.Println(err)
		return err
	}

	changeAddress, err := Address.DecodeAddress(u.ChangeAddress)
	if err != nil {
		log.Println(err)
		return err
	}

	escrowContractAddress, err := Address.DecodeAddress(cfg.ADAMarketplace.AdaP2PBuyEscrow.Address)
	if err != nil {
		log.Println(err)
		return err
	}

	treasuryAddress, err := Address.DecodeAddress(cfg.MainTreasury.Address)
	if err != nil {
		log.Println(err)
		return err
	}

	treasuryInfo := &model.TreasuryInfo{
		Address: treasuryAddress,
		Datum:   utility.CreateSimpleDatum(constants.INDEX_ONE, constants.ADA_P2P_BUY_FEE_TYPE),
	}

	makerFee, err := utility.CalculateFee(u.Precision, u.OrderThreshold, u.MinOrderAmount, u.MakerMinFee, u.MakerPct, u.OrderAmount)
	if err != nil {
		log.Println(err)
		return err
	}
	takerFee, err := utility.CalculateFee(u.Precision, u.OrderThreshold, u.MinOrderAmount, u.TakerMinFee, u.TakerPct, u.OrderAmount)
	if err != nil {
		log.Println(err)
		return err
	}
	collateralAmount, err := utility.CalculateFee(u.Precision, u.OrderThreshold, u.MinOrderAmount, u.MinCollateral, u.CollateralPct, u.OrderAmount)
	if err != nil {
		log.Println(err)
		return err
	}
	// Dispute fee configuration currently comes from backend/on-chain state; here we keep it zero.
	var disputeFee int64 = 0

	orderInfo := &model.Order{
		OrderInfo: model.OrderInfo{
			OrderId:       u.OrderId,
			OrderAmount:   u.OrderAmount,
			MakerAddress:  ma,
			TakerAddress:  ta,
			MakerDeadline: u.MakerDeadline,
			TakerDeadline: u.TakerDeadline,
		},
		BrokerageInfo: model.BrokerageInfo{
			Precision:      u.Precision,
			CollateralPct:  u.CollateralPct,
			MakerPct:       u.MakerPct,
			TakerPct:       u.TakerPct,
			CancelPct:      u.CancelPct,
			MinCollateral:  u.MinCollateral,
			MakerMinFee:    u.MakerMinFee,
			TakerMinFee:    u.TakerMinFee,
			CancelMinFee:   u.CancelMinFee,
			MinOrderAmount: u.MinOrderAmount,
			OrderThreshold: u.OrderThreshold,
			CancelPenalty:  u.CancelPenalty,
		},
		// Dispute resolution happens only for active/committed orders.
		TradeState: constants.COMMITTED_ORDER_STATUS,
		OrderTxInfo: model.OrderTxInfo{
			EscrowContractAddress: escrowContractAddress,
			EscrowContractRefUtxo: model.EUTxO{
				TxID:      cfg.ADAMarketplace.AdaP2PBuyEscrow.RefTxID,
				TxIDIndex: cfg.ADAMarketplace.AdaP2PBuyEscrow.RefTxIDx,
			},
			StateTokenPolicyId: cfg.ADAMarketplace.APBST.PolicyID,
			StateTokenRefUtxo: model.EUTxO{
				TxID:      cfg.ADAMarketplace.APBST.RefTxID,
				TxIDIndex: cfg.ADAMarketplace.APBST.RefTxIDx,
			},
			MakerFee:         makerFee,
			TakerFee:         takerFee,
			CollateralAmount: collateralAmount,
			DisputeFee:       disputeFee,
			ChangeAddress:    changeAddress,
			UserUtxos:        u.UserUTxOs,
			CollateralUtxo:   u.CollateralUTxO,
			OrderUtxo: model.EUTxO{
				TxID:      u.OrderUtxo.TxID,
				TxIDIndex: u.OrderUtxo.TxIDIndex,
			},
		},
	}

	// Dispute parameters are taken from on-chain datum / backend state;
	// here we pass an empty redeemer shell and let the validator enforce details.
	takerWinDisputeRedeemer := &model.MakerWinDisputeRedeemer{}

	cborString, txHash, err := ada_p2p_buy.TakerWinDispute(orderInfo, treasuryInfo, takerWinDisputeRedeemer, config.GetAdaP2PBuyAdminWallet())
	if cborString == "" || txHash == "" || err != nil {
		return errors.TxError(c)
	}

	res := &viewmodel.TxResponse{
		TxCBOR: cborString,
		TxID:   txHash,
	}

	c.Status(fiber.StatusOK)
	return enc.Encode(res)
}
