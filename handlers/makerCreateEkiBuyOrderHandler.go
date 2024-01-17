package handlers

import (
	"encoding/json"
	"io"
	"log"
	"time"

	"ekival-canvas/config"
	"ekival-canvas/constants"
	"ekival-canvas/errors"
	"ekival-canvas/txBuilders/token_p2p_buy"
	"ekival-canvas/vars"
	"ekival-canvas/viewmodel"

	"github.com/Salvionied/apollo/serialization/Address"
	"github.com/gofiber/fiber/v2"
)

func MakerCreateEkiBuyOrderHandler(c *fiber.Ctx) error {
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

	ma, err := Address.DecodeAddress(vars.MakerAddress)
	if err != nil {
		log.Println(err)
		return err
	}

	ta, err := Address.DecodeAddress(vars.TakerAddress)
	if err != nil {
		log.Println(err)
		return err
	}

	changeAddress, err := Address.DecodeAddress(u.ChangeAddress)
	if err != nil {
		log.Println(err)
		return err
	}

	escrowContractAddress, err := Address.DecodeAddress(cfg.EkiP2PBuyEscrow.Address)
	if err != nil {
		log.Println(err)
		return err
	}

	currentTime := time.Now().UnixNano() / int64(time.Millisecond)
	var md int64 = currentTime + (vars.MakerDeadline * 1000)
	var td int64 = currentTime + (vars.TakerDeadline * 1000)

	// var order *viewmodel.Order = &viewmodel.Order{}
	orderInfo := &viewmodel.Order{
		OrderInfo: viewmodel.OrderInfo{
			TradeTokenName:     vars.TEKITokenName,
			TradeTokenPolicyId: vars.TEKIPolicyId,
			OrderId:            vars.OrderId,
			OrderAmount:        vars.OrderAmount,
			MakerAddress:       ma,
			TakerAddress:       ta,
			MakerDeadline:      md,
			TakerDeadline:      td,
		},
		BrokerageInfo: viewmodel.BrokerageInfo{
			Precision:      vars.Precision,
			CollateralPct:  vars.CollateralPct,
			MakerPct:       vars.MakerPct,
			TakerPct:       vars.TakerPct,
			CancelPct:      vars.CancelPct,
			MinCollateral:  vars.MinCollateral,
			MakerMinFee:    vars.MakerMinFee,
			TakerMinFee:    vars.TakerMinFee,
			CancelMinFee:   vars.CancelMinFee,
			MinOrderAmount: vars.MinOrderAmount,
			OrderThreshold: vars.OrderThreshold,
			CancelPenalty:  vars.CancelPenalty,
			AdaCollateral:  vars.AdaCollateral,
		},
		TradeState: constants.UNCOMMITTED_ORDER_STATUS,

		OrderTxInfo: viewmodel.OrderTxInfo{
			EscrowContractAddress: escrowContractAddress,
			EscrowContractRefUtxo: viewmodel.EUTxO{
				TxID:      cfg.EkiP2PBuyEscrow.RefTxID,
				TxIDIndex: cfg.EkiP2PBuyEscrow.RefTxIDx,
			},
			StateTokenPolicyId: cfg.EKI_BST.PolicyID,
			StateTokenRefUtxo: viewmodel.EUTxO{
				TxID:      cfg.EKI_BST.RefTxID,
				TxIDIndex: cfg.EKI_BST.RefTxIDx,
			},
			ChangeAddress:  changeAddress,
			UserUtxos:      u.UserUTxOs,
			CollateralUtxo: u.CollateralUTxO,
		},
	}

	cborString, txHash, err := token_p2p_buy.MakerCreateOrder(orderInfo, config.GetEKI_BSTAdminWallet())
	if cborString == "" || txHash == "" || err != nil {
		return errors.TxError(c)
	}
	if err != nil {
		return errors.TxError(c)
	}

	res := &viewmodel.TxResponse{}
	res.TxCBOR = cborString
	res.TxID = txHash

	c.Status(fiber.StatusOK)
	return enc.Encode(res)

}
