package models

import (
	"time"

	"github.com/google/uuid"
)

// Achievement Table to log all new achievements
type Achievement struct {
	ID          uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primary_key"`
	UserID      uuid.UUID `gorm:"type:uuid;not null"`
	Name        string    `gorm:"not null"`
	Description string
	AchievedAt  time.Time
}
