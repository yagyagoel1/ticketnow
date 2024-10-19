package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/yagyagoel1/ticketnow/internal/api/handler"
	"gorm.io/gorm"
)

func SetupUserRoutes(router fiber.Router, db *gorm.DB) {
	UserHandler := handler.UserHandler{DB: db}
	router.Post("/signup", UserHandler.signupUser)
	router.Post("/signin", UserHandler.signinUser)
}