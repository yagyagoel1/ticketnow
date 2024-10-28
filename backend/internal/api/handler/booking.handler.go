package handler

import (
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/yagyagoel1/ticketnow/internal/models"
	"github.com/yagyagoel1/ticketnow/internal/validators"
	errorhandler "github.com/yagyagoel1/ticketnow/pkg/errorHandler"
	"gorm.io/gorm"
)

type BookingHandler struct {
	DB *gorm.DB
}

func (r *BookingHandler) GetBookings(c *fiber.Ctx) error {
	user, ok := c.Locals("user").(models.User)
	if !ok {
		return errorhandler.Request(nil, c, "unauthorized")
	}
	bookings := []models.Booking{}
	err := r.DB.Model(&models.Booking{}).Where("userId=?", user.Id).Preload("Show", func(db *gorm.DB) *gorm.DB {
		return db.Select("id,name,location,image")
	}).Preload("TicketCount", func(db *gorm.DB) *gorm.DB {
		return db.Select("id, ticketCountCategory")
	}).Preload("TicketCount.TicketType", func(db *gorm.DB) *gorm.DB {
		return db.Select("name,price")
	}).Find(&bookings).Error
	if err != nil {
		return errorhandler.Request(nil, c, "there was some problem fetching the data")
	}
	// var responseData []map[string]interface{}

	// for _, booking := range bookings {
	// 	responseData = append(responseData, map[string]interface{}{
	// 		"id":        booking.Id,
	// 		"totalCost": booking.TotalCost,
	// 		"show": map[string]interface{}{
	// 			"id":       booking.Show.Id,
	// 			"name":     booking.Show.Name,
	// 			"location": booking.Show.Location,
	// 			"image":    booking.Show.Image,
	// 		},
	// 		"ticketType": map[string]interface{}{
	// 			"id":        booking.TicketType.Id,
	// 			"type_name": booking.TicketType.Name,
	// 			"price":     booking.TicketType.Price,
	// 		},
	// 	})
	// }

	// Send the filtered data in the response
	c.Status(http.StatusOK).JSON(fiber.Map{
		"success": true,
		"message": "All the booking details fetched successfully",
		"data":    bookings,
	})
	return nil
}

func (r *BookingHandler) PostBooking(c *fiber.Ctx) error {

	request := new(validators.PostBooking)
	err := c.BodyParser(request)
	if err != nil {
		return errorhandler.Request(nil, c, "There was some problem while parsing the data")
	}
	err = validate.Struct(request)
	if err != nil {
		return errorhandler.Request(nil, c, "validation failed")
	}
	var show models.Show
	err = r.DB.Model(&models.Show{}).Where("id=?", request.ShowId).Preload("ticketTypes", func(db *gorm.DB) *gorm.DB {
		return db.Select("id,name,count,price")
	}).First(&show).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return errorhandler.Request(nil, c, "there was no show with the given id")
	}
	if err != nil {
		return errorhandler.Request(nil, c, "there was some error while fetching the show with the given id ")
	}
	if show.ShowTiming.Before(time.Now()) {
		return errorhandler.Request(nil, c, "the show is already over")
	}
	totalCost := 0.0
	for ticketTypeIdStr, quantity := range request.TicketTypes {
		quantityInt, ok := quantity.(int)
		if !ok {
			return errorhandler.Request(nil, c, "invalid ticket quantity specified")

		}
		ticketTypeId, err := strconv.ParseUint(ticketTypeIdStr, 10, 32)
		if err != nil {
			return errorhandler.Request(nil, c, "invalid ticket type id ")

		}
		var ticketTypeCount struct {
			Count uint   `json:"count"`
			Name  string `json:"name"`
		}
		err = r.DB.Model(&models.TicketType{}).Where("id=?", ticketTypeId).Select("TicketType.count ,TicketType.name").First(&ticketTypeCount).Error
		if err != nil {
			return errorhandler.Request(nil, c, "invalid ticket type ")
		}

		var lockedTickets []struct {
			Id                  uint `json:"id"`
			TicketCountCategory int  `json:"ticketCountCategory"`
		}

		err = r.DB.Model(&models.BookingLock{}).
			Where("showId = ? AND LockTime > ?", request.ShowId, time.Now()).
			Preload("TicketCount", func(db *gorm.DB) *gorm.DB {
				return db.Select("ticketTypeId, SUM(ticketCountCategory) as ticketCountCategory").
					Where("ticketTypeId = ?", ticketTypeId).
					Group("ticketTypeId")
			}).
			Select("BookingLock.id, TicketCount.ticketCountCategory").
			Find(&lockedTickets).Error
		if err != nil {
			return errorhandler.Request(err, c, "error while fetching the locked tickets")
		}
		var totalSeatsBooked int
		for _, record := range lockedTickets {
			totalSeatsBooked += record.TicketCountCategory
		}
		var bookingTickets []struct {
			Id                  uint `json:"id"`
			TicketCountCategory int  `json:"ticketCountCategory"`
		}

		err = r.DB.Model(&models.Booking{}).
			Where("showId = ?", request.ShowId).
			Preload("TicketCount", func(db *gorm.DB) *gorm.DB {
				return db.Select("ticketTypeId, SUM(ticketCountCategory) as ticketCountCategory").
					Where("ticketTypeId = ?", ticketTypeId).
					Group("ticketTypeId")
			}).
			Select("Booking.id, TicketCount.ticketCountCategory").
			Find(&bookingTickets).Error
		if err != nil {
			return errorhandler.Request(err, c, "error while fetching the booking tickets")
		}
		for _, record := range bookingTickets {
			totalSeatsBooked += record.TicketCountCategory
		}
		QuantityLeft := int(ticketTypeCount.Count) - totalSeatsBooked
		if QuantityLeft < quantityInt {
			return errorhandler.Request(nil, c, "the quantity of the ticket is not available")
		}
		totalCost += float64(quantityInt) * show.TicketTypes[ticketTypeId].Price

	}
	booking := models.BookingLock{
		UserId:      c.Locals("user").(models.User).Id,
		ShowId:      uint(request.ShowId),
		TicketCount: []models.TicketCount{},
		TotalCost:   totalCost,
		LockTime:    time.Now().Add(10 * time.Minute),
	}
	for ticketTypeIdStr, quantity := range request.TicketTypes {
		quantityInt, ok := quantity.(int)
		if !ok {
			return errorhandler.Request(nil, c, "invalid ticket quantity specified")

		}
		ticketTypeId, err := strconv.ParseUint(ticketTypeIdStr, 10, 32)
		if err != nil {
			return errorhandler.Request(nil, c, "invalid ticket type id ")

		}
		booking.TicketCount = append(booking.TicketCount, models.TicketCount{
			TicketTypeId:        uint(ticketTypeId),
			TicketCountCategory: quantityInt,
		})
	}
	err = r.DB.Create(&booking).Error
	if err != nil {
		return errorhandler.Request(err, c, "failed to create booking")
	}
	c.Status(http.StatusCreated).JSON(fiber.Map{
		"success": true,
		"message": "lock booking created successfully",
		"data":    nil,
	})
	return nil

}
