package controllers

import (
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/karunapo/flush-wars-backend/db"
	"github.com/karunapo/flush-wars-backend/models"
	"github.com/karunapo/flush-wars-backend/xp"
)

// GetXPSummary returns the user's XP, level, streak, and milestones
func GetXPSummary(c *fiber.Ctx) error {
	log.Println("[GetXPSummary] Start")

	// TEMP: Replace with real user ID from auth
	userID, err := uuid.Parse("2f9f3c05-75b0-4935-9d89-f074715f5c19")
	if err != nil {
		log.Printf("[GetXPSummary] Invalid user ID: %v", err)
		return fiber.NewError(fiber.StatusInternalServerError, "Invalid user ID")
	}

	// Fetch poop logs sorted by date
	var logs []models.PoopLog
	if err := db.DB.Where("user_id = ?", userID).Order("timestamp asc").Find(&logs).Error; err != nil {
		log.Printf("[GetXPSummary] Failed to fetch logs: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to fetch logs"})
	}

	var streak int
	var lastLogTime time.Time

	for _, logEntry := range logs {
		if lastLogTime.IsZero() {
			streak = 1
		} else {
			diff := logEntry.Timestamp.Sub(lastLogTime).Hours()

			switch {
			case diff >= 24 && diff <= 48:
				streak++
			case diff > 48:
				streak = 1
				// diff < 24 → same day: streak unchanged
			}
		}
		lastLogTime = logEntry.Timestamp
	}

	// Fetch user record
	var user models.User
	if err := db.DB.First(&user, "id = ?", userID).Error; err != nil {
		log.Printf("[GetXPSummary] Failed to fetch user: %v", err)
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to fetch user")
	}

	level := xp.CalculateLevel(user.XP)

	log.Printf("[GetXPSummary] XP: %d | Level: %d | Streak: %d", user.XP, level, streak)

	return c.JSON(fiber.Map{
		"total_xp": user.XP,
		"level":    level,
		"streak":   streak,
	})
}
