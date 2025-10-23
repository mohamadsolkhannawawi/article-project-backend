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
	log.Println("Running Migrations...")
	err := db.AutoMigrate(&models.User{}, &models.Tag{}, &models.Post{})
	if err != nil {
		log.Printf("ERROR: Failed to migrate database: %v\n", err)
	} else {
		log.Println("Database Migrated Successfully!")
	}
}

func setupRoutes(app *fiber.App) {
	api := app.Group("/api")

	// --- Public Auth Routes ---
	api.Post("/register", handlers.RegisterUser)
	api.Post("/login", handlers.LoginUser)

	// --- Public Post Routes ---
	api.Get("/posts", handlers.GetPosts)
	// ⭐ IMPORTANT: /posts/my MUST come BEFORE /posts/:id
	// Otherwise /posts/my will be caught by /posts/:id route (my treated as ID parameter)
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

	app.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"message": "Welcome to KataGenzi API!",
			"status":  "ok",
		})
	})

	api.Use(func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status":  "error",
			"message": "API endpoint not found",
		})
	})
}

func init() {
	log.Println("Initializing application...")

	// Note: .env tidak akan dibaca di Vercel, gunakan Environment Variables
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Println("WARNING: DATABASE_URL is not set!")
	} else {
		log.Println("✓ DATABASE_URL is set")
	}

	// Connect to database
	database.ConnectDB()

	// Initialize Cloudinary
	utils.InitCloudinary()

	// Run migrations
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

	// Setup routes
	setupRoutes(app)

	log.Println("✓ Application initialized successfully")
}

// Handler adalah entry point untuk Vercel
func Handler(w http.ResponseWriter, r *http.Request) {
	adaptor.FiberApp(app)(w, r)
}
