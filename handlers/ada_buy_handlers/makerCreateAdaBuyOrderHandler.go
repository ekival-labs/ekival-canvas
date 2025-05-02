package ada_buy_handlers

import (
	"encoding/json"
	"io"
	"log"
	"time"

	"github.com/ekival-labs/ekival-canvas/config"
	"github.com/ekival-labs/ekival-canvas/constants"
	"github.com/ekival-labs/ekival-canvas/errors"
	"github.com/ekival-labs/ekival-canvas/model"
	"github.com/ekival-labs/ekival-canvas/txBuilders/ada_p2p_buy"
	"github.com/ekival-labs/ekival-canvas/utility"
	"github.com/ekival-labs/ekival-canvas/viewmodel"

	"github.com/Salvionied/apollo/serialization/Address"
	"github.com/gofiber/fiber/v2"
)

func MakerCreateAdaBuyOrderHandler(c *fiber.Ctx) error {
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

	var mra model.MaybeAddress
	if u.MakerRepAddress != "" {
		addr, err := Address.DecodeAddress(u.MakerRepAddress)
		if err != nil {
			log.Println(err)
			return err
		}
		mra = model.WithAddress{
			Address: addr,
		}
	} else {
		mra = model.Nothing{}
	}

	var tra model.MaybeAddress
	if u.MakerRepAddress != "" {
		addr, err := Address.DecodeAddress(u.TakerRepAddress)
		if err != nil {
			log.Println(err)
			return err
		}
		tra = model.WithAddress{
			Address: addr,
		}
	} else {
		tra = model.Nothing{}
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

	currentTime := time.Now().UnixNano() / int64(time.Millisecond)
	var md int64 = currentTime + (u.MakerDeadline * 1000)
	var td int64 = currentTime + (u.TakerDeadline * 1000)

	makerFee := utility.CalculateFee(u.Precision, u.OrderAmount, u.OrderThreshold, u.MakerPct, u.MakerMinFee)
	collateralAmount := utility.CalculateFee(u.Precision, u.OrderAmount, u.OrderThreshold, u.CollateralPct, u.MinCollateral)

	orderInfo := &model.Order{
		OrderInfo: model.OrderInfo{
			OrderId:         u.OrderId,
			OrderAmount:     u.OrderAmount,
			MakerAddress:    ma,
			MakerRepAddress: mra,
			TakerAddress:    ta,
			TakerRepAddress: tra,
			MakerDeadline:   md,
			TakerDeadline:   td,
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
		TradeState: constants.UNCOMMITTED_ORDER_STATUS,

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
			CollateralAmount: collateralAmount,
			ChangeAddress:    changeAddress,
			UserUtxos:        u.UserUTxOs,
			CollateralUtxo:   u.CollateralUTxO,
		},
	}

	cborString, txHash, err := ada_p2p_buy.MakerCreateOrder(orderInfo, config.GetAPBSTAdminWallet())
	if cborString == "" || txHash == "" || err != nil {
		return errors.TxError(c)
	}
	res := &viewmodel.TxResponse{}
	res.TxCBOR = cborString
	res.TxID = txHash

	c.Status(fiber.StatusOK)
	return enc.Encode(res)

}
