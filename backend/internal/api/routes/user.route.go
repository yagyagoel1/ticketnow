package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/yagyagoel1/ticketnow/internal/api/handler"
	"github.com/yagyagoel1/ticketnow/internal/api/middleware"
	"gorm.io/gorm"
)

func SetupUserRoutes(router fiber.Router, db *gorm.DB) {
	UserHandler := handler.UserHandler{DB: db}
	router.Post("/signup", UserHandler.SignupUser)
	router.Post("/signin", UserHandler.SigninUser)
	router.Get("/profile", middleware.Auth(db), UserHandler.GetProfile)
	router.Put("/profile", middleware.Auth(db), UserHandler.PutProfile)
	router.Put("/profile/password", middleware.Auth(db), UserHandler.PutPassword)
}
