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

// GetUserProfile returns basic user info (username and email)
func GetUserProfile(c *fiber.Ctx) error {
	// TODO: Replace hardcoded UUID with authenticated user ID
	userID, err := uuid.Parse("2f9f3c05-75b0-4935-9d89-f074715f5c19")
	if err != nil {
		log.Printf("[GetUserProfile] Invalid hardcoded user ID: %v", err)
		return fiber.NewError(fiber.StatusInternalServerError, "Invalid user configuration")
	}

	var user models.User
	if err := db.DB.First(&user, "id = ?", userID).Error; err != nil {
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
	// TODO: Replace hardcoded UUID with authenticated user ID
	userID, err := uuid.Parse("2f9f3c05-75b0-4935-9d89-f074715f5c19")
	if err != nil {
		log.Printf("[UpdateUserProfile] Invalid hardcoded user ID: %v", err)
		return fiber.NewError(fiber.StatusInternalServerError, "Invalid user configuration")
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
