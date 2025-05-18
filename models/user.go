// models/user.go
// Package models contains the data models for the application.
package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// User represents a user in the application.
type User struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey"`
	Email     string    `gorm:"unique;not null"`
	Username  string    `gorm:"unique;not null"`
	Password  string    `gorm:"not null"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

// BeforeCreate is a GORM hook that generates a UUID before inserting a new record.
func (u *User) BeforeCreate(_ *gorm.DB) (err error) {
	u.ID = uuid.New()
	return
}
