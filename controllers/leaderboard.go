package controllers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/karunapo/flush-wars-backend/db"
	"github.com/karunapo/flush-wars-backend/models"
	"github.com/karunapo/flush-wars-backend/services/achievement"
)

// LeaderboardEntry to display the leaderboard
type LeaderboardEntry struct {
	Username string `json:"username"`
	XP       int    `json:"xp"`
	Title    string `json:"title"`
	Rank     int    `json:"rank"`
}

// GetLeaderboard displays the top users
func GetLeaderboard(c *fiber.Ctx) error {
	var users []models.User

	// Get top 10 users ordered by XP
	if err := db.DB.Order("xp desc").Limit(10).Find(&users).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Could not fetch leaderboard"})
	}

	entries := []LeaderboardEntry{}
	for i, user := range users {
		entries = append(entries, LeaderboardEntry{
			Username: user.Username,
			XP:       user.XP,
			Title:    achievement.GetTitle(i), // Assign anime-inspired title
			Rank:     i + 1,
		})
	}

	return c.JSON(entries)
}
