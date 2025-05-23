package controllers

import (
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/karunapo/flush-wars-backend/db"
	"github.com/karunapo/flush-wars-backend/models"
	"github.com/karunapo/flush-wars-backend/services/achievement"
	"github.com/karunapo/flush-wars-backend/services/challenge"
	"github.com/karunapo/flush-wars-backend/services/xp"
)

var validStoolTypes = map[string]bool{
	"Normal": true, "Hard": true, "Runny": true, "Soft": true,
}

type createPoopLogInput struct {
	StoolType string `json:"stool_type"`
	Timestamp string `json:"timestamp"`
	Notes     string `json:"notes,omitempty"`
}

// CreatePoopLog creates the logs
func CreatePoopLog(c *fiber.Ctx) error {
	log.Println("[CreatePoopLog] Request received")

	var input createPoopLogInput
	if err := c.BodyParser(&input); err != nil {
		log.Printf("[CreatePoopLog] Failed to parse body: %v", err)
		return fiber.NewError(fiber.StatusBadRequest, "Invalid input")
	}
	log.Printf("[CreatePoopLog] Parsed input: %+v", input)

	if !validStoolTypes[input.StoolType] {
		log.Printf("[CreatePoopLog] Invalid stool type: %s", input.StoolType)
		return fiber.NewError(fiber.StatusBadRequest, "Invalid stool type")
	}

	timestamp, err := time.Parse(time.RFC3339, input.Timestamp)
	if err != nil || timestamp.After(time.Now()) {
		log.Printf("[CreatePoopLog] Invalid timestamp: %v", input.Timestamp)
		return fiber.NewError(fiber.StatusBadRequest, "Invalid timestamp")
	}

	userIDVal := c.Locals("userID")
	userID, ok := userIDVal.(uuid.UUID)
	if !ok {
		log.Println("[CreatePoopLog] Failed to get user ID from context")
	}

	var lastLog models.PoopLog
	var lastLogTime *time.Time
	if err := db.DB.Where("user_id = ?", userID).Order("timestamp desc").First(&lastLog).Error; err == nil {
		lastLogTime = &lastLog.Timestamp
		log.Printf("[CreatePoopLog] Found last poop log at: %v", *lastLogTime)
	} else {
		log.Printf("[CreatePoopLog] No previous poop log found")
	}

	var user models.User
	if err := db.DB.First(&user, "id = ?", userID).Error; err != nil {
		log.Printf("[CreatePoopLog] Failed to fetch user: %v", err)
		return fiber.NewError(fiber.StatusInternalServerError, "User not found")
	}
	log.Printf("[CreatePoopLog] User before streak update: %+v", user)

	updateStreak(&user, timestamp, lastLogTime)
	log.Printf("[CreatePoopLog] User after streak update: %+v", user)

	xpGained := xp.CalculateXP(input.StoolType, timestamp, lastLogTime)
	log.Printf("[CreatePoopLog] XP calculated: %d", xpGained)

	newLog := models.PoopLog{
		UserID:    userID,
		StoolType: input.StoolType,
		Timestamp: timestamp,
		Notes:     input.Notes,
		XPGained:  xpGained,
	}

	if err := db.DB.Create(&newLog).Error; err != nil {
		log.Printf("[CreatePoopLog] Error saving poop log: %v", err)
		return fiber.NewError(fiber.StatusInternalServerError, "Could not save poop log")
	}
	log.Printf("[CreatePoopLog] Poop log created: %+v", newLog)

	user.XP += xpGained
	if err := db.DB.Save(&user).Error; err != nil {
		log.Printf("[CreatePoopLog] Failed to save user XP: %v", err)
		return fiber.NewError(fiber.StatusInternalServerError, "Could not update XP")
	}
	log.Printf("[CreatePoopLog] User XP updated: %d", user.XP)

	achievementMessage := achievement.CheckAndAwardAchievements(userID)
	log.Printf("[CreatePoopLog] Achievement check result: %v", achievementMessage)

	completedChallengeName, err := challenge.CheckForNewlyCompletedChallenge(user)
	if err != nil {
		log.Printf("[CreatePoopLog] Error checking challenge completion: %v", err)
	}

	log.Println("[CreatePoopLog] Poop log created successfully")
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message":     "Poop log successfully created",
		"xp_gained":   xpGained,
		"poop_log":    newLog,
		"achievement": achievementMessage,
		"total_xp":    user.XP,
		"challenge":   completedChallengeName,
	})
}

func updateStreak(user *models.User, newTime time.Time, lastTime *time.Time) {
	if lastTime == nil {
		user.CurrentStreak = 1
		log.Println("[updateStreak] First poop log — starting new streak")
		return
	}
	diff := newTime.Sub(*lastTime).Hours()
	switch {
	case diff >= 24 && diff <= 48:
		user.CurrentStreak++
		log.Println("[updateStreak] Streak incremented")
	case diff > 48:
		user.CurrentStreak = 1
		log.Println("[updateStreak] Streak reset")
	default:
		log.Println("[updateStreak] Same day or close — streak unchanged")
	}
}

