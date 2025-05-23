package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/karunapo/flush-wars-backend/controllers"
	"github.com/karunapo/flush-wars-backend/middleware"
)

// SetupRoutes initializes all application routes and connects them to their handlers.
func SetupRoutes(app *fiber.App) {
	// Public route for user registration.
	app.Post("/register", controllers.Register)
	app.Post("/login", controllers.Login)

	poop := app.Group("/api/poop-log")

	poop.Post("/", middleware.RequireAuth, controllers.CreatePoopLog)
	poop.Get("/history", middleware.RequireAuth, controllers.GetPoopLogHistory)
	poop.Get("/:id", middleware.RequireAuth, controllers.GetPoopLogByID)
	poop.Put("/:id", middleware.RequireAuth, controllers.UpdatePoopLogByID)
	poop.Delete("/:id", middleware.RequireAuth, controllers.DeletePoopLogByID)

	app.Get("/api/achievements", middleware.RequireAuth, controllers.GetAchievements)

	leaderboard := app.Group("/api/leaderboard")
	leaderboard.Get("/", middleware.RequireAuth, controllers.GetLeaderboard)
	leaderboard.Get("/weekly", middleware.RequireAuth, controllers.GetWeeklyLeaderboard)

	user := app.Group("/api/user")
	user.Get("/profile", middleware.RequireAuth, middleware.RequireAuth, controllers.GetUserProfile)
	user.Put("/profile", middleware.RequireAuth, controllers.UpdateUserProfile)
	user.Get("/xp-summary", middleware.RequireAuth, controllers.GetXPSummary)
	user.Get("/level", middleware.RequireAuth, controllers.GetUserLevel)
	user.Get("/level-progress", middleware.RequireAuth, controllers.GetLevelProgress)
	user.Get("/streak", middleware.RequireAuth, controllers.GetUserStreak)

	analytics := app.Group("/api/analytics")
	analytics.Get("/weekly", middleware.RequireAuth, controllers.GetAnalyticsWeekly)
	analytics.Get("/monthly", middleware.RequireAuth, controllers.GetAnalyticsMonthly)
	analytics.Get("/yearly", middleware.RequireAuth, controllers.GetAnalyticsYearly)
	analytics.Get("/trends", middleware.RequireAuth, controllers.GetAnalyticsTrends)

	challenges := app.Group("/api/challenges")

	challenges.Get("/", controllers.GetAllChallenges)
	challenges.Post("/:id/join", middleware.RequireAuth, controllers.JoinChallenge)
	challenges.Get("/user", middleware.RequireAuth, controllers.GetUserChallenges)
	challenges.Get("/:id/progress", middleware.RequireAuth, controllers.GetChallengeProgress)
	challenges.Post("/:id/complete", middleware.RequireAuth, controllers.CompleteChallenge)
}
