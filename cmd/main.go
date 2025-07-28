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
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading/finding .env file")
	}

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

	api := app.Group("/api")
	authMiddleware := middlewares.JWTAuth(userRepo)

	routes.RegisterUserRoutes(api, userController)
	routes.RegisterExpenseRoutes(api, expenseController, authMiddleware)

	// password, err := utils.HashPassword("123456")
	// if err != nil {
	// 	log.Fatal("Something went wrong while hashing password")
	// }

	// err = userRepo.CreateUser(models.User{FirstName: "Vivek", LastName: "Kolhe", Email: "vivek@example.com", Password: password})
	// if err != nil {
	// 	log.Fatal(err)
	// }

	log.Fatal(app.Listen(":3000"))
}
