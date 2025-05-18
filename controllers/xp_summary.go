package controllers

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/karunapo/flush-wars-backend/db"
	"github.com/karunapo/flush-wars-backend/models"
	"github.com/karunapo/flush-wars-backend/xp"
)

// GetXPSummary for a user with level and milestone
func GetXPSummary(c *fiber.Ctx) error {
	// TEMP: Hardcoded user ID
	userID, _ := uuid.Parse("8d970d62-8fdb-4d00-a578-47f4977f14ca")

	// Fetch logs
	var logs []models.PoopLog
	if err := db.DB.Where("user_id = ?", userID).Order("timestamp asc").Find(&logs).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to fetch logs"})
	}

	var totalXP int
	var streak int
	var lastDay time.Time

	for _, log := range logs {
		if lastDay.IsZero() || log.Timestamp.Sub(lastDay).Hours() <= 48 {
			streak++
		} else {
			streak = 1 // reset
		}
		lastDay = log.Timestamp
	}

	var user models.User
	if err := db.DB.First(&user, "id = ?", userID).Error; err != nil {
		return err
	}

	level := xp.CalculateLevel(user.XP)
	milestones := xp.GetMilestones(totalXP, streak)

	return c.JSON(fiber.Map{
		"total_xp":   user.XP,
		"level":      level,
		"streak":     streak,
		"milestones": milestones,
	})
}
