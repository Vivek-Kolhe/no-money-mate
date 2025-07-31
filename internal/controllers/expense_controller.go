package controllers

import (
	"strconv"
	"time"

	"github.com/Vivek-Kolhe/no-money-mate/internal/models"
	"github.com/Vivek-Kolhe/no-money-mate/internal/services"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/v2/mongo"
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
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to add an expense",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Expense added successfully",
		"expense": expense,
	})
}

func (ec *ExpenseController) GetExpenses(c *fiber.Ctx) error {
	user, ok := c.Locals("user").(*models.User)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized access",
		})
	}

	now := time.Now()
	month, year := now.Month(), now.Year()

	if monthQuery := c.Query("month"); monthQuery != "" {
		monthInt, err := strconv.Atoi(monthQuery)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Month and year query param should be an integer",
			})
		}

		if monthInt >= 1 && monthInt <= 12 {
			month = time.Month(monthInt)
		} else {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Month can only be between 1 and 12",
			})
		}
	}

	if yearQuery := c.Query("year"); yearQuery != "" {
		yearInt, err := strconv.Atoi(yearQuery)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Month and year query param should be an integer",
			})
		}

		if yearInt > 0 {
			year = yearInt
		} else {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Year should be greater than 0",
			})
		}
	}

	expenses, err := ec.service.GetExpensesByUserID(user.ID, month, year)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch expenses",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data": expenses,
	})
}

func (ec *ExpenseController) DeleteExpense(c *fiber.Ctx) error {
	user, ok := c.Locals("user").(*models.User)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized access",
		})
	}

	id := c.Params("id")
	expenseID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid expense ID",
		})
	}

	err = ec.service.DeleteExpense(expenseID, user.ID)
	if err != nil {
		if err.Error() == mongo.ErrNoDocuments.Error() {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Expense not found",
			})
		}

		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to delete expense",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Expense deleted successfully",
	})
}
