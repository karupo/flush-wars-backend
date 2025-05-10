package db

import (
	"os"
	"strings"
	"testing"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func TestInit_DBMock(t *testing.T) {
	// Set env vars manually
	os.Setenv("DB_HOST", "localhost")
	os.Setenv("DB_USER", "testuser")
	os.Setenv("DB_PASSWORD", "testpass")
	os.Setenv("DB_NAME", "testdb")
	os.Setenv("DB_PORT", "5432")

	// Mock the DB connection
	mockDB := &gorm.DB{}
	openDB = func(dsn string) (*gorm.DB, error) {
		// Optional: assert the dsn contains expected values
		if !strings.Contains(dsn, "localhost") {
			t.Errorf("Unexpected DSN: %s", dsn)
		}
		return mockDB, nil
	}
	defer func() {
		openDB = func(dsn string) (*gorm.DB, error) {
			return gorm.Open(postgres.Open(dsn), &gorm.Config{})
		}
	}() // restore after test

	err := Init(false)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if DB != mockDB {
		t.Error("Expected DB to match mockDB")
	}
}
