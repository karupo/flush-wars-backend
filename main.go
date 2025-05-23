// Package main contains the main go code
package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/karunapo/flush-wars-backend/db"
	"github.com/karunapo/flush-wars-backend/models"
	"github.com/karunapo/flush-wars-backend/routes"
)

func main() {
	// Auto-create table
	// Initialize a new Fiber app
	// Define a simple health check route
	app := InitApp()

	// Start the server on port 3000
	log.Fatal(app.Listen(":3000"))
}

// InitApp - Initialize the database connection, Initialize a new Fiber app
func InitApp() *fiber.App {
	// Initialize the database connection
	db.Init(true)

	// Auto-migrate the models (creating/updating tables)
	db.DB.AutoMigrate(&models.User{})
	db.DB.AutoMigrate(&models.PoopLog{})
	db.DB.AutoMigrate(&models.Achievement{})
	db.DB.AutoMigrate(&models.Challenge{}, &models.UserChallenge{})
	db.SeedChallenges()

	// Initialize a new Fiber app
	app := fiber.New()
	routes.SetupRoutes(app)

	// Define a simple health check route
	app.Get("/healthz", func(c *fiber.Ctx) error {
		return c.SendString("OK")
	})

	return app
}
