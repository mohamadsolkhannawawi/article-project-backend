package database

import (
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func ConnectDB() {
	var err error

	dsn := os.Getenv("DATABASE_URL")

	if dsn == "" {
		log.Println("ERROR: DATABASE_URL environment variable is not set!")
		log.Println("Please set DATABASE_URL in Vercel Environment Variables")
		// Don't use log.Fatal() - let it continue so we can see the error
		return
	}

	log.Println("Attempting to connect to database...")

	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})

	if err != nil {
		log.Printf("ERROR: Failed to connect to database: %v", err)
		// Don't crash - just log the error
		return
	}

	// Test connection
	sqlDB, err := DB.DB()
	if err != nil {
		log.Printf("ERROR: Failed to get database instance: %v", err)
		return
	}

	err = sqlDB.Ping()
	if err != nil {
		log.Printf("ERROR: Failed to ping database: %v", err)
		return
	}

	log.Println("✓ Successfully connected to database")
}