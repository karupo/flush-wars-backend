package controllers

import (
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/karunapo/flush-wars-backend/db"
	"github.com/karunapo/flush-wars-backend/models"
)

type ChallengeResponse struct {
	ID          uuid.UUID `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	StartDate   time.Time `json:"start_date"`
	EndDate     time.Time `json:"end_date"`
	RewardXP    int       `json:"reward_xp"`
}

func GetAllChallenges(c *fiber.Ctx) error {
	log.Println("[GetAllChallenges] Request received")
	var challenges []models.Challenge
	if err := db.DB.Find(&challenges).Error; err != nil {
		log.Printf("[GetAllChallenges] Error fetching challenges: %v", err)
		return c.Status(500).JSON(fiber.Map{"error": "Could not fetch challenges"})
	}
	log.Printf("[GetAllChallenges] Found %d challenges", len(challenges))
	return c.JSON(challenges)
}

func JoinChallenge(c *fiber.Ctx) error {
	log.Println("[JoinChallenge] Request received")
	challengeID := c.Params("id")
	userIDVal := c.Locals("userID")
	userID, ok := userIDVal.(uuid.UUID)
	if !ok {
		log.Println("[JoinChallenge] Failed to get user ID from context")
	}

	log.Printf("[JoinChallenge] User %s joining challenge %s", userID, challengeID)

	var existing models.UserChallenge
	if err := db.DB.Where("user_id = ? AND challenge_id = ?", userID, challengeID).First(&existing).Error; err == nil {
		log.Println("[JoinChallenge] User already joined this challenge")
		return c.Status(400).JSON(fiber.Map{"error": "Already joined"})
	}

	join := models.UserChallenge{
		UserID:      userID,
		ChallengeID: uuid.MustParse(challengeID),
		JoinedAt:    time.Now(),
		Completed:   false,
	}

	if err := db.DB.Create(&join).Error; err != nil {
		log.Printf("[JoinChallenge] Error joining challenge: %v", err)
		return c.Status(500).JSON(fiber.Map{"error": "Could not join challenge"})
	}
	log.Println("[JoinChallenge] Challenge joined successfully")
	return c.JSON(fiber.Map{"message": "Challenge joined"})
}

func GetUserChallenges(c *fiber.Ctx) error {
	log.Println("[GetUserChallenges] Request received")
	userIDVal := c.Locals("userID")
	userID, ok := userIDVal.(uuid.UUID)
	if !ok {
		log.Println("[GetUserChallenges] Failed to get user ID from context")
	}

	var userChallenges []models.UserChallenge
	if err := db.DB.Where("user_id = ?", userID).Find(&userChallenges).Error; err != nil {
		log.Printf("[GetUserChallenges] Error fetching user challenges: %v", err)
		return c.Status(500).JSON(fiber.Map{"error": "Could not fetch user challenges"})
	}
	log.Printf("[GetUserChallenges] Found %d user challenges", len(userChallenges))
	return c.JSON(userChallenges)
}

func GetChallengeProgress(c *fiber.Ctx) error {
	log.Println("[GetChallengeProgress] Request received")
	challengeID := c.Params("id")
	userIDVal := c.Locals("userID")
	userID, ok := userIDVal.(uuid.UUID)
	if !ok {
		log.Println("[GetChallengeProgress] Failed to get user ID from context")
	}

	var progress models.UserChallenge
	if err := db.DB.Where("user_id = ? AND challenge_id = ?", userID, challengeID).First(&progress).Error; err != nil {
		log.Printf("[GetChallengeProgress] Error fetching progress: %v", err)
		return c.Status(404).JSON(fiber.Map{"error": "Challenge not joined or not found"})
	}
	log.Printf("[GetChallengeProgress] Progress fetched for challenge %s", challengeID)
	return c.JSON(progress)
}

func CompleteChallenge(c *fiber.Ctx) error {
	log.Println("[CompleteChallenge] Request received")
	challengeID := c.Params("id")
	userIDVal := c.Locals("userID")
	userID, ok := userIDVal.(uuid.UUID)
	if !ok {
		log.Println("[CompleteChallenge] Failed to get user ID from context")
	}

	var progress models.UserChallenge
	if err := db.DB.Where("user_id = ? AND challenge_id = ?", userID, challengeID).First(&progress).Error; err != nil {
		log.Printf("[CompleteChallenge] Error finding challenge: %v", err)
		return c.Status(404).JSON(fiber.Map{"error": "Challenge not found"})
	}

	if progress.Completed {
		log.Println("[CompleteChallenge] Challenge already completed")
		return c.Status(400).JSON(fiber.Map{"error": "Challenge already completed"})
	}

	progress.Completed = true
	now := time.Now()
	progress.CompletedAt = &now

	if err := db.DB.Save(&progress).Error; err != nil {
		log.Printf("[CompleteChallenge] Error updating challenge: %v", err)
		return c.Status(500).JSON(fiber.Map{"error": "Failed to mark as complete"})
	}

	log.Printf("[CompleteChallenge] Challenge %s marked as complete for user %s", challengeID, userID)
	return c.JSON(fiber.Map{"message": "Challenge completed"})
}
