package handler

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/go-playground/validator"
	"github.com/gofiber/fiber/v2"
	"github.com/yagyagoel1/ticketnow/internal/models"
	"github.com/yagyagoel1/ticketnow/internal/utils"
	"github.com/yagyagoel1/ticketnow/internal/validators"
	errorhandler "github.com/yagyagoel1/ticketnow/pkg/errorHandler"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

var validate = validator.New()

type UserHandler struct {
	DB *gorm.DB
}

// todo add email based auth
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
	var count int64
	err = r.DB.Model(&models.User{}).Where("email = ?", request.Email).Count(&count).Error

	if err != nil {
		return errorhandler.Request(nil, c, "there was some problem creating the record")
	}
	if count > 0 {
		return errorhandler.Request(nil, c, "the user exist with this email already")
	}
	user := &models.User{
		Name:     request.Name,
		Email:    request.Email,
		Password: request.Password,
	}
	err = r.DB.Create(user).Error
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
	request := new(validators.SigninUser)

	err := c.BodyParser(request)
	if err != nil {
		return errorhandler.Request(err, c, "There was some problem while parsing the data")

	}
	err = validate.Struct(request)
	if err != nil {
		return errorhandler.Request(err, c, "validation failed")
	}
	user := models.User{}
	err = r.DB.Model(&models.User{}).Where("email=?", request.Email).First(&user).Error
	if err != nil {
		return errorhandler.Request(err, c, fmt.Sprintf("cannot find the user with the email %s", request.Email))

	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(request.Password))
	if err != nil {
		return errorhandler.Request(err, c, "Password doesnt matched")
	}
	token, err := utils.GenerateToken(user)
	if err != nil {
		return errorhandler.Request(err, c, "there was some problem generating the token")
	}
	cookie := new(fiber.Cookie)
	cookie.Name = os.Getenv("AUTH_TOKEN_NAME")
	cookie.Value = token
	cookie.Expires = time.Now().Add(24 * time.Hour)
	c.Cookie(cookie)
	err = r.DB.Model(&models.User{}).Where("email=?", request.Email).Update("token", token).Error
	if err != nil {
		return errorhandler.Request(nil, c, "error updating token to db")
	}
	c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"message": "Login successful",
		"user":    nil,
	})
	return nil
}

func (r *UserHandler) GetProfile(c *fiber.Ctx) error {
	user, ok := c.Locals("user").(models.User)
	if !ok {
		return errorhandler.Request(nil, c, "unauthorized")
	}
	c.Status(http.StatusOK).JSON(&fiber.Map{

		"success": true,
		"message": "user retreived Successfully",
		"user": fiber.Map{
			"id":    user.Id,
			"name":  user.Name,
			"email": user.Email,
		}})
	return nil
}

func (r *UserHandler) PutProfile(c *fiber.Ctx) error {
	user, ok := c.Locals("user").(models.User)
	if !ok {
		return errorhandler.Request(nil, c, "unauthorized")
	}
	request := new(validators.UpdateProfile)
	err := c.BodyParser(request)
	if err != nil {
		return errorhandler.Request(nil, c, "There was some problem while parsing the data")

	}
	err = validate.Struct(request)
	if err != nil {
		return errorhandler.Request(nil, c, "validation failed")
	}
	err = r.DB.Model(&models.User{}).Where("id=?", user.Id).Updates(map[string]interface{}{
		"name": request.Name,
	}).Error
	if err != nil {
		return errorhandler.Request(nil, c, "There was some problem while updating the record")
	}
	c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"message": "the name of the user updated",
		"data":    nil,
	})
	return nil
}
func (r *UserHandler) PutPassword(c *fiber.Ctx) error {
	user, ok := c.Locals("user").(models.User)
	if !ok {
		return errorhandler.Request(nil, c, "unauthorized")
	}
	request := new(validators.UpdatePassword)
	err := c.BodyParser(request)
	if err != nil {
		return errorhandler.Request(nil, c, "There was some problem while parsing the data")

	}
	err = validate.Struct(request)
	if err != nil {
		return errorhandler.Request(nil, c, "validation failed")
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(request.OldPassword))
	if err != nil {
		return errorhandler.Request(err, c, "Password doesnt matched")
	}
	err = r.DB.Model(&models.User{}).Where("id=?", user.Id).Updates(map[string]interface{}{
		"password": request.NewPassword,
	}).Error
	if err != nil {
		return errorhandler.Request(nil, c, "There was some problem while updating the record")
	}
	c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"message": "the password of the user updated",
		"data":    nil,
	})
	return nil
}
