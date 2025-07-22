package utils

import (
	"errors"
	"os"
	"time"

	"github.com/Vivek-Kolhe/no-money-mate/internal/models"
	"github.com/golang-jwt/jwt/v5"
)

var (
	ErrInvalidSigningMethod = errors.New("invalid or unexpected signing method")
	ErrInvalidToken         = errors.New("invalid or expired token")
	ErrInvalidClaims        = errors.New("invalid claims")
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

func ValidateToken(tokenString string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrInvalidSigningMethod
		}
		return []byte(os.Getenv("JWT_SECRET")), nil
	})

	if err != nil || !token.Valid {
		return nil, ErrInvalidToken
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, ErrInvalidClaims
	}

	return claims, nil
}
