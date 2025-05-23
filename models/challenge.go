package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Challenge struct {
	ID          uuid.UUID      `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	Name        string         `json:"name"`
	Description string         `json:"description"`
	GoalType    string         `json:"goal_type"`  // e.g., "logs", "streak", "xp"
	GoalValue   int            `json:"goal_value"` // e.g., 7 logs, 3-day streak, 100 XP
	StartDate   time.Time      `json:"start_date"`
	EndDate     time.Time      `json:"end_date"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}
