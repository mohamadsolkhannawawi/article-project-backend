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

	log.Printf("DEBUG: Attempting to connect with DATABASE_URL = [%s]\n", dsn)

	// Open connection to database
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Printf("ERROR: Failed to connect to database using DSN [%s]: %v\n", dsn, err)
		panic(err)
	}

	log.Println("Successfully connected to database.")
}