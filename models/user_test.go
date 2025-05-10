package models

import (
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestNewUserInitialization(t *testing.T) {
	id := uuid.New()
	now := time.Now()

	user := User{
		ID:        id,
		Email:     "test@example.com",
		Username:  "testuser",
		Password:  "securepassword",
		CreatedAt: now,
	}

	if user.ID != id {
		t.Errorf("expected ID %v, got %v", id, user.ID)
	}

	if user.Email != "test@example.com" {
		t.Errorf("expected Email test@example.com, got %v", user.Email)
	}

	if user.Username != "testuser" {
		t.Errorf("expected Username testuser, got %v", user.Username)
	}

	if user.Password != "securepassword" {
		t.Errorf("expected Password securepassword, got %v", user.Password)
	}

	if !user.CreatedAt.Equal(now) {
		t.Errorf("expected CreatedAt %v, got %v", now, user.CreatedAt)
	}
}
