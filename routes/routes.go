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

	poop := app.Group("/api/poop-log")

	poop.Post("/", controllers.CreatePoopLog)
	poop.Get("/history", controllers.GetPoopLogHistory)
	poop.Get("/:id", controllers.GetPoopLogByID)
	poop.Put("/:id", controllers.UpdatePoopLogByID)
	poop.Delete("/:id", controllers.DeletePoopLogByID)

	app.Get("/api/achievements", controllers.GetAchievements)
	app.Get("/api/leaderboard", controllers.GetLeaderboard)

	user := app.Group("/api/user")
	user.Get("/profile", controllers.GetUserProfile)
	user.Put("/profile", controllers.UpdateUserProfile)
	user.Get("/xp-summary", controllers.GetXPSummary)
	user.Get("/level", controllers.GetUserLevel)
	user.Get("/level-progress", controllers.GetLevelProgress)
}
