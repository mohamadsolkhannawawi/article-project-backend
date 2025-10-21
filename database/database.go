package database

import (
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// DB is the global database connection instance with GORM ORM package will be shared across the application
var DB *gorm.DB

// ConnectDB is a function to initialize the database connection
func ConnectDB() {
	var err error

	// Get DSN (Data Source Name) from environment variable
	dsn := os.Getenv("DATABASE_URL")

	// Open connection to database
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	log.Println("Successfully connected to database.")
}