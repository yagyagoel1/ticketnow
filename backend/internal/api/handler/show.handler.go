package handler

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/yagyagoel1/ticketnow/internal/models"
	"github.com/yagyagoel1/ticketnow/internal/validators"
	errorhandler "github.com/yagyagoel1/ticketnow/pkg/errorHandler"
	"gorm.io/gorm"
)

type ShowHandler struct {
	DB *gorm.DB
}

func (r *ShowHandler) GetAllShows(c *fiber.Ctx) error {
	shows := []models.Show{}
	err := r.DB.Model(&models.Show{}).
		Where("show_timing > ?", time.Now()).
		Select("id, name, description, image, location, show_timing").
		Find(&shows).
		Error
	if err != nil {
		return errorhandler.Request(nil, c, "error in fetching shows")
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"message": "All the shows fetched successfully",
		"data":    shows})

}
func (r *ShowHandler) GetShow(c *fiber.Ctx) error {
	show := models.Show{}

	err := r.DB.Model(&models.Show{}).Where("id=?", c.Params("id")).Preload("TicketTypes", func(db *gorm.DB) *gorm.DB {
		return db.Select("id,Name,Price,Count")
	}).Find(&show).Error
	if err != nil {
		return errorhandler.Request(nil, c, "error in fetching show")
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"message": "Show fetched successfully",
		"data":    show})

}

func (r *ShowHandler) PostShow(c *fiber.Ctx) error {
	if c.Locals("user") == nil {
		return errorhandler.Request(nil, c, "unauthorized")
	}
	user, ok := c.Locals("user").(models.User)
	if !ok {
		return errorhandler.Request(nil, c, "invalid user data")
	}

	if !user.IsAdmin || !user.CreateEvent {
		return errorhandler.Request(nil, c, "unauthorized")
	}

	request := new(validators.CreateShow)

	err := c.BodyParser(request)
	if err != nil {
		return errorhandler.Request(err, c, "error in parsing the request")
	}
	err = validate.Struct(request)
	if err != nil {
		return errorhandler.Request(err, c, "validation failed")
	}
	show := models.Show{
		Name:        request.Name,
		Description: request.Description,
		Image:       request.Image,
		Location:    request.Location,
		ShowTiming:  request.ShowTiming,
		TicketTypes: []models.TicketType{},
		UserId:      user.Id,
	}
	for _, ticketType := range request.TicketTypes {
		show.TicketTypes = append(show.TicketTypes, models.TicketType{
			Name:  ticketType.Name,
			Price: ticketType.Price,
			Count: ticketType.Count,
		})
	}

	err = r.DB.Create(&show).Error
	if err != nil {
		return errorhandler.Request(err, c, "error in creating the show")
	}
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"success": true,
		"message": "Show created successfully",
		"data":    show})
}

func (r *ShowHandler) DeleteShow(c *fiber.Ctx) error {
	if c.Locals("user") == nil {
		return errorhandler.Request(nil, c, "unauthorized")
	}
	user, ok := c.Locals("user").(models.User)
	if !ok {
		return errorhandler.Request(nil, c, "invalid user data")
	}

	if !user.IsAdmin || !user.CreateEvent {
		return errorhandler.Request(nil, c, "unauthorized")
	}
	show := models.Show{}
	err := r.DB.Model(&models.Show{}).Where("id=?", c.Params("id")).First(&show).Error
	if err != nil {
		return errorhandler.Request(err, c, "error in fetching show")
	}
	err = r.DB.Delete(&show).Error
	if err != nil {
		return errorhandler.Request(err, c, "error in deleting show")
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"message": "Show deleted successfully",
		"data":    nil})
}
