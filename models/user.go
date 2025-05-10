package models

import (
	"time"

	"github.com/google/uuid"
)

// User model
type User struct {
	ID        uuid.UUID
	Email     string
	Username  string
	Password  string
	CreatedAt time.Time
}
