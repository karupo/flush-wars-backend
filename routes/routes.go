package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/karunapo/flush-wars-backend/controllers"
)

// SetupRoutes initializes all application routes and connects them to their handlers.
func SetupRoutes(app *fiber.App) {
	// Public route for user registration.
	app.Post("/register", controllers.Register)
	app.Post("/login", controllers.Login)
	app.Post("/api/poop-log", controllers.CreatePoopLog)
}
