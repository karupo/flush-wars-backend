package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/karunapo/flush-wars-backend/db"
	"github.com/karunapo/flush-wars-backend/models"
)

func main() {
	db.Init()
	db.DB.AutoMigrate(&models.User{}) // Auto-create table
	// Initialize a new Fiber app
	app := fiber.New()

	// Define a simple health check route
	app.Get("/healthz", func(c *fiber.Ctx) error {
		return c.SendString("OK")
	})

	// Start the server on port 3000
	log.Fatal(app.Listen(":3000"))
}
