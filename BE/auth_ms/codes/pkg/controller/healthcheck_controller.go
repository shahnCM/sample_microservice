package controller

import (
	"auth_ms/pkg/helper/common"
	"fmt"

	"github.com/gofiber/fiber/v2"
)

func ServerAlive(ctx *fiber.Ctx) error {
	serverAddress := ctx.Context().LocalAddr()
	msg := fmt.Sprintf("Server running on: %s", serverAddress)
	return common.SuccessResponse(ctx, 200, msg, nil, nil)
}
