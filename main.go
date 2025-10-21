package main

import (
	"log"

	// Import package database, handlers, middleware, models
	"github.com/mohamadsolkhannawawi/article-backend/database"
	"github.com/mohamadsolkhannawawi/article-backend/handlers"
	"github.com/mohamadsolkhannawawi/article-backend/middleware"
	"github.com/mohamadsolkhannawawi/article-backend/models"

	// Import Fiber web framework to start the server later
	"github.com/gofiber/fiber/v2"
	// Import godotenv to load .env file
	"github.com/joho/godotenv"
	// Import GORM for database migrations
	"gorm.io/gorm"
)

// runMigrations is a helper function to run database migrations
func runMigrations(db *gorm.DB) {
	log.Println("Running Migrations...")
	err := db.AutoMigrate(&models.User{}, &models.Tag{}, &models.Post{})
	if err != nil {
		log.Fatal("Failed to migrate database:", err)
	}
	log.Println("Database Migrated Successfully!")
}

// setupRoutes is a helper function to register all routes
func setupRoutes(app *fiber.App) {
	// Route Group for /api
	api := app.Group("/api")

	// --- Public Routes ---
	// Auth Routes
	// POST /api/register
	api.Post("/register", handlers.RegisterUser)
	// POST /api/login
	api.Post("/login", handlers.LoginUser)

	// --- Protected Routes (require authentication) ---
	api.Get("/profile", middleware.AuthRequired(), func(c *fiber.Ctx) error {
		// Get user data stored by middleware
		userID := c.Locals("userID")
		email := c.Locals("userEmail")
		fullName := c.Locals("userFullName")

		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"status": "success",
			"message": "Profile data retrieved successfully",
			"data": fiber.Map{
				"id":        userID,
				"email":     email,
				"full_name": fullName,
			},
		})
	})

	// --- POST Routes (CRUD for Posts) ---
	/// POST /api/posts - Create a new post
	api.Post("/posts", middleware.AuthRequired(), handlers.CreatePost) // Protected route
}

func main() {
	// Load .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file:", err)
	}

	// Connect to the database
	database.ConnectDB()

	// Run Migrations
	runMigrations(database.DB)

	// Create a new Fiber app
	app := fiber.New()

	// Define simple route for testing
	app.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"message": "Hello, World!",
		})
	})

	// --- CORS Middleware ---
	app.Use(func(c *fiber.Ctx) error {
		c.Set("Access-Control-Allow-Origin", "*")
		c.Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Set("Access-Control-Allow-Headers", "Origin, Content-Type, Accept, Authorization")
		
		// Handle preflight OPTIONS requests
		if c.Method() == "OPTIONS" {
			return c.SendStatus(fiber.StatusNoContent)
		}
		return c.Next()
	})

	// Register routes
	setupRoutes(app)

	// Handle 404 - Not Found
	app.Use(func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status":  "error",
			"message": "Endpoint not found",
		})
	})

	// Start the server
	log.Println("Server is running on http://localhost:3000")
	log.Fatal(app.Listen(":3000"))
}