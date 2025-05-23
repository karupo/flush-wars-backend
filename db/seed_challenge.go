package db

import (
	"errors"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/karunapo/flush-wars-backend/models"
	"gorm.io/gorm"
)

func SeedChallenges() {
	sampleChallenges := []models.Challenge{
		{
			ID:          uuid.New(),
			Name:        "3-Day Streak",
			Description: "Log your poop for 3 consecutive days.",
			GoalType:    "streak",
			GoalValue:   3,
			StartDate:   time.Now(),
			EndDate:     time.Now().AddDate(0, 1, 0), // 1 month duration
		},
		{
			ID:          uuid.New(),
			Name:        "7 Logs in a Week",
			Description: "Log your poop 7 times within one week.",
			GoalType:    "logs",
			GoalValue:   7,
			StartDate:   time.Now(),
			EndDate:     time.Now().AddDate(0, 0, 7), // 7 days
		},
		{
			ID:          uuid.New(),
			Name:        "100 XP Milestone",
			Description: "Earn 100 XP through poop logging.",
			GoalType:    "xp",
			GoalValue:   100,
			StartDate:   time.Now(),
			EndDate:     time.Now().AddDate(0, 1, 0),
		},
		{
			ID:          uuid.New(),
			Name:        "5-Day Healthy Streak",
			Description: "Maintain a 5-day streak of healthy stool logs.",
			GoalType:    "streak",
			GoalValue:   5,
			StartDate:   time.Now(),
			EndDate:     time.Now().AddDate(0, 1, 0),
		},
	}

	for _, challenge := range sampleChallenges {
		var existing models.Challenge
		err := DB.Where("name = ?", challenge.Name).First(&existing).Error
		if err == nil {
			// Challenge with this name exists, skip insert
			log.Printf("[SeedChallenges] Challenge already exists: %s", challenge.Name)
			continue
		} else if !errors.Is(err, gorm.ErrRecordNotFound) {
			// Some other error occurred
			log.Printf("[SeedChallenges] Error checking challenge %s: %v", challenge.Name, err)
			continue
		}

		// Insert new challenge
		if err := DB.Create(challenge).Error; err != nil {
			log.Printf("[SeedChallenges] Failed to insert challenge %s: %v", challenge.Name, err)
		} else {
			log.Printf("[SeedChallenges] Challenge inserted: %s", challenge.Name)
		}
	}

}
