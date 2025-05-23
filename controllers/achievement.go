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
	userIDVal := c.Locals("userID")
	userID, ok := userIDVal.(uuid.UUID)
	if !ok {
		log.Println("[GetAchievements] Failed to get user ID from context")
		return fiber.NewError(fiber.StatusUnauthorized, "Unauthorized")
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
