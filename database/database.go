package database

import (
	"log"

	"github.com/mohamadsolkhannawawi/article-backend/config"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func ConnectDB() {
	var err error

	dsn := config.AppConfig.DatabaseURL

	log.Println("Connecting to database...")

	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Warn),
	})

	if err != nil {
		log.Printf("ERROR: Failed to connect to database: %v", err)
		return
	}

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

	log.Println("✓ Database connected successfully")
}