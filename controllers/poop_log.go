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

// CreatePoopLogInput is the structure for the incoming JSON for poop log
type CreatePoopLogInput struct {
	StoolType string `json:"stool_type"`
	Timestamp string `json:"timestamp"`
	Notes     string `json:"notes,omitempty"`
}

// CreatePoopLog handles the creation of a new poop log entry
func CreatePoopLog(c *fiber.Ctx) error {
	// Parse request body
	var input CreatePoopLogInput
	if err := c.BodyParser(&input); err != nil {
		log.Printf("Failed to parse poop log: %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid input"})
	}

	// Validate stool type
	validStoolTypes := []string{"Normal", "Hard", "Runny", "Soft"}
	isValidType := false
	for _, v := range validStoolTypes {
		if input.StoolType == v {
			isValidType = true
			break
		}
	}
	if !isValidType {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid stool type"})
	}

	// Validate timestamp (ensure it's not in the future)
	timestamp, err := time.Parse(time.RFC3339, input.Timestamp)
	if err != nil || timestamp.After(time.Now()) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid timestamp"})
	}

	// TEMP: Hardcoded user UUID for development/testing
	userID, err := uuid.Parse("8d970d62-8fdb-4d00-a578-47f4977f14ca") // Use a real UUID from your users table
	if err != nil {
		log.Printf("Invalid hardcoded user ID: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Server config error"})
	}

	previousLog := models.PoopLog{}
	var lastLogTime *time.Time

	if err := db.DB.Where("user_id = ?", userID).Order("timestamp desc").First(&previousLog).Error; err == nil {
		lastLogTime = &previousLog.Timestamp
	}

	xp := xp.CalculateXP(input.StoolType, timestamp, lastLogTime)

	// Create a new PoopLog
	poopLog := models.PoopLog{
		UserID:    userID,
		StoolType: input.StoolType,
		Timestamp: timestamp,
		Notes:     input.Notes,
		XPGained:  xp,
	}

	// Store poop log in DB
	if err := db.DB.Create(&poopLog).Error; err != nil {
		log.Printf("Error saving poop log: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Could not store poop log"})
	}

	// Return success response
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message":   "Poop log successfully created",
		"xp_gained": poopLog.XPGained,
		"poop_log":  poopLog,
	})
}
