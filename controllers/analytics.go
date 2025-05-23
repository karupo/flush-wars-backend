package controllers

import (
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/karunapo/flush-wars-backend/db"
	"github.com/karunapo/flush-wars-backend/models"
)

// AnalyticsResponse holds summary statistics for the analytics endpoint.
type AnalyticsResponse struct {
	TotalLogs   int               `json:"total_logs"`
	TotalXP     int               `json:"total_xp"`
	ByStoolType map[string]int    `json:"by_stool_type"`
	Trend       []DailyTrendEntry `json:"trend,omitempty"` // Only for trends
}

// DailyTrendEntry represents a single entry in a trend graph.
type DailyTrendEntry struct {
	Date     string `json:"date"`
	LogCount int    `json:"log_count"`
	XP       int    `json:"xp"`
}

// GetAnalyticsWeekly returns weekly poop log trends for the user.
func GetAnalyticsWeekly(c *fiber.Ctx) error {
	log.Println("[Analytics] Weekly analytics requested")
	return getAnalyticsSince(c, 7)
}

// GetAnalyticsMonthly returns monthly poop log trends for the user.
func GetAnalyticsMonthly(c *fiber.Ctx) error {
	log.Println("[Analytics] Monthly analytics requested")
	return getAnalyticsSince(c, 30)
}

// GetAnalyticsYearly returns yearly poop log trends for the user.
func GetAnalyticsYearly(c *fiber.Ctx) error {
	log.Println("[Analytics] Yearly analytics requested")
	return getAnalyticsSince(c, 365)
}

// GetAnalyticsTrends returns daily trends in poop activity.
func GetAnalyticsTrends(c *fiber.Ctx) error {
	log.Println("[Analytics] Trends requested (last 30 days)")

	userID := getUserID(c)
	start := time.Now().AddDate(0, 0, -30)

	var logs []models.PoopLog
	if err := db.DB.Where("user_id = ? AND timestamp >= ?", userID, start).
		Order("timestamp asc").
		Find(&logs).Error; err != nil {
		log.Printf("[Analytics] DB error fetching logs: %v", err)
		return fiber.NewError(fiber.StatusInternalServerError, "Could not fetch logs")
	}

	log.Printf("[Analytics] %d logs fetched", len(logs))

	trend := map[string]*DailyTrendEntry{}
	for _, logEntry := range logs {
		day := logEntry.Timestamp.Format("2006-01-02")
		if trend[day] == nil {
			trend[day] = &DailyTrendEntry{Date: day}
		}
		trend[day].LogCount++
		trend[day].XP += logEntry.XPGained
	}

	var trendData []DailyTrendEntry
	for date := start; !date.After(time.Now()); date = date.AddDate(0, 0, 1) {
		day := date.Format("2006-01-02")
		entry := trend[day]
		if entry == nil {
			trendData = append(trendData, DailyTrendEntry{Date: day})
		} else {
			trendData = append(trendData, *entry)
		}
	}

	log.Println("[Analytics] Trend data compiled")

	return c.JSON(AnalyticsResponse{
		Trend: trendData,
	})
}

func getAnalyticsSince(c *fiber.Ctx, days int) error {
	userID := getUserID(c)
	start := time.Now().AddDate(0, 0, -days)

	var logs []models.PoopLog
	if err := db.DB.Where("user_id = ? AND timestamp >= ?", userID, start).
		Find(&logs).Error; err != nil {
		log.Printf("[Analytics] Error fetching logs for last %d days: %v", days, err)
		return fiber.NewError(fiber.StatusInternalServerError, "Could not fetch logs")
	}

	log.Printf("[Analytics] %d logs fetched for last %d days", len(logs), days)

	byType := map[string]int{}
	totalXP := 0
	for _, l := range logs {
		byType[l.StoolType]++
		totalXP += l.XPGained
	}

	response := AnalyticsResponse{
		TotalLogs:   len(logs),
		TotalXP:     totalXP,
		ByStoolType: byType,
	}

	log.Printf("[Analytics] Response: %+v", response)
	return c.JSON(response)
}

func getUserID(c *fiber.Ctx) uuid.UUID {
	userIDVal := c.Locals("userID")
	userID, ok := userIDVal.(uuid.UUID)
	if !ok {
		log.Println("[Analytics] Failed to get user ID from context")
	}

	return userID
}
