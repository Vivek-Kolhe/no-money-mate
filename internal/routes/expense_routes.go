package routes

import (
	"github.com/Vivek-Kolhe/no-money-mate/internal/controllers"
	"github.com/gofiber/fiber/v2"
)

func RegisterExpenseRoutes(app fiber.Router, controller *controllers.ExpenseController, middleware fiber.Handler) {
	expenseGroup := app.Group("/expense", middleware)

	expenseGroup.Post("/", controller.AddExpense)
	expenseGroup.Delete("/:id", controller.DeleteExpense)
	expenseGroup.Get("/all", controller.GetExpenses)
}
