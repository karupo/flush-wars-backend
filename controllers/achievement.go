package controllers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/karunapo/flush-wars-backend/db"
	"github.com/karunapo/flush-wars-backend/models"
)

// GetAchievements list all the achievements of the user
func GetAchievements(c *fiber.Ctx) error {
	userID := "2f9f3c05-75b0-4935-9d89-f074715f5c19" // Replace with real user ID from auth later
	var achievement []models.Achievement
	if err := db.DB.Where("user_id = ?", userID).Find(&achievement).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to fetch achievements"})
	}
	return c.JSON(achievement)
}
