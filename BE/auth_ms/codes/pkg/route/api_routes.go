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

	auth.Post("/fresh", controller.Login)
	auth.Put("/revoke", controller.Logout)
	auth.Post("/register", controller.Register)
	auth.Post("/refresh", controller.Refresh)
	auth.Put("/verify", controller.Verify)
}
