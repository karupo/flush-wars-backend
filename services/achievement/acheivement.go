package achievement

import (
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/karunapo/flush-wars-backend/db"
	"github.com/karunapo/flush-wars-backend/models"
	"github.com/karunapo/flush-wars-backend/services/xp"
)

// CheckAndAwardAchievements evaluates user progress and awards new achievements if criteria are met.
// Returns a message indicating the result.
func CheckAndAwardAchievements(userID uuid.UUID) string {
	var messages []string

	// Count poop logs as before
	var count int64
	if err := db.DB.Model(&models.PoopLog{}).Where("user_id = ?", userID).Count(&count).Error; err != nil {
		log.Printf("[Achievement] Failed to count poop logs for user %s: %v", userID, err)
		return "Error checking achievements"
	}

	// Poop log milestones (unchanged)
	switch count {
	case 1:
		messages = append(messages, awardAchievement(userID, "First Log", "You logged your first poop!"))
	case 5:
		messages = append(messages, awardAchievement(userID, "Fifth Log", "You logged your fifth poop!"))
	}

	// Get user and current level
	var user models.User
	if err := db.DB.First(&user, "id = ?", userID).Error; err != nil {
		log.Printf("[Achievement] Failed to fetch user %s: %v", userID, err)
		return "Error fetching user for level achievements"
	}
	level := xp.CalculateLevel(user.XP)

	// Compose achievement name for the level
	levelAchievementName := fmt.Sprintf("Level %d Reached", level)
	levelAchievementDesc := fmt.Sprintf("You reached level %d!", level)

	// Award level achievement only if not awarded yet
	messages = append(messages, awardAchievement(userID, levelAchievementName, levelAchievementDesc))

	// Combine messages and return
	return combineMessages(messages)
}

func combineMessages(msgs []string) string {
	unique := make(map[string]bool)
	var result []string
	for _, msg := range msgs {
		if !unique[msg] {
			result = append(result, msg)
			unique[msg] = true
		}
	}
	return fmt.Sprintf("Achievements: %v", result)
}

// awardAchievement creates a new achievement entry for the user if it does not already exist.
func awardAchievement(userID uuid.UUID, name, description string) string {
	var exists int64
	if err := db.DB.Model(&models.Achievement{}).
		Where("user_id = ? AND name = ?", userID, name).
		Count(&exists).Error; err != nil {
		log.Printf("[Achievement] Failed to check if achievement exists for user %s: %v", userID, err)
		return "Error awarding achievement"
	}

	if exists > 0 {
		log.Printf("[Achievement] User %s already has achievement: %s", userID, name)
		return "Achievement already awarded"
	}

	newAchievement := models.Achievement{
		UserID:      userID,
		Name:        name,
		Description: description,
		AchievedAt:  time.Now(),
	}

	if err := db.DB.Create(&newAchievement).Error; err != nil {
		log.Printf("[Achievement] Failed to create achievement for user %s: %v", userID, err)
		return "Failed to award achievement"
	}

	log.Printf("[Achievement] Awarded '%s' to user %s", name, userID)
	return description
}

// GetTitle returns a rank-based title for leaderboard display.
func GetTitle(rank int) string {
	switch rank {
	case 0:
		return "Lord of the Bowels"
	case 1:
		return "Duke of Digestion"
	case 2:
		return "Knight of the Flush"
	case 3:
		return "Guardian of the Gut"
	case 4:
		return "Prince of Poop"
	default:
		return "Warrior of the Wipe"
	}
}
