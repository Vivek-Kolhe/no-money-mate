package routes

import (
	"github.com/Vivek-Kolhe/no-money-mate/internal/controllers"
	"github.com/gofiber/fiber/v2"
)

func RegisterExpenseRoutes(app fiber.Router, controller *controllers.ExpenseController, middleware fiber.Handler) {
	userGroup := app.Group("/expense", middleware)

	userGroup.Post("/", controller.AddExpense)
	userGroup.Get("/all", controller.GetExpenses)
}
