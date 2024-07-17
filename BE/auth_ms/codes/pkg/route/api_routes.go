package route

import (
	"auth_ms/pkg/controller"

	"github.com/gofiber/fiber/v2"
)

func InitApiRoutes(app *fiber.App) {
	app.Get("/health-check", controller.ServerAlive)

	v1 := app.Group("/v1")
	auth := v1.Group("/auth/token")

	auth.Post("/fresh", controller.FreshToken)
	auth.Put("/revoke", controller.RevokeToken)
	auth.Put("/verify", controller.VerifyToken)
	auth.Post("/refresh", controller.RefreshToken)
	auth.Post("/register", controller.RegisterUser)
}
