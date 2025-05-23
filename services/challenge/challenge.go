package challenge

import (
	"log"
	"time"

	"github.com/karunapo/flush-wars-backend/db"
	"github.com/karunapo/flush-wars-backend/models"
)

func isChallengeCompleted(user models.User, challenge models.Challenge) bool {
	log.Printf("[IsChallengeCompleted] Checking completion for user %s on challenge %s (GoalType: %s, GoalValue: %d)",
		user.ID, challenge.Name, challenge.GoalType, challenge.GoalValue)

	switch challenge.GoalType {
	case "logs":
		var count int64
		if err := db.DB.Model(&models.PoopLog{}).
			Where("user_id = ?", user.ID).
			Count(&count).Error; err != nil {
			log.Printf("[IsChallengeCompleted] Error counting poop logs for user %s: %v", user.ID, err)
			return false
		}
		log.Printf("[IsChallengeCompleted] User %s has %d poop logs, goal is %d", user.ID, count, challenge.GoalValue)
		if int(count) >= challenge.GoalValue {
			log.Printf("[IsChallengeCompleted] User %s completed logs goal", user.ID)
			return true
		}
		return false

	case "streak":
		log.Printf("[IsChallengeCompleted] User %s current streak is %d, goal is %d", user.ID, user.CurrentStreak, challenge.GoalValue)
		if user.CurrentStreak >= challenge.GoalValue {
			log.Printf("[IsChallengeCompleted] User %s completed streak goal", user.ID)
			return true
		}
		return false

	case "xp":
		log.Printf("[IsChallengeCompleted] User %s XP is %d, goal is %d", user.ID, user.XP, challenge.GoalValue)
		if user.XP >= challenge.GoalValue {
			log.Printf("[IsChallengeCompleted] User %s completed XP goal", user.ID)
			return true
		}
		return false

	default:
		log.Printf("[IsChallengeCompleted] Unknown goal type: %s", challenge.GoalType)
		return false
	}
}

func CheckForNewlyCompletedChallenge(user models.User) (string, error) {
	var challenges []models.Challenge
	// Get active challenges
	if err := db.DB.Where("start_date <= NOW() AND end_date >= NOW()").Find(&challenges).Error; err != nil {
		log.Printf("[CheckForNewlyCompletedChallenge] Failed to fetch challenges: %v", err)
		return "", err
	}

	for _, ch := range challenges {
		// Check if user already completed this challenge
		var userCh models.UserChallenge
		err := db.DB.Where("user_id = ? AND challenge_id = ?", user.ID, ch.ID).First(&userCh).Error
		if err != nil || userCh.Completed {
			// Skip if not found or already completed
			continue
		}

		if isChallengeCompleted(user, ch) {
			// Mark challenge as complete in DB
			userCh.Completed = true
			now := time.Now()
			userCh.CompletedAt = &now
			userCh.UserID = user.ID
			userCh.ChallengeID = ch.ID
			if err := db.DB.Save(&userCh).Error; err != nil {
				log.Printf("[CheckForNewlyCompletedChallenge] Failed to mark challenge complete: %v", err)
				return "", err
			}
			log.Printf("[CheckForNewlyCompletedChallenge] User %s completed challenge %s", user.ID, ch.Name)
			return ch.Name, nil
		}
	}

	return "", nil
}
