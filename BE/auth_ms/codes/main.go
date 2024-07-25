package main

import (
	"auth_ms/pkg/config/fiberconfig"
	"auth_ms/pkg/helper/safeasync"
	"auth_ms/pkg/migration"
	"auth_ms/pkg/provider/database/mariadb10"
	"auth_ms/pkg/route"
	"log"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/joho/godotenv"
)

func main() {
	log.Println("SERVER INITIATING")

	if err := godotenv.Load(".env"); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	} else {
		log.Println("Loaded .env file")
	}

	timeZone := os.Getenv("APP_TIMEZONE")
	localTimeZone, err := time.LoadLocation(timeZone)
	if err != nil {
		log.Fatalf("Must set APP_TIMEZONE in .env: %v", err)
	}
	time.Local = localTimeZone

	// if err := queue.Init(); err != nil {
	// 	log.Printf("Failed to initialize RabbitMQ queue: %v", err)
	// } else {
	// 	log.Println("Connected to RabbitMQ")
	// }

	if err := mariadb10.ConnectToMariaDb10(); err != nil {
		log.Printf("Failed to initialize Database: %v", err)
	} else {
		safeasync.Run(migration.RunMigration)
		log.Println("Connected to Database")
	}

	app := fiber.New(fiberconfig.FiberConfig())
	app.Use(logger.New(logger.Config{Format: "[${ip}]:${port} ${status} - ${method} ${path}\n"}))
	app.Use(recover.New(fiberconfig.RecoveryConfig()))

	route.InitApiRoutes(app)

	port := os.Getenv("SERVER_PORT")
	if port == "" {
		port = "8000"
	}

	log.Println(app.Listen(":" + port))
}
