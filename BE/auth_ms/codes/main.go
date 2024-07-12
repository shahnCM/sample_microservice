package main

import (
	"auth_ms/pkg/config"
	"auth_ms/pkg/queue"
	"auth_ms/pkg/route"
	"fmt"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/joho/godotenv"
)

func main() {
	log.Println("SERVER INITIATING")
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	if err := queue.Init(); err != nil {
		log.Println(fmt.Sprintf("Failed to initialize RabbitMQ queue: %v", err))
	} else {
		log.Println("Connected to RabbitMQ")
	}

	app := fiber.New(config.FiberConfig())
	app.Use(logger.New(logger.Config{
		Format: "[${ip}]:${port} ${status} - ${method} ${path}\n",
	}))
	app.Use(recover.New(config.RecoveryConfig()))
	route.InitApiRoutes(app)

	port := os.Getenv("SERVER_PORT")
	if port == "" {
		port = "8000"
	}

	log.Println(app.Listen(":" + port))
}
