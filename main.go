package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/ivandi1980/my-gofiber/models"
	"github.com/ivandi1980/my-gofiber/service"
	"github.com/ivandi1980/my-gofiber/storage"
	"github.com/joho/godotenv"
	"log"
	"os"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	config := &storage.Config{
		Host:     os.Getenv("DB_HOST"),
		Port:     os.Getenv("DB_PORT"),
		Password: os.Getenv("DB_PASSWORD"),
		User:     os.Getenv("DB_USERNAME"),
		DBName:   os.Getenv("DB_NAME"),
		SSLMode:  os.Getenv("DB_SSLMODE"),
	}

	db, err := storage.NewConnection(config)
	if err != nil {
		log.Fatal("Could not connect to database")
	}

	err = models.MigrateBooks(db)
	if err != nil {
		log.Fatal("Could not migrate Database")
	}

	r := &service.Repository{
		DB: db,
	}

	app := fiber.New()
	r.SetupRoutes(app)

	app.Listen(":1234")
}
