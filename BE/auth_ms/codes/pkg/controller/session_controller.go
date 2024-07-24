package controller

import (
	"auth_ms/pkg/helper/common"
	"auth_ms/pkg/service"

	"github.com/gofiber/fiber/v2"
)

func Get(ctx *fiber.Ctx) error {
	id, _ := common.ParseRouteParams(ctx, "id")
	pageIntValP, sizeIntValP, offsetIntValP := common.ProcessPageAndSize(ctx)

	sessionService := service.NewSessionService(nil)
	result, _ := sessionService.GetUserSessions(&id, sizeIntValP, offsetIntValP)
	responseP := common.PaginateResponse("", common.CurrentUrl(ctx), *pageIntValP, *sizeIntValP, result, 100)

	return common.HandleResponse(ctx, responseP)
}
