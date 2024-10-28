package main

import (
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"github.com/yagyagoel1/ticketnow/internal/api/routes"
	"github.com/yagyagoel1/ticketnow/internal/models"
	errorhandler "github.com/yagyagoel1/ticketnow/pkg/errorHandler"
	"github.com/yagyagoel1/ticketnow/pkg/storage"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal(err)

	}
	config := &storage.Config{
		Host:     os.Getenv("DB_HOST"),
		Port:     os.Getenv("DB_PORT"),
		Password: os.Getenv("DB_PASSWORD"),
		User:     os.Getenv("DB_USER"),
		SSLMode:  os.Getenv("DB_SSLMODE"),
		DBName:   os.Getenv("DB_NAME"),
	}
	db, err := storage.NewConnection(config)
	models.AutoMigrate(db)
	errorhandler.Fatal(err)
	app := fiber.New()
	routes.SetupUserRoutes(app.Group("/api/user"), db)
	routes.SetupBookingRoutes(app.Group("/api/booking"), db)
	routes.SetupShowRoutes(app.Group("/api/show"), db)
	log.Fatal(app.Listen(os.Getenv("PORT")))
}
