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
		log.Fatal("ERROR: DATABASE_URL environment variable is not set!")
	}

	log.Println("Attempting to connect to database...")

	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})

	if err != nil {
		log.Fatalf("ERROR: Failed to connect to database: %v", err)
	}

	// Test connection
	sqlDB, err := DB.DB()
	if err != nil {
		log.Fatalf("ERROR: Failed to get database instance: %v", err)
	}

	err = sqlDB.Ping()
	if err != nil {
		log.Fatalf("ERROR: Failed to ping database: %v", err)
	}

	log.Println("✓ Successfully connected to database")
}
