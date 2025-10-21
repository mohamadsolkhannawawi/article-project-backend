package main

import (
	"log"

	// Import package database dan models
	"github.com/mohamadsolkhannawawi/article-backend/database"
	"github.com/mohamadsolkhannawawi/article-backend/models"

	"github.com/joho/godotenv"
)

// Init will be called before main()
func init() {
	// Load .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	
	// Connect to Database
	database.ConnectDB()
}

func main() {
	log.Println("Running Migrations...")

	// Run AutoMigrate to create/update tables based on our models
	// GORM will automatically create tables, foreign keys, and constraints
	// based on our Model structs.
	// Important order: User and Tag must exist before Post.
	err := database.DB.AutoMigrate(&models.User{}, &models.Tag{}, &models.Post{})
	if err != nil {
		log.Fatal("Failed to migrate database:", err)
	}

	log.Println("Database Migrated Successfully!")

	// Here we would start the server (e.g., using Gin or another framework)
	log.Println("Setup complete. Server will be started here...")
}