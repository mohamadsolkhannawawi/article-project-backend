package handler

import (
	"log"
	"net/http"
	"os"

	"github.com/mohamadsolkhannawawi/article-backend/database"
	"github.com/mohamadsolkhannawawi/article-backend/handlers"
	"github.com/mohamadsolkhannawawi/article-backend/middleware"
	"github.com/mohamadsolkhannawawi/article-backend/models"
	"github.com/mohamadsolkhannawawi/article-backend/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/adaptor"
	"gorm.io/gorm"
)

var app *fiber.App

func runMigrations(db *gorm.DB) {
	if db == nil {
		log.Println("WARNING: Cannot run migrations - database is nil")
		return
	}
	
	log.Println("Running Migrations...")
	err := db.AutoMigrate(&models.User{}, &models.Tag{}, &models.Post{})
	if err != nil {
		log.Printf("ERROR: Failed to migrate database: %v\n", err)
	} else {
		log.Println("✓ Database Migrated Successfully!")
	}
}

func setupRoutes(app *fiber.App) {
	// Root handler
	app.Get("/", func(c *fiber.Ctx) error {
		// Check if DB is connected
		dbStatus := "disconnected"
		if database.DB != nil {
			dbStatus = "connected"
		}
		
		return c.JSON(fiber.Map{
			"message": "Welcome to KataGenzi API!",
			"status":  "ok",
			"database": dbStatus,
		})
	})

	// API group
	api := app.Group("/api")

	// --- Public Auth Routes ---
	api.Post("/register", handlers.RegisterUser)
	api.Post("/login", handlers.LoginUser)

	// --- Public Post Routes ---
	api.Get("/posts", handlers.GetPosts)
	// ⭐ IMPORTANT: /posts/my MUST come BEFORE /posts/:id
	api.Get("/posts/my", middleware.AuthRequired(), handlers.GetMyPosts)
	api.Get("/posts/:id", handlers.GetPostByID)

	// --- Protected Post Routes ---
	api.Post("/posts", middleware.AuthRequired(), handlers.CreatePost)
	api.Put("/posts/:id", middleware.AuthRequired(), handlers.UpdatePost)
	api.Delete("/posts/:id", middleware.AuthRequired(), handlers.DeletePost)

	// --- Protected Admin Routes ---
	api.Get("/admin/posts", middleware.AuthRequired(), handlers.GetAdminPosts)
	api.Post("/upload", middleware.AuthRequired(), handlers.UploadImage)

	// --- Protected User Routes ---
	api.Get("/profile", middleware.AuthRequired(), func(c *fiber.Ctx) error {
		userID := c.Locals("userID")
		email := c.Locals("userEmail")
		fullName := c.Locals("userFullName")
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"status":  "success",
			"message": "Profile data",
			"data":    fiber.Map{"id": userID, "email": email, "full_name": fullName},
		})
	})

	// 404 Handler for API
	api.Use(func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status":  "error",
			"message": "API endpoint not found",
		})
	})
}

func init() {
	log.Println("========================================")
	log.Println("=== VERCEL INIT START ===")
	log.Println("========================================")

	// Check all env vars
	dbURL := os.Getenv("DATABASE_URL")
	cloudName := os.Getenv("CLOUDINARY_CLOUD_NAME")
	cloudKey := os.Getenv("CLOUDINARY_API_KEY")
	cloudSecret := os.Getenv("CLOUDINARY_API_SECRET")
	jwtSecret := os.Getenv("JWT_SECRET")

	log.Printf("Environment Variables Check:")
	log.Printf("  DATABASE_URL: %v", dbURL != "")
	log.Printf("  CLOUDINARY_CLOUD_NAME: %v", cloudName != "")
	log.Printf("  CLOUDINARY_API_KEY: %v", cloudKey != "")
	log.Printf("  CLOUDINARY_API_SECRET: %v", cloudSecret != "")
	log.Printf("  JWT_SECRET: %v", jwtSecret != "")

	// Connect to database (won't crash if fails now)
	log.Println("Connecting to database...")
	database.ConnectDB()

	// Initialize Cloudinary
	log.Println("Initializing Cloudinary...")
	utils.InitCloudinary()

	// Run migrations (only if DB is connected)
	runMigrations(database.DB)

	// Create Fiber app
	app = fiber.New(fiber.Config{
		DisableStartupMessage: true,
	})

	// CORS
	app.Use(func(c *fiber.Ctx) error {
		c.Set("Access-Control-Allow-Origin", "*")
		c.Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Set("Access-Control-Allow-Headers", "Origin, Content-Type, Accept, Authorization")

		if c.Method() == "OPTIONS" {
			return c.SendStatus(fiber.StatusNoContent)
		}
		return c.Next()
	})

	setupRoutes(app)

	log.Println("========================================")
	log.Println("=== VERCEL INIT COMPLETE ===")
	log.Println("========================================")
}

// Handler adalah entry point untuk Vercel
func Handler(w http.ResponseWriter, r *http.Request) {
	log.Printf("Request: %s %s", r.Method, r.URL.Path)
	adaptor.FiberApp(app)(w, r)
}