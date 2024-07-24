package middleware

import (
	"auth_ms/pkg/action"
	"auth_ms/pkg/enum"
	"auth_ms/pkg/helper/common"
	"strings"

	"github.com/gofiber/fiber/v2"
)

func AllowAdmin(ctx *fiber.Ctx) error {

	jwtToken, _ := common.ParseHeader(ctx, "Authorization", "Bearer ")
	jwtToken = strings.Split(jwtToken, "Bearer ")[1]
	if jwtToken == "" {
		return common.ErrorResponse(ctx, fiber.ErrUnauthorized.Code, fiber.ErrUnauthorized.Message)
	}
	userTokenDataP, err := action.Verify(&jwtToken)
	if err != nil {
		return common.ErrorResponse(ctx, err.Code, err.Message)
	}

	if (*userTokenDataP.(*map[string]any))["role"] != enum.ADMIN {
		return common.ErrorResponse(ctx, fiber.ErrUnauthorized.Code, fiber.ErrUnauthorized.Message)
	}

	return ctx.Next()
}
