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

func (ic *IncomeController) GetIncomes(c *fiber.Ctx) error {
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

	incomes, err := ic.service.GetIncomeByUserID(user.ID, month, year)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch incomes",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data": incomes,
	})
}

func (ic *IncomeController) DeleteIncome(c *fiber.Ctx) error {
	user, ok := c.Locals("user").(*models.User)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized access",
		})
	}

	id := c.Params("id")
	incomeID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid income ID",
		})
	}

	err = ic.service.DeleteIncome(incomeID, user.ID)
	if err != nil {
		if err.Error() == mongo.ErrNoDocuments.Error() {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Income not found",
			})
		}

		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to delete income",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Income deleted successfully",
	})
}
