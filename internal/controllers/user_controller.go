package controllers

import (
	"github.com/Vivek-Kolhe/no-money-mate/internal/models"
	"github.com/Vivek-Kolhe/no-money-mate/internal/services"
	"github.com/gofiber/fiber/v2"
)

type UserController struct {
	service *services.UserService
}

func NewUserController(service *services.UserService) *UserController {
	return &UserController{service: service}
}

func (uc *UserController) RegisterUser(c *fiber.Ctx) error {
	var user models.User

	if err := c.BodyParser(&user); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid JSON input",
		})
	}

	if user.Email == "" || user.FirstName == "" || user.LastName == "" || user.Password == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Name, email and password is required",
		})
	}

	err := uc.service.CreateUser(user)
	if err != nil {
		if err.Error() == "user already exists" {
			return c.Status(fiber.StatusConflict).JSON(fiber.Map{
				"error": "User with this email already exists",
			})
		}

		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create user",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "User registered successfully",
	})
}
