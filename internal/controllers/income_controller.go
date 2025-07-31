package controllers

import (
	"time"

	"github.com/Vivek-Kolhe/no-money-mate/internal/models"
	"github.com/Vivek-Kolhe/no-money-mate/internal/services"
	"github.com/gofiber/fiber/v2"
)

type IncomeController struct {
	service *services.IncomeService
}

func NewIncomeController(service *services.IncomeService) *IncomeController {
	return &IncomeController{service: service}
}

func (ic *IncomeController) AddIncome(c *fiber.Ctx) error {
	var income models.Income

	if err := c.BodyParser(&income); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid JSON input",
		})
	}

	if income.Amount <= 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Amount should be greater than zero",
		})
	}

	if income.Icon == "" || income.Source == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Income source or icon is missing",
		})
	}

	user, ok := c.Locals("user").(*models.User)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized user or invalid token",
		})
	}

	if income.Date.IsZero() {
		income.Date = time.Now().UTC().Truncate(24 * time.Hour)
	}

	income.UserID = user.ID
	insertedID, err := ic.service.CreateIncome(income)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to add an income source",
		})
	}
	income.ID = insertedID

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Income added successfully",
		"income":  income,
	})
}
