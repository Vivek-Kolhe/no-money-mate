package middlewares

import (
	"log"
	"strings"

	"github.com/Vivek-Kolhe/no-money-mate/internal/repository"
	"github.com/Vivek-Kolhe/no-money-mate/internal/utils"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

func JWTAuth(repo *repository.UserRepository) fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Missing or invalid Authorization header",
			})
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		claims, err := utils.ValidateToken(tokenString)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		email, ok := claims["email"].(string)
		if !ok || email == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid token payload",
			})
		}

		user, err := repo.FindByEmail(email)
		if err != nil {
			if err.Error() == mongo.ErrNoDocuments.Error() {
				return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
					"error": "Account not found",
				})
			}

			log.Panic(err.Error())
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Internal server error! Something went wrong",
			})
		}

		c.Locals("user", user)

		return c.Next()
	}
}
