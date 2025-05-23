package controllers

import (
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/karunapo/flush-wars-backend/db"
	"github.com/karunapo/flush-wars-backend/models"
	"github.com/karunapo/flush-wars-backend/services/xp"
)

// GetXPSummary returns the user's XP, level, streak, and milestones
func GetXPSummary(c *fiber.Ctx) error {
	log.Println("[GetXPSummary] Start")

	userIDVal := c.Locals("userID")
	userID, ok := userIDVal.(uuid.UUID)
	if !ok {
		log.Println("[GetXPSummary] Failed to get user ID from context")
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

func GetUserLevel(c *fiber.Ctx) error {
	log.Println("[GetUserLevel] Start")

	userIDVal := c.Locals("userID")
	userID, ok := userIDVal.(uuid.UUID)
	if !ok {
		log.Println("[GetUserLevel] Failed to get user ID from context")
	}

	// Fetch user
	var user models.User
	if err := db.DB.First(&user, "id = ?", userID).Error; err != nil {
		log.Printf("[GetUserLevel] Failed to fetch user: %v", err)
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to fetch user")
	}

	level := xp.CalculateLevel(user.XP)

	log.Printf("[GetUserLevel] XP: %d | Level: %d", user.XP, level)

	return c.JSON(fiber.Map{
		"level": level,
	})
}

func GetLevelProgress(c *fiber.Ctx) error {
	log.Println("[GetLevelProgress] Start")

	userIDVal := c.Locals("userID")
	userID, ok := userIDVal.(uuid.UUID)
	if !ok {
		log.Println("[GetLevelProgress] Failed to get user ID from context")
	}

	var user models.User
	if err := db.DB.First(&user, "id = ?", userID).Error; err != nil {
		log.Printf("[GetLevelProgress] Failed to fetch user: %v", err)
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to fetch user")
	}

	currentXP := user.XP
	level := xp.CalculateLevel(currentXP)
	currentLevelXP := xp.XPForLevel(level)
	nextLevelXP := xp.XPForLevel(level + 1)
	xpToNext := nextLevelXP - currentXP

	log.Printf("[GetLevelProgress] XP: %d | Level: %d | XP → Next: %d", currentXP, level, xpToNext)

	return c.JSON(fiber.Map{
		"total_xp":         currentXP,
		"level":            level,
		"current_level_xp": currentLevelXP,
		"next_level_xp":    nextLevelXP,
		"xp_to_next_level": xpToNext,
	})
}