// GetPoopLogHistory lists history
func GetPoopLogHistory(c *fiber.Ctx) error {
	log.Println("[GetPoopLogHistory] Request received")

	userIDVal := c.Locals("userID")
	userID, ok := userIDVal.(uuid.UUID)
	if !ok {
		log.Println("[GetPoopLogHistory] Failed to get user ID from context")
	}

	filter := c.Query("filter")
	page, _ := strconv.Atoi(c.Query("page", "1"))
	pageSize, _ := strconv.Atoi(c.Query("page_size", "3"))
	sortOrder := strings.ToLower(c.Query("sort", "desc"))
	if sortOrder != "asc" {
		sortOrder = "desc"
	}

	log.Printf("[GetPoopLogHistory] Query params: filter=%s, page=%d, page_size=%d, sort=%s", filter, page, pageSize, sortOrder)

	offset := (page - 1) * pageSize
	query := db.DB.Where("user_id = ?", userID)

	switch filter {
	case "epic":
		query = query.Where("stool_type = ?", "Normal")
	case "fail":
		query = query.Where("stool_type IN ?", []string{"Runny", "Hard", "Soft"})
	}

	var logs []models.PoopLog
	if err := query.Order("timestamp " + sortOrder).Limit(pageSize).Offset(offset).Find(&logs).Error; err != nil {
		log.Printf("[GetPoopLogHistory] DB error: %v", err)
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to retrieve poop logs")
	}

	result := make([]fiber.Map, len(logs))
	for i, logEntry := range logs {
		result[i] = fiber.Map{
			"date":       logEntry.Timestamp,
			"stool_type": logEntry.StoolType,
			"xp_gained":  logEntry.XPGained,
		}
	}
	log.Printf("[GetPoopLogHistory] Returned %d logs", len(result))

	return c.JSON(fiber.Map{
		"page":       page,
		"page_size":  pageSize,
		"sort_order": sortOrder,
		"logs":       result,
	})
}

// GetPoopLogByID returns a specific poop log by ID.
func GetPoopLogByID(c *fiber.Ctx) error {
	log.Println("[GetPoopLogByID] Request received")

	logID := c.Params("id")
	log.Printf("[GetPoopLogByID] ID: %s", logID)

	userIDVal := c.Locals("userID")
	userID, ok := userIDVal.(uuid.UUID)
	if !ok {
		log.Println("[GetPoopLogHistory] Failed to get user ID from context")
	}

	poopLogID, err := uuid.Parse(logID)
	if err != nil {
		log.Printf("[GetPoopLogByID] Invalid UUID: %v", err)
		return fiber.NewError(fiber.StatusBadRequest, "Invalid poop log ID")
	}

	var logEntry models.PoopLog
	if err := db.DB.First(&logEntry, "id = ? and user_id =?", poopLogID, userID).Error; err != nil {
		log.Printf("[GetPoopLogByID] Poop log not found: %v", err)
		return fiber.NewError(fiber.StatusNotFound, "Poop log not found")
	}

	log.Printf("[GetPoopLogByID] Returning log: %+v", logEntry)

	return c.JSON(fiber.Map{
		"id":         logEntry.ID,
		"timestamp":  logEntry.Timestamp,
		"stool_type": logEntry.StoolType,
		"xp_gained":  logEntry.XPGained,
	})
}

// UpdatePoopLogByID updates a specific poop log by ID.
func UpdatePoopLogByID(c *fiber.Ctx) error {
	log.Println("[UpdatePoopLogByID] Request received")

	logID := c.Params("id")
	log.Printf("[UpdatePoopLogByID] ID: %s", logID)

	poopLogID, err := uuid.Parse(logID)
	if err != nil {
		log.Printf("[UpdatePoopLogByID] Invalid UUID: %v", err)
		return fiber.NewError(fiber.StatusBadRequest, "Invalid poop log ID")
	}

	type UpdatePoopLogInput struct {
		StoolType string `json:"stool_type"`
	}

	var input UpdatePoopLogInput
	if err := c.BodyParser(&input); err != nil {
		log.Printf("[UpdatePoopLogByID] Failed to parse input: %v", err)
		return fiber.NewError(fiber.StatusBadRequest, "Invalid input")
	}

	log.Printf("[UpdatePoopLogByID] New stool type: %s", input.StoolType)

	var logEntry models.PoopLog
	if err := db.DB.First(&logEntry, "id = ?", poopLogID).Error; err != nil {
		log.Printf("[UpdatePoopLogByID] Log not found: %v", err)
		return fiber.NewError(fiber.StatusNotFound, "Poop log not found")
	}

	logEntry.StoolType = input.StoolType
	// TODO: Recalculate XP if applicable

	if err := db.DB.Save(&logEntry).Error; err != nil {
		log.Printf("[UpdatePoopLogByID] Failed to save changes: %v", err)
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to update poop log")
	}

	log.Printf("[UpdatePoopLogByID] Updated log: %+v", logEntry)

	return c.JSON(fiber.Map{
		"message": "Poop log updated successfully",
		"log":     logEntry,
	})
}

// DeletePoopLogByID deletes a specific poop log by ID.
func DeletePoopLogByID(c *fiber.Ctx) error {
	log.Println("[DeletePoopLogByID] Request received")

	logID := c.Params("id")
	log.Printf("[DeletePoopLogByID] ID: %s", logID)

	poopLogID, err := uuid.Parse(logID)
	if err != nil {
		log.Printf("[DeletePoopLogByID] Invalid UUID: %v", err)
		return fiber.NewError(fiber.StatusBadRequest, "Invalid poop log ID")
	}

	var logEntry models.PoopLog
	if err := db.DB.First(&logEntry, "id = ?", poopLogID).Error; err != nil {
		log.Printf("[DeletePoopLogByID] Log not found: %v", err)
		return fiber.NewError(fiber.StatusNotFound, "Poop log not found")
	}

	if err := db.DB.Delete(&logEntry).Error; err != nil {
		log.Printf("[DeletePoopLogByID] Failed to delete: %v", err)
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to delete poop log")
	}

	log.Printf("[DeletePoopLogByID] Deleted log ID: %s", poopLogID)

	return c.JSON(fiber.Map{
		"message": "Poop log deleted successfully",
		"log_id":  poopLogID,
	})
}
