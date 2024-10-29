package handler

import (
	"errors"
	"fmt"
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

	bookings := []models.BookingLock{}
	err := r.DB.Model(&models.BookingLock{}).
		Where("user_id = ?", user.Id).
		Preload("Show", func(db *gorm.DB) *gorm.DB {
			return db.Select("id, name, location, image")
		}).
		Preload("TicketCount.TicketType", func(db *gorm.DB) *gorm.DB {
			return db.Select("id, name, price")
		}).
		Find(&bookings).Error

	if err != nil {
		return errorhandler.Request(nil, c, "there was some problem fetching the data")
	}

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
		fmt.Println("error in parsing the request", err)
		return errorhandler.Request(nil, c, "There was some problem while parsing the data")
	}
	err = validate.Struct(request)
	if err != nil {
		return errorhandler.Request(nil, c, "validation failed")
	}
	var show models.Show
	err = r.DB.Preload("TicketTypes").
		Where("id = ?", request.ShowId).
		First(&show).Error
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
	fmt.Println("request.TicketTypes", request.TicketTypes)
	for ticketTypeIdStr, quantity := range request.TicketTypes {
		fmt.Println("ticketTypeIdStr", ticketTypeIdStr)
		fmt.Println("quantity", quantity)

		quantityInt := quantity
		fmt.Println("quantityInt", quantityInt)
		fmt.Printf("%T/n", quantity)

		ticketTypeId, err := strconv.ParseUint(ticketTypeIdStr, 10, 32)
		if err != nil {
			return errorhandler.Request(nil, c, "invalid ticket type id ")

		}
		var ticketTypeCount struct {
			Count uint   `json:"count"`
			Name  string `json:"name"`
		}
		err = r.DB.Model(&models.TicketType{}).Where("id=?", ticketTypeId).Select("ticket_types.count, ticket_types.name").First(&ticketTypeCount).Error
		if err != nil {
			return errorhandler.Request(nil, c, "invalid ticket type ")
		}

		var lockedTickets []struct {
			Id                  uint `json:"id"`
			TicketCountCategory int  `json:"ticketCountCategory"`
		}

		err = r.DB.Model(&models.BookingLock{}).
			Joins("LEFT JOIN ticket_counts ON ticket_counts.booking_lock_id = booking_locks.id").
			Where("booking_locks.show_id = ? AND booking_locks.lock_time > ? AND ticket_counts.ticket_type_id = ?",
				request.ShowId,
				time.Now(),
				ticketTypeId,
			).
			Select("booking_locks.id, ticket_counts.ticket_count_category as ticket_count_category").
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
			Select("bookings.id, COALESCE(subquery.total_count, 0) as ticket_count_category").
			Joins(`LEFT JOIN (
				SELECT booking_id, SUM(ticket_count_category) as total_count 
				FROM ticket_counts 
				WHERE ticket_type_id = ?
				GROUP BY booking_id
			) as subquery ON subquery.booking_id = bookings.id`, ticketTypeId).
			Where("show_id = ?", request.ShowId).
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
		fmt.Println("QuantityLeft", QuantityLeft)
		fmt.Println("quantityInt", quantityInt)
		fmt.Println("ticketTypeCount.Count", ticketTypeCount.Count)
		fmt.Println("show.tickettypes", show.TicketTypes)
		for _, ticketType := range show.TicketTypes {
			if uint64(ticketType.Id) == ticketTypeId {
				totalCost += float64(quantityInt) * ticketType.Price
				break
			}

		}

	}
	booking := models.BookingLock{
		UserId:    c.Locals("user").(models.User).Id,
		ShowId:    uint(request.ShowId),
		TotalCost: totalCost,
		LockTime:  time.Now().Add(10 * time.Minute),
	}

	tx := r.DB.Begin()
	if tx.Error != nil {
		return errorhandler.Request(tx.Error, c, "failed to start transaction")
	}

	if err := tx.Create(&booking).Error; err != nil {
		tx.Rollback()
		return errorhandler.Request(err, c, "failed to create booking lock")
	}

	for ticketTypeIdStr, quantity := range request.TicketTypes {
		ticketTypeId, err := strconv.ParseUint(ticketTypeIdStr, 10, 32)
		if err != nil {
			tx.Rollback()
			return errorhandler.Request(nil, c, "invalid ticket type id")
		}

		ticketCount := models.TicketCount{
			BookingLockId:       &booking.Id,
			TicketTypeId:        uint(ticketTypeId),
			TicketCountCategory: quantity,
		}

		if err := tx.Omit("BookingId").Create(&ticketCount).Error; err != nil {
			tx.Rollback()
			return errorhandler.Request(err, c, "failed to create ticket count")
		}
	}

	if err := tx.Commit().Error; err != nil {
		return errorhandler.Request(err, c, "failed to commit transaction")
	}
	c.Status(http.StatusCreated).JSON(fiber.Map{
		"success": true,
		"message": "lock booking created successfully",
		"data":    booking,
	})
	return nil

}
