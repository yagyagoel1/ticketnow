package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/yagyagoel1/ticketnow/internal/api/handler"
	"github.com/yagyagoel1/ticketnow/internal/api/middleware"
	"gorm.io/gorm"
)

func SetupBookingRoutes(router fiber.Router, db *gorm.DB) {
	BookingHandler := handler.BookingHandler{DB: db}
	router.Get("/bookings", middleware.Auth(db), BookingHandler.GetBookings)
	router.Post("/bookings", middleware.Auth(db), BookingHandler.PostBooking)
}
