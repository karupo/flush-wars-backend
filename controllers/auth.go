package controllers

import (
	"os"
	"time"

	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/karunapo/flush-wars-backend/db"
	"github.com/karunapo/flush-wars-backend/models"
	"golang.org/x/crypto/bcrypt"
)

// RegisterInput represents the expected payload for user registration.
type RegisterInput struct {
	Email    string `json:"email"`
	Username string `json:"username"`
	Password string `json:"password"`
}

// Register handles user registration: validates input, hashes password, and stores user.
func Register(c *fiber.Ctx) error {
	var input RegisterInput
	log.Printf("Login attempt for email: %s", input.Email)

	// Parse the incoming JSON request body into the input struct.
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid input"})
	}

	// Basic validation to ensure required fields are not empty.
	if input.Email == "" || input.Username == "" || input.Password == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "All fields are required"})
	}

	// Securely hash the user's password using bcrypt.
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Could not hash password"})
	}

	// Create a new User instance with the provided data and hashed password.
	user := models.User{
		Email:    input.Email,
		Username: input.Username,
		Password: string(hashedPassword),
	}

	// Attempt to store the user in the database.
	if err := db.DB.Create(&user).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Could not create user"})
	}

	// Return a success response with selected user info (excluding password).
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "User registered successfully",
		"user": fiber.Map{
			"id":       user.ID,
			"email":    user.Email,
			"username": user.Username,
		},
	})
}

// LoginInput
type LoginInput struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// Login
func Login(c *fiber.Ctx) error {
	var input LoginInput
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid input"})
	}

	var user models.User
	if err := db.DB.Where("email = ?", input.Email).First(&user).Error; err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid email or password"})
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid email or password"})
	}

	// Generate JWT token
	secret := os.Getenv("JWT_SECRET")
	claims := jwt.MapClaims{
		"user_id":  user.ID,
		"email":    user.Email,
		"username": user.Username,
		"exp":      time.Now().Add(time.Hour * 72).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(secret))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Could not sign token"})
	}

	return c.JSON(fiber.Map{
		"token": signedToken,
	})
}
