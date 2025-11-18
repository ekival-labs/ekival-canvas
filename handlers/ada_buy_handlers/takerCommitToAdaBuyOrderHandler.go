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

func TakerCommitToAdaBuyOrderHandler(c *fiber.Ctx) error {
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
	//The EOF check if body is empty
	case err == io.EOF:
		c.Status(fiber.StatusBadRequest)
		err := enc.Encode(errors.GeneralError("Nothing Was Sent"))
		if err != nil {
			log.Println(err)
			return errors.ServerErrorHandler(c, err)
		}
		return err
	}

	// Check if Necessary fields are empty or not,
	// if it is empty it'll send an error message
	if u.IsValid() != nil {
		c.Status(fiber.StatusBadRequest)
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

	takerFee := utility.CalculateFee(u.Precision, u.OrderAmount, u.OrderThreshold, u.MakerPct, u.MakerMinFee)
	makerFee := utility.CalculateFee(u.Precision, u.OrderAmount, u.OrderThreshold, u.MakerPct, u.MakerMinFee)
	collateralAmount := utility.CalculateFee(u.Precision, u.OrderAmount, u.OrderThreshold, u.CollateralPct, u.MinCollateral)

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
			ChangeAddress:    changeAddress,
			UserUtxos:        u.UserUTxOs,
			CollateralUtxo:   u.CollateralUTxO,
			OrderUtxo: model.EUTxO{
				TxID:      u.OrderUtxo.TxID,
				TxIDIndex: u.OrderUtxo.TxIDIndex,
			},
		},
	}

	cborString, txHash, err := ada_p2p_buy.TakerCommitToOrder(orderInfo, config.GetAdaP2PBuyAdminWallet())
	if cborString == "" || txHash == "" || err != nil {
		return errors.TxError(c)
	}

	res := &viewmodel.TxResponse{}
	res.TxCBOR = cborString
	res.TxID = txHash

	c.Status(fiber.StatusOK)
	return enc.Encode(res)

}
