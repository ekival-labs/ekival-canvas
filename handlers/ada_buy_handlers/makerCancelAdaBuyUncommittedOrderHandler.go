package ada_buy_handlers

import (
	"encoding/json"
	"io"
	"log"

	"ekival-canvas/config"
	"ekival-canvas/errors"
	"ekival-canvas/model"
	"ekival-canvas/txBuilders/ada_p2p_buy"
	"ekival-canvas/utility"
	"ekival-canvas/viewmodel"

	"github.com/Salvionied/apollo/serialization/Address"
	"github.com/gofiber/fiber/v2"
)

// MakerCancelAdaBuyUncommittedOrderHandler cancels an uncommitted ADA buy order by the maker.
func MakerCancelAdaBuyUncommittedOrderHandler(c *fiber.Ctx) error {
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

	changeAddress, err := Address.DecodeAddress(u.ChangeAddress)
	if err != nil {
		log.Println(err)
		return err
	}

	makerFee := utility.CalculateFee(u.Precision, u.OrderAmount, u.OrderThreshold, u.MakerPct, u.MakerMinFee)
	collateralAmount := utility.CalculateFee(u.Precision, u.OrderAmount, u.OrderThreshold, u.CollateralPct, u.MinCollateral)

	orderInfo := &model.Order{
		OrderInfo: model.OrderInfo{
			OrderId:       u.OrderId,
			OrderAmount:   u.OrderAmount,
			MakerAddress:  ma,
			MakerDeadline: u.MakerDeadline,
			TakerDeadline: u.TakerDeadline,
		},
		OrderTxInfo: model.OrderTxInfo{
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
			OrderUtxo: model.EUTxO{
				TxID:      u.OrderUtxo.TxID,
				TxIDIndex: u.OrderUtxo.TxIDIndex,
			},
		},
	}

	cborString, txHash, err := ada_p2p_buy.MakerCancelUncommittedOrder(orderInfo, config.GetAdaP2PBuyAdminWallet())
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
