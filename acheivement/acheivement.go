package acheivement

import (
	"time"

	"github.com/google/uuid"
	"github.com/karunapo/flush-wars-backend/db"
	"github.com/karunapo/flush-wars-backend/models"
)

// CheckAndAwardAchievements Example: Award "First Log" if user has only one log
func CheckAndAwardAchievements(userID uuid.UUID) string {
	var count int64
	db.DB.Model(&models.PoopLog{}).Where("user_id = ?", userID).Count(&count)
	if count == 1 {
		return createAchievement(userID, "First Log", "You logged your first poop!")
	}

	if count == 5 {
		return createAchievement(userID, "Fifth Log", "You logged your fifth poop!")
	}
	// Add more achievement rules below
	return "No new acheivements yet"
}

func createAchievement(userID uuid.UUID, name string, description string) string {
	var exists int64
	db.DB.Model(&models.Achievement{}).Where("user_id = ? AND name = ?", userID, name).Count(&exists)
	if exists == 0 {
		db.DB.Create(&models.Achievement{
			UserID:      userID,
			Name:        name,
			Description: description,
			AchievedAt:  time.Now(),
		})
	}

	return description
}
