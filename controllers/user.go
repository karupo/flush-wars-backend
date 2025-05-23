package controllers

import (
	"log"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/karunapo/flush-wars-backend/db"
	"github.com/karunapo/flush-wars-backend/models"
)

type UserProfile struct {
	UserName string `json:"username"`
	Email    string `json:"email"`
}

func GetUserProfile(c *fiber.Ctx) error {
	userIDVal := c.Locals("userID")
	userID, ok := userIDVal.(uuid.UUID)
	if !ok {
		log.Println("[GetUserProfile] Failed to get user ID from context")
		return fiber.NewError(fiber.StatusUnauthorized, "Unauthorized")
	}

	var user models.User
	if err := db.DB.First(&user, "id = ?", userID).Error; err != nil {
		log.Printf("[GetUserProfile] Could not fetch user %s: %v", userID, err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Could not fetch user profile",
		})
	}

	userProfile := UserProfile{
		UserName: user.Username,
		Email:    user.Email,
	}

	return c.JSON(userProfile)
}

// UpdateUserProfile allows the user to update their username
func UpdateUserProfile(c *fiber.Ctx) error {
	userIDVal := c.Locals("userID")
	userID, ok := userIDVal.(uuid.UUID)
	if !ok {
		log.Println("[UpdateUserProfile] Failed to get user ID from context")
		return fiber.NewError(fiber.StatusUnauthorized, "Unauthorized")
	}

	type UpdateUsernameInput struct {
		Username string `json:"username"`
	}

	var input UpdateUsernameInput
	if err := c.BodyParser(&input); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid input format")
	}

	if strings.TrimSpace(input.Username) == "" {
		return fiber.NewError(fiber.StatusBadRequest, "Username cannot be empty")
	}

	var user models.User
	if err := db.DB.First(&user, "id = ?", userID).Error; err != nil {
		return fiber.NewError(fiber.StatusNotFound, "User not found")
	}

	user.Username = input.Username

	if err := db.DB.Save(&user).Error; err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to update username")
	}

	return c.JSON(fiber.Map{
		"message": "Username updated successfully",
		"user": fiber.Map{
			"username": user.Username,
			"email":    user.Email,
		},
	})
}

func GetUserStreak(c *fiber.Ctx) error {
	userIDVal := c.Locals("userID")
	userID, ok := userIDVal.(uuid.UUID)
	if !ok {
		log.Println("[GetUserStreak] Failed to get user ID from context")
		return fiber.NewError(fiber.StatusUnauthorized, "Unauthorized")
	}

	var user models.User
	if err := db.DB.First(&user, "id = ?", userID).Error; err != nil {
		log.Printf("[GetUserStreak] User not found: %v", err)
		return fiber.NewError(fiber.StatusNotFound, "User not found")
	}

	log.Printf("[GetUserStreak] User %s current streak: %d", userID, user.CurrentStreak)

	return c.JSON(fiber.Map{
		"user_id":        userID,
		"current_streak": user.CurrentStreak,
	})
}
