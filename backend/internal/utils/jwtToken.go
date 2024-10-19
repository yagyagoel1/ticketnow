package utils

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/yagyagoel1/ticketnow/internal/models"
)

func GenerateToken(user models.User) (string, error) {
	claims := jwt.MapClaims{
		"id":    user.Id,
		"email": user.Email,
		"exp":   time.Now().Add(time.Hour * 24).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(os.Getenv("JWT_SECRET")))
}
