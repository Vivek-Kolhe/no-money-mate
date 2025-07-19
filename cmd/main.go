package main

import (
	"log"

	"github.com/Vivek-Kolhe/no-money-mate/internal/models"
	"github.com/Vivek-Kolhe/no-money-mate/internal/repository"
	"github.com/gofiber/fiber/v3"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading/finding .env file")
	}

	db := models.ConnectDB()
	userRepo := repository.NewUserRepository(db)

	err = userRepo.CreateUser(models.User{FirstName: "Vivek", LastName: "Kolhe", Email: "vivek@example.com", Password: "123456"})
	if err != nil {
		log.Fatal(err)
	}

	app := fiber.New()

	log.Fatal(app.Listen(":3000"))
}
