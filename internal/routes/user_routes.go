package routes

import (
	"github.com/Vivek-Kolhe/no-money-mate/internal/controllers"
	"github.com/gofiber/fiber/v2"
)

func RegisterUserRoutes(app fiber.Router, controller *controllers.UserController) {
	userGroup := app.Group("/user")

	userGroup.Post("/register", controller.RegisterUser)
	userGroup.Post("/login", controller.LoginUser)
}
