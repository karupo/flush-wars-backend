package middleware

import (
	"log"
	"os"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

// RequireAuth validates JWT from the Authorization header and sets user ID in context.
func RequireAuth(c *fiber.Ctx) error {
	log.Printf("[RequireAuth] JWT_SECRET: %s", os.Getenv("JWT_SECRET"))

	authHeader := c.Get("Authorization")
	log.Printf("[RequireAuth] Authorization header: %s", authHeader)

	if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
		log.Println("[RequireAuth] Missing or invalid authorization header")
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Missing or invalid authorization header"})
	}

	tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
	log.Printf("[RequireAuth] Token string: %s", tokenStr)

	token, err := jwt.Parse(tokenStr, func(_ *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_SECRET")), nil
	})

	if err != nil || !token.Valid {
		log.Printf("[RequireAuth] Invalid token: %v", err)
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid or expired token"})
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		log.Println("[RequireAuth] Failed to cast claims to jwt.MapClaims")
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid token claims"})
	}
	log.Printf("[RequireAuth] Token claims: %+v", claims)

	claimUserID, exists := claims["user_id"]
	if !exists {
		log.Println("[RequireAuth] user_id claim missing")
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid token claims"})
	}

	userIDStr, ok := claimUserID.(string)
	if !ok {
		log.Printf("[RequireAuth] user_id claim is not a string: %v", claimUserID)
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid user ID in token"})
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		log.Printf("[RequireAuth] Malformed user ID: %v", err)
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Malformed user ID"})
	}

	log.Printf("[RequireAuth] Authenticated user ID: %s", userID.String())

	// Store user ID in context
	c.Locals("userID", userID)

	return c.Next()
}
