// Package db contains the db connection and call for the application.
package db

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// DB is the global database connection pool.
var DB *gorm.DB

var openDB = func(dsn string) (*gorm.DB, error) {
	return gorm.Open(postgres.Open(dsn), &gorm.Config{})
}

// Init - Connect to DB
var Init = func(loadEnv bool) error {
	if loadEnv {
		if err := godotenv.Load(); err != nil {
			return fmt.Errorf("error loading .env file: %w", err)
		}
	}

	// Retrieve environment variables
	host := os.Getenv("DB_HOST")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")
	port := os.Getenv("DB_PORT")

	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		host, user, password, dbname, port,
	)

	db, err := openDB(dsn)
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}

	DB = db
	log.Println("Connected to PostgreSQL database")
	return nil
}
