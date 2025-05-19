package controllers

import (
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/karunapo/flush-wars-backend/achievement"
	"github.com/karunapo/flush-wars-backend/db"
	"github.com/karunapo/flush-wars-backend/models"
	"github.com/karunapo/flush-wars-backend/xp"
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

	userID, err := uuid.Parse("2f9f3c05-75b0-4935-9d89-f074715f5c19") // Replace later
	if err != nil {
		log.Printf("[CreatePoopLog] Invalid hardcoded user ID: %v", err)
		return fiber.NewError(fiber.StatusInternalServerError, "Invalid user configuration")
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

	log.Println("[CreatePoopLog] Poop log created successfully")
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message":     "Poop log successfully created",
		"xp_gained":   xpGained,
		"poop_log":    newLog,
		"achievement": achievementMessage,
		"total_xp":    user.XP,
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

	userID, _ := uuid.Parse("2f9f3c05-75b0-4935-9d89-f074715f5c19") // Replace later
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
