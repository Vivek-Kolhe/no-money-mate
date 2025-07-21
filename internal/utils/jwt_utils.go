package utils

import (
	"os"
	"time"

	"github.com/Vivek-Kolhe/no-money-mate/internal/models"
	"github.com/golang-jwt/jwt/v5"
)

func GenerateToken(user models.User) (string, error) {
	claims := jwt.MapClaims{
		"id":    user.ID,
		"email": user.Email,
		"exp":   time.Now().Add(72 * time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(os.Getenv("JWT_SECRET")))
}
