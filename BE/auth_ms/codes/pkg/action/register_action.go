package action

import (
	"auth_ms/pkg/dto/request"
	"auth_ms/pkg/helper/common"
	"auth_ms/pkg/service"

	"github.com/gofiber/fiber/v2"
)

func Register(userP *request.UserRegistrationDto) *fiber.Error {

	hashedPassword, err := common.GenerateHash(&userP.Password)
	if err != nil {
		return fiber.ErrInternalServerError
	}
	userP.Password = *hashedPassword

	// Store Username & Password
	userService := service.NewUserService(nil)
	_, err = userService.StoreUser(userP)
	if err != nil {
		// return common.ErrorResponse(ctx, fiber.ErrUnprocessableEntity.Code, strings.Split(err.Error(), ":")[1])
		return fiber.ErrUnprocessableEntity
	}

	return nil
}
