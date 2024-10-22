package middleware

import (
	"fmt"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt"
	"github.com/yagyagoel1/ticketnow/internal/models"
	errorhandler "github.com/yagyagoel1/ticketnow/pkg/errorHandler"
	"gorm.io/gorm"
)

func Auth(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		tokenString := c.Cookies(os.Getenv("AUTH_TOKEN_NAME"))
		if tokenString == "" {
			return errorhandler.Request(nil, c, "error while fetching token please login again.")
		}

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(os.Getenv("JWT_SECRET")), nil
		})

		if err != nil {
			return errorhandler.Request(nil, c, "error while authenticating")
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			if float64(time.Now().Unix()) > claims["exp"].(float64) {
				return errorhandler.Request(nil, c, "token is expired, please login again")
			}

			var user models.User
			if err := db.First(&user, claims["sub"]).Error; err != nil {
				return errorhandler.Request(nil, c, "user not found")
			}

			c.Locals("user", user)
			return c.Next()
		}

		return errorhandler.Request(nil, c, "invalid token")
	}
}
