package routes

import (
	"github.com/Vivek-Kolhe/no-money-mate/internal/controllers"
	"github.com/gofiber/fiber/v2"
)

func RegisterDashboardRoutes(app fiber.Router, controller *controllers.DashboardController, middleware fiber.Handler) {
	dashboardGroup := app.Group("/dashboard", middleware)

	dashboardGroup.Get("/", controller.GetDashboardData)
}
