package main

import (
	"log"

	"github.com/Vivek-Kolhe/no-money-mate/internal/controllers"
	"github.com/Vivek-Kolhe/no-money-mate/internal/middlewares"
	"github.com/Vivek-Kolhe/no-money-mate/internal/models"
	"github.com/Vivek-Kolhe/no-money-mate/internal/repository"
	"github.com/Vivek-Kolhe/no-money-mate/internal/routes"
	"github.com/Vivek-Kolhe/no-money-mate/internal/services"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func main() {
	// err := godotenv.Load()
	// if err != nil {
	// 	log.Fatal("Error loading/finding .env file")
	// }

	db := models.ConnectDB()
	app := fiber.New()

	app.Use(cors.New(cors.Config{
		AllowOrigins: "*", // Allow all origins
		AllowMethods: "GET,POST,PUT,DELETE,OPTIONS",
		AllowHeaders: "Origin, Content-Type, Accept, Authorization",
	}))

	userRepo := repository.NewUserRepository(db)
	userService := services.NewUserService(userRepo)
	userController := controllers.NewUserController(userService)

	expenseRepo := repository.NewExpenseRepository(db)
	expenseService := services.NewExpenseService(expenseRepo)
	expenseController := controllers.NewExpenseController(expenseService)

	incomeRepo := repository.NewIncomeRepository(db)
	incomeService := services.NewIncomeService(incomeRepo)
	incomeController := controllers.NewIncomeController(incomeService)

	dashboardRepo := repository.NewDashboardRepository(db)
	dashboardService := services.NewDashboardService(dashboardRepo)
	dashboardController := controllers.NewDashboardController(dashboardService)

	api := app.Group("/api")
	authMiddleware := middlewares.JWTAuth(userRepo)

	routes.RegisterUserRoutes(api, userController, authMiddleware)
	routes.RegisterExpenseRoutes(api, expenseController, authMiddleware)
	routes.RegisterIncomeRoutes(api, incomeController, authMiddleware)
	routes.RegisterDashboardRoutes(api, dashboardController, authMiddleware)

	app.Get("/health", func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"status":  "ok",
			"message": "alive",
		})
	})

	log.Fatal(app.Listen(":3000"))
}
