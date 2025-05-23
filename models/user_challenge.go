package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// UserChallenge links a user to a challenge and tracks their progress and completion status.
type UserChallenge struct {
	ID          uuid.UUID      `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	UserID      uuid.UUID      `gorm:"type:uuid;index" json:"user_id"`
	ChallengeID uuid.UUID      `gorm:"type:uuid;index" json:"challenge_id"`
	JoinedAt    time.Time      `json:"joined_at"`
	Progress    int            `json:"progress"`
	Completed   bool           `json:"completed"`
	CompletedAt *time.Time     `json:"completed_at,omitempty"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}
