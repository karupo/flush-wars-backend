package controllers

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
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

// WeeklyLeaderboardEntry represents a leaderboard entry for a user in a given week.
type WeeklyLeaderboardEntry struct {
	UserID   uuid.UUID `json:"-"`
	Username string    `json:"username"`
	XP       int       `json:"xp"`
	Title    string    `json:"title"`
	Rank     int       `json:"rank"`
}

// GetWeeklyLeaderboard returns the top users based on XP for the current week.
func GetWeeklyLeaderboard(c *fiber.Ctx) error {
	type Result struct {
		UserID uuid.UUID
		XP     int
	}

	oneWeekAgo := time.Now().AddDate(0, 0, -7)

	var results []Result

	// Query: sum XP grouped by user_id for logs in last 7 days, ordered desc, limit 10
	err := db.DB.
		Model(&models.PoopLog{}).
		Select("user_id, SUM(xp_gained) as xp").
		Where("timestamp >= ?", oneWeekAgo).
		Group("user_id").
		Order("xp DESC").
		Limit(10).
		Scan(&results).Error
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to get weekly leaderboard"})
	}

	entries := make([]WeeklyLeaderboardEntry, 0, len(results))

	for i, r := range results {
		var user models.User
		if err := db.DB.First(&user, "id = ?", r.UserID).Error; err != nil {
			// Skip user if not found
			continue
		}
		entries = append(entries, WeeklyLeaderboardEntry{
			Username: user.Username,
			XP:       r.XP,
			Title:    achievement.GetTitle(i),
			Rank:     i + 1,
		})
	}

	return c.JSON(entries)
}
