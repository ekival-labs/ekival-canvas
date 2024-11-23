package offers_handlers

import (
	"encoding/json"
	"io"
	"log"

	"ekival-canvas/config"
	"ekival-canvas/constants"
	"ekival-canvas/errors"
	"ekival-canvas/model"
	"ekival-canvas/txBuilders/offers_tx_builders"
	"ekival-canvas/utility"
	"ekival-canvas/viewmodel"

	"github.com/Salvionied/apollo/serialization/Address"
	"github.com/gofiber/fiber/v2"
)

func MakerCreateOfferHandler(c *fiber.Ctx) error {
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

	ma, err := Address.DecodeAddress(u.Address)
	if err != nil {
		log.Println(err)
		return errors.ServerErrorHandler(c, err)
	}

	changeAddress, err := Address.DecodeAddress(u.ChangeAddress)
	if err != nil {
		log.Println(err)
		return errors.ServerErrorHandler(c, err)
	}

	orderInfo := &model.AdaOfferTxInfo{
		MakerAddress:      ma,
		ChangeAddress:     changeAddress,
		UserUtxos:         u.UserUTxOs,
		CollateralUtxo:    u.CollateralUTxO,
		EkivalFeeLovelace: constants.EKIVAL_FEE,
	}

	treasuryAddress, err := Address.DecodeAddress(cfg.MainTreasury.Address)
	if err != nil {
		log.Println(err)
		return errors.ServerErrorHandler(c, err)
	}
	treasuryInfo := &model.TreasuryInfo{
		Address: treasuryAddress,
		Datum:   utility.CreateSimpleDatum(constants.INDEX_ONE, constants.ADA_P2P_SELL_FEE_TYPE),
	}

	cborString, txHash, err := offers_tx_builders.MakerCreateAdaOffer(orderInfo, treasuryInfo, config.GetOfferAdminWallet())
	if cborString == "" || txHash == "" || err != nil {
		return errors.TxError(c)
	}


	res := &viewmodel.TxResponse{}
	res.TxCBOR = cborString
	res.TxID = txHash

	c.Status(fiber.StatusOK)
	return enc.Encode(res)

}
