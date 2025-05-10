package db

import (
	"fmt"
	"os"
	"testing"

	"gorm.io/gorm"
)

// mockOpenDB is a mock function for gorm.Open to control database connection during tests.
var mockOpenDB = func(dsn string) (*gorm.DB, error) {
	// Simulate a successful connection for the expected DSN
	if dsn == fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		"localhost", "testuser", "testpassword", "testdb", "5432",
	) {
		return &gorm.DB{}, nil // Simulate a successful connection
	}
	return nil, fmt.Errorf("failed to connect to mock database with dsn: %s", dsn)
}

func TestInit_Success(t *testing.T) {
	// Backup the original openDB and restore it after the test
	originalOpenDB := openDB
	defer func() { openDB = originalOpenDB }()

	// Use the mock openDB for this test
	openDB = mockOpenDB

	// Set necessary environment variables
	os.Setenv("DB_HOST", "localhost")
	os.Setenv("DB_USER", "testuser")
	os.Setenv("DB_PASSWORD", "testpassword")
	os.Setenv("DB_NAME", "testdb")
	os.Setenv("DB_PORT", "5432")
	defer func() {
		os.Unsetenv("DB_HOST")
		os.Unsetenv("DB_USER")
		os.Unsetenv("DB_PASSWORD")
		os.Unsetenv("DB_NAME")
		os.Unsetenv("DB_PORT")
	}()

	err := Init(false)
	if err != nil {
		t.Errorf("Init failed: %v", err)
	}
	if DB == nil {
		t.Errorf("DB is nil after successful initialization")
	}
}
