package controllers

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/karunapo/flush-wars-backend/db"
	"github.com/karunapo/flush-wars-backend/models"
)

// GetAchievements lists all achievements for the authenticated user
func GetAchievements(c *fiber.Ctx) error {
	// In production, extract userID from JWT/session instead of hardcoding
	userIDStr := "2f9f3c05-75b0-4935-9d89-f074715f5c19"
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		log.Printf("Invalid UUID: %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid user ID",
		})
	}

	var achievements []models.Achievement
	if err := db.DB.Where("user_id = ?", userID).Find(&achievements).Error; err != nil {
		log.Printf("DB error fetching achievements: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch achievements",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"user_id":      userID,
		"achievements": achievements,
	})
}
