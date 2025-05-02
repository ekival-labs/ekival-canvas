package submit_handlers

import (
	"log"

	"github.com/ekival-labs/ekival-canvas/errors"
	"github.com/ekival-labs/ekival-canvas/txBuilders/tx_submits"
	"github.com/ekival-labs/ekival-canvas/viewmodel"

	"github.com/gofiber/fiber/v2"
)

func TxSubmitHandler(c *fiber.Ctx) error {
	var u viewmodel.TxResponse

	c.Set("Content-Type", "application/json")
	c.Set("Access-Control-Allow-Origin", "*")
	c.Set("Access-Control-Allow-Methods", "POST")

	body := c.Body()
	if len(body) == 0 {
		log.Println("Empty request body")
		return c.Status(fiber.StatusBadRequest).JSON(errors.GeneralError("Nothing Was Sent"))
	}

	if err := c.BodyParser(&u); err != nil {
		log.Println("BodyParser error:", err)
		return errors.BadRequestErrorHandler(c, err)
	}

	if err := u.IsValid(); err != nil {
		log.Println("Validation error:", err)
		return c.Status(fiber.StatusBadRequest).JSON(errors.FieldError("txCBOR"))
	}

	response, err := tx_submits.TxSubmitter(&u)
	if err != nil {
		log.Println("TxSubmitter error:", err)
		return errors.TxError(c)
	}

	res := &viewmodel.TxResponse{TxCBOR: response}
	log.Printf("Transaction submitted successfully: %+v\n", res)

	return c.Status(fiber.StatusOK).JSON(res)
}
