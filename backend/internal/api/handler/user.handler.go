package handler

import (
	"github.com/go-playground/validator"
	"github.com/gofiber/fiber/v2"
	"github.com/yagyagoel1/ticketnow/internal/validators"
	errorhandler "github.com/yagyagoel1/ticketnow/pkg/errorHandler"
	"gorm.io/gorm"
)

var validate = validator.New()

type UserHandler struct {
	DB *gorm.DB
}

func (r *UserHandler) SignupUser(c *fiber.Ctx) error {
	request := new(validators.CreateUserReq)
	err := c.BodyParser(request)
	if err != nil {
		return errorhandler.Request(err, c, "There was some problem while parsing the data")

	}
	err = validate.Struct(request)
	if err != nil {
		return errorhandler.Request(err, c, "validation failed")
	}
	err = r.DB.Create(request).Error
	if err != nil {
		return errorhandler.Request(err, c, "there was some problem creating the record")
	}
	c.Status(fiber.StatusOK).JSON(&fiber.Map{
		"success": true,
		"message": "the user signed up successfully",
		"data":    nil,
	})
	return nil
}

func (r *UserHandler) SigninUser(c *fiber.Ctx) error {
	return nil
}
