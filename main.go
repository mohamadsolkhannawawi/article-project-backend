package handler

import (
	"log"
	"net/http"

	// Import package database, handlers, middleware, models, utils
	"github.com/mohamadsolkhannawawi/article-backend/database"
	"github.com/mohamadsolkhannawawi/article-backend/handlers"
	"github.com/mohamadsolkhannawawi/article-backend/middleware"
	"github.com/mohamadsolkhannawawi/article-backend/models"
	"github.com/mohamadsolkhannawawi/article-backend/utils"

	// Import Fiber web framework to start the server later
	"github.com/gofiber/fiber/v2"
	// Import adaptor for http compatibility
	"github.com/gofiber/fiber/v2/middleware/adaptor"
	// Import godotenv to load .env file
	"github.com/joho/godotenv"
	// Import GORM for database migrations
	"gorm.io/gorm"
)

var app *fiber.App // Make app a global variable

func runMigrations(db *gorm.DB) {
	log.Println("Running Migrations...")
	err := db.AutoMigrate(&models.User{}, &models.Tag{}, &models.Post{})
	if err != nil {
		// Log fatal might stop the serverless function cold, maybe just log error?
		log.Printf("ERROR: Failed to migrate database: %v\n", err)
	} else {
		log.Println("Database Migrated Successfully!")
	}
}

func setupRoutes(app *fiber.App) {
	api := app.Group("/api")

	// --- Auth Routes (Public) ---
	api.Post("/register", handlers.RegisterUser)
	api.Post("/login", handlers.LoginUser)

	// --- Posts Routes (Public) ---
	api.Get("/posts", handlers.GetPosts)
	api.Get("/posts/:id", handlers.GetPostByID)

	// --- Protected Post Routes ---
	api.Post("/posts", middleware.AuthRequired(), handlers.CreatePost)
	api.Put("/posts/:id", middleware.AuthRequired(), handlers.UpdatePost)
	api.Delete("/posts/:id", middleware.AuthRequired(), handlers.DeletePost)

	// --- Protected Admin Routes ---
	api.Get("/profile", middleware.AuthRequired(), func(c *fiber.Ctx) error {
        userID := c.Locals("userID")
		email := c.Locals("userEmail")
		fullName := c.Locals("userFullName")
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"status": "success", "message": "Profile data",
			"data": fiber.Map{"id": userID, "email": email, "full_name": fullName},
		})
     })
	api.Get("/admin/posts", middleware.AuthRequired(), handlers.GetAdminPosts)

	// --- Upload Route (Protected) ---
	api.Post("/upload", middleware.AuthRequired(), handlers.UploadImage)

    // --- Add a simple root handler ---
    app.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"message": "Welcome to KataMerpati API!",
		})
	})

    // Handle 404 - Not Found for API routes specifically
	api.Use(func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status":  "error",
			"message": "API endpoint not found",
		})
	})
}

// init function runs once when the serverless function starts
func init() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Note: .env file not found, relying on Vercel env vars.")
	}

	database.ConnectDB()
	utils.InitCloudinary()
	runMigrations(database.DB) // Run migrations on init

	app = fiber.New()

	// CORS Middleware
	app.Use(func(c *fiber.Ctx) error {
		// TODO: Replace "*" with your frontend Vercel URL for production
		c.Set("Access-Control-Allow-Origin", "*")
		c.Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Set("Access-Control-Allow-Headers", "Origin, Content-Type, Accept, Authorization")
		
		if c.Method() == "OPTIONS" {
			return c.SendStatus(fiber.StatusNoContent)
		}
		return c.Next()
	})

	setupRoutes(app)
}

// Handler is the exported function Vercel will call
func Handler(w http.ResponseWriter, r *http.Request) {
	adaptor.FiberApp(app)(w, r)
}