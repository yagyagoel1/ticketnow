package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/yagyagoel1/ticketnow/internal/api/handler"
	"github.com/yagyagoel1/ticketnow/internal/api/middleware"
	"gorm.io/gorm"
)

func SetupShowRoutes(router fiber.Router, db *gorm.DB) {
	ShowHandler := handler.ShowHandler{DB: db}
	router.Get("/shows", middleware.Auth(db), ShowHandler.GetAllShows)
	router.Get("/show/:id", middleware.Auth(db), ShowHandler.GetShow)
	router.Post("/show", middleware.Auth(db), ShowHandler.PostShow)
	router.Delete("/show/:id", middleware.Auth(db), ShowHandler.DeleteShow)
}
