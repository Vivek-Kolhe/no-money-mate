package controllers

import (
	"log"
	"time"

	"github.com/Vivek-Kolhe/no-money-mate/internal/models"
	"github.com/Vivek-Kolhe/no-money-mate/internal/services"
	"github.com/gofiber/fiber/v2"
)

type ExpenseController struct {
	service *services.ExpenseService
}

func NewExpenseController(service *services.ExpenseService) *ExpenseController {
	return &ExpenseController{service: service}
}

func (ec *ExpenseController) AddExpense(c *fiber.Ctx) error {
	var expense models.Expense

	if err := c.BodyParser(&expense); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid JSON input",
		})
	}

	if expense.Amount <= 0.0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Amount should be greater than zero",
		})
	}

	if expense.Category == "" || expense.Icon == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Expense category or icon is missing",
		})
	}

	user, ok := c.Locals("user").(*models.User)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized user or invalid token",
		})
	}

	if expense.Date.IsZero() {
		expense.Date = time.Now().UTC().Truncate(24 * time.Hour)
	}

	expense.UserID = user.ID

	if err := ec.service.CreateExpense(expense); err != nil {
		log.Panic(err.Error())
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to add an expense",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Expense added successfully",
		"expense": expense,
	})
}
