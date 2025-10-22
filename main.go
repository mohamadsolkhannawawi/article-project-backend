package main

import (
	"log"
	"net/http"
	"os"

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
	// POST /api/register - Public route to register a new user
	api.Post("/register", handlers.RegisterUser)
	// POST /api/login - Public route to login an existing user
	api.Post("/login", handlers.LoginUser)

	// GET /api/posts - Public route to get all published posts with pagination
	api.Get("/posts", handlers.GetPosts)

	// GET /api/posts/my - Get posts created by authenticated user (Protected) - MUST BE BEFORE /:id route
	api.Get("/posts/my", middleware.AuthRequired(), handlers.GetMyPosts)

	// GET /api/posts/:id - Get a single post by ID (Public)
	api.Get("/posts/:id", handlers.GetPostByID)


	// --- Protected Routes (require authentication) ---
	// GET /api/profile - Protected route to get user profile
	api.Get("/profile", middleware.AuthRequired(), func(c *fiber.Ctx) error {
		// Get user data stored by middleware
		userID := c.Locals("userID")
		email := c.Locals("userEmail")
		fullName := c.Locals("userFullName")
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"status":  "success",
			"message": "Profile data retrieved successfully",
			"data": fiber.Map{
				"id":        userID,
				"email":     email,
				"full_name": fullName,
			},
		})
	})

	// GET /api/admin/posts - Get posts, including trashed (Protected)
	api.Get("/admin/posts", middleware.AuthRequired(), handlers.GetAdminPosts)

	// POST /api/upload - Upload an image (Protected)
	api.Post("/upload", middleware.AuthRequired(), handlers.UploadImage)

	// --- POST Routes (CRUD for Posts) ---
	/// POST /api/posts - Protected route to create a new post
	api.Post("/posts", middleware.AuthRequired(), handlers.CreatePost)

	// PUT /api/posts/:id - Update a post (Protected)
	api.Put("/posts/:id", middleware.AuthRequired(), handlers.UpdatePost)

	// DELETE /api/posts/:id - Soft delete a post (Protected)
	api.Delete("/posts/:id", middleware.AuthRequired(), handlers.DeletePost)

	// --- Add a simple root handler ---
    // Needed because vercel.json routes "/" to main.go
    app.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"message": "Welcome to KataGenzi API!",
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
	// Load .env (Vercel injects env vars, but this is good for local)
	err := godotenv.Load()
	if err != nil {
		log.Println("Note: .env file not found, relying on Vercel env vars.")
	}

	// Connect to DB
	database.ConnectDB()

	// Initialize Cloudinary
	utils.InitCloudinary()

    // Run migrations (only if needed - Vercel might reuse instances)
    // Consider moving migrations to a separate script/process for production
	runMigrations(database.DB)

	// --- Create and configure the Fiber app ---
	app = fiber.New()

	// CORS Middleware (Allow frontend origin in production)
    // Replace "*" with your frontend Vercel URL in production
	app.Use(func(c *fiber.Ctx) error {
		// TODO: Replace "*" with your frontend Vercel URL for production
        // Example: c.Set("Access-Control-Allow-Origin", "https://your-frontend-url.vercel.app")
		c.Set("Access-Control-Allow-Origin", "*") 
		c.Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Set("Access-Control-Allow-Headers", "Origin, Content-Type, Accept, Authorization")
		
		if c.Method() == "OPTIONS" {
			return c.SendStatus(fiber.StatusNoContent)
		}
		return c.Next()
	})

	// Setup all routes
	setupRoutes(app)
}

// Handler is the exported function Vercel will call
func Handler(w http.ResponseWriter, r *http.Request) {
    // Use adaptor.FiberApp to convert Fiber handler to net/http handler
	adaptor.FiberApp(app)(w, r)
}

// main function is NO LONGER USED by Vercel for serving requests,
// but it's useful for local development using `go run main.go` or `air`.
func main() {
    // We still initialize everything in init() for local dev
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}
	log.Printf("Starting LOCAL server on port %s...\n", port)
	log.Fatal(app.Listen(":" + port))
}