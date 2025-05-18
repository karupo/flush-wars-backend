package models

import (
	"time"

	"github.com/google/uuid"
)

// PoopLog model
type PoopLog struct {
	ID        uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	UserID    uuid.UUID `gorm:"type:uuid;not null"`
	StoolType string    `gorm:"not null"`
	Timestamp time.Time `gorm:"not null"`
	Notes     string    `gorm:"type:text"`
	XPGained  int       `gorm:"not null"`
	User      User      `gorm:"foreignKey:ID"`
}
