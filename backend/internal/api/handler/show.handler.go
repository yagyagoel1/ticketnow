package handler

import (
	"encoding/json"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/yagyagoel1/ticketnow/internal/models"
	"github.com/yagyagoel1/ticketnow/internal/utils"
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

	err := r.DB.Preload("TicketTypes").First(&show, c.Params("id")).Error
	if err != nil {
		return errorhandler.Request(err, c, "error in fetching show")
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

	// if !user.IsAdmin || !user.CreateEvent {
	// 	return errorhandler.Request(nil, c, "unauthorized")
	// }

	form, err := c.MultipartForm()
	if err != nil {
		return errorhandler.Request(err, c, "error in parsing the form data")
	}
	request := new(validators.CreateShow)
	request.Name = form.Value["name"][0]
	request.Description = form.Value["description"][0]
	request.Location = form.Value["location"][0]
	request.ShowTiming, _ = time.Parse(time.RFC3339, form.Value["showTiming"][0])
	jsonStr := form.Value["ticketTypes"][0]
	var ticketTypes []validators.TicketType
	err = json.Unmarshal([]byte(jsonStr), &ticketTypes)
	if err != nil {
		return errorhandler.Request(nil, c, "error in parsing the ticket types")
	}
	request.TicketTypes = ticketTypes

	err = validate.Struct(request)
	if err != nil {
		return errorhandler.Request(err, c, "validation failed")
	}
	file, err := c.FormFile("image")
	if err != nil {
		return errorhandler.Request(err, c, "error in uploading image")
	}
	show := models.Show{
		Name:        request.Name,
		Description: request.Description,
		Location:    request.Location,
		ShowTiming:  request.ShowTiming,
		UserId:      user.Id,
		TicketTypes: make([]models.TicketType, 0, len(request.TicketTypes)),
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
	showId := show.Id
	fileUrl, err := utils.UploadToS3(file, strconv.FormatUint(uint64(showId), 10))
	if err != nil {
		return errorhandler.Request(err, c, "error in uploading image")
	}
	err = r.DB.Model(&models.Show{}).Where("id=?", showId).Update("image", fileUrl).Error
	if err != nil {
		return errorhandler.Request(err, c, "error in updating image url")

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

	showID := c.Params("id")

	tx := r.DB.Begin()

	show := models.Show{}
	if err := tx.First(&show, showID).Error; err != nil {
		tx.Rollback()
		return errorhandler.Request(err, c, "error in fetching show")
	}

	if err := tx.Where("show_id = ?", showID).Delete(&models.TicketType{}).Error; err != nil {
		tx.Rollback()
		return errorhandler.Request(err, c, "error deleting ticket types")
	}

	if err := tx.Delete(&show).Error; err != nil {
		tx.Rollback()
		return errorhandler.Request(err, c, "error deleting show")
	}

	tx.Commit()

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"message": "Show deleted successfully",
		"data":    nil,
	})
}
