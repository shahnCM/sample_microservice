package route

import (
	"auth_ms/pkg/controller"
	"auth_ms/pkg/middleware"

	"github.com/gofiber/fiber/v2"
)

func InitApiRoutes(app *fiber.App) {
	app.Get("/health-check", controller.ServerAlive)

	v1 := app.Group("auth/api/v1")

	auth := v1.Group("/token")
	auth.Post("/fresh", controller.FreshToken)
	auth.Put("/revoke", controller.RevokeToken)
	auth.Put("/verify", controller.VerifyToken)
	auth.Post("/refresh", controller.RefreshToken)
	auth.Post("/register", controller.RegisterUser)

	session := v1.Group("/sessions/")
	session.Get("/users/:id", middleware.AllowAdmin, controller.Get)
}
