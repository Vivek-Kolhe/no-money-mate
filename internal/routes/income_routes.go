package routes

import (
	"github.com/Vivek-Kolhe/no-money-mate/internal/controllers"
	"github.com/gofiber/fiber/v2"
)

func RegisterIncomeRoutes(app fiber.Router, controller *controllers.IncomeController, middleware fiber.Handler) {
	incomeGroup := app.Group("/income", middleware)

	incomeGroup.Post("/", controller.AddIncome)
	// incomeGroup.Delete(":/id",)
}
