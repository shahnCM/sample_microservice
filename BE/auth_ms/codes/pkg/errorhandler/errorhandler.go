package errorhandler

import (
	"auth_ms/pkg/helper/common"
	"auth_ms/pkg/provider/database/mariadb10"
	"errors"
	"log"

	"github.com/gofiber/fiber/v2"
)

func CustomFiberErrorHandler(ctx *fiber.Ctx, err error) error {

	mariadb10.TransactionRollback()
	log.Println("CustomeFiberErrorHandler recovered from panic: ")

	// Status code defaults to 500
	code := 500
	message := err.Error() //"Internal Server Error"

	// Retrieve the custom status code if it's a *fiber.Error
	var e *fiber.Error
	if errors.As(err, &e) {
		code = e.Code
		message = e.Message
	}

	return common.ErrorResponse(ctx, code, message)
}
