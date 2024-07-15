package route

import (
	"auth_ms/pkg/controller"

	"github.com/gofiber/fiber/v2"
)

func InitApiRoutes(app *fiber.App) {

	app.Get("/health-check", controller.ServerAlive)

	/**
	* @group API V1
	* @basepath /v1
	 */
	v1 := app.Group("/v1")

	auth := v1.Group("/auth/token")

	auth.Post("/login", controller.Login)
	auth.Post("/register", controller.Register)
	auth.Put("/verify", controller.Verify)
	auth.Get("/revoke", controller.Revoke)
	auth.Get("/refresh", controller.Refresh)
}
