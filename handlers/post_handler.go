package handlers

import (
	"log"
	"strconv"
	"time"

	// Import database package for DB instance access
	"github.com/mohamadsolkhannawawi/article-backend/database"
	// Import models package for Post struct
	"github.com/mohamadsolkhannawawi/article-backend/models"

	// Import Fiber web framework
	"github.com/gofiber/fiber/v2"
	// Import UUID package
	"github.com/google/uuid"
)

// CreatePostRequest is the struct for parsing and validating the create post request body
type CreatePostRequest struct {
	Title            string   `json:"title" validate:"required,min=20"`  
	Content          string   `json:"content" validate:"required,min=200"` 
	Category         string   `json:"category" validate:"required,min=3"`
	Status           string   `json:"status" validate:"required,oneof=publish draft thrash"` 
	FeaturedImageURL string   `json:"featured_image_url" validate:"omitempty,url"` // URL allow empty or valid URL
	Tags             []string `json:"tags" validate:"omitempty,dive,min=1"` // "dive" for validating each tag
}

// CreatePost is the handler for the POST /api/posts endpoint
func CreatePost(c *fiber.Ctx) error {
	// 1. Parse and validate the request body
	req := new(CreatePostRequest)
	if err := c.BodyParser(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status": "error", "message": "Invalid request body", "error": err.Error(),
		})
	}

	// Use the validator initialized in auth_handler.go
	if err := validate.Struct(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status": "error", "message": "Validation failed", "error": err.Error(),
		})
	}

	// 2. Get Author ID from middleware
	// Convert from 'interface{}' to 'string', then parse to UUID
	authorIDString, ok := c.Locals("userID").(string)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status": "error", "message": "Invalid user data in token",
		})
	}
	authorID, err := uuid.Parse(authorIDString)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status": "error", "message": "Invalid user ID format",
		})
	}

	// 3. Logic for handling Tags
	var tags []*models.Tag // Use slice of pointers to models.Tag
	// Loop through each tag name sent from the frontend
	for _, tagName := range req.Tags {
		var tag models.Tag
		// Try to find the tag, if not found, create a new one (FirstOrCreate)
		// This is very efficient and prevents duplicate tags
		result := database.DB.Where("name = ?", tagName).FirstOrCreate(&tag, models.Tag{
			ID:   uuid.New(),
			Name: tagName,
		})

		if result.Error != nil {
			log.Println("Error finding/creating tag:", result.Error)
			// skip tag if error occurs and log the error
			continue
		}
		// Append tag as pointer to slice
		tags = append(tags, &tag)  // Perbaikan di sini: gunakan &tag untuk mendapatkan pointer
	}

	// 4. Create new Post instance
	newPost := models.Post{
		ID:               uuid.New(),
		Title:            req.Title,
		Content:          req.Content,
		Category:         req.Category,
		Status:           req.Status,
		FeaturedImageURL: req.FeaturedImageURL,
		AuthorID:         authorID,
		Tags:             tags, // GORM will automatically fill the 'post_tags' table
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
	}

	// 5. Save post to database
	if err := database.DB.Create(&newPost).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status": "error", "message": "Failed to create post", "error": err.Error(),
		})
	}

	// 6. Load Author and Tags relations for response
	// (By default GORM does not automatically load relations on Create)
	// We will load them manually to ensure the JSON response is complete.
	database.DB.Preload("Author").Preload("Tags").First(&newPost, newPost.ID)

	// 7. Return the newly created post
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"status":  "success",
		"message": "Post created successfully",
		"data":    newPost,
	})
}

// GetPosts is the handler for the GET /api/posts endpoint
func GetPosts(c *fiber.Ctx) error {
	// 1. Parse query parameters for pagination and filtering
	// Set defaults
	limit, _ := strconv.Atoi(c.Query("limit", "10"))
	offset, _ := strconv.Atoi(c.Query("offset", "0"))
	status := c.Query("status", "") // Get status filter

	var posts []models.Post
	var total int64

	// 2. Build the database query
	// We start with a base query and add to it
	query := database.DB.Model(&models.Post{}).
		Preload("Author"). // Load the related Author
		Preload("Tags")    // Load the related Tags

	// 3. Apply status filter if provided
	if status != "" {
		// This allows filtering for 'publish', 'draft', or 'thrash'
		query = query.Where("status = ?", status)
	}

	// 4. Get the total count of records matching the filter
	// We need this for pagination metadata
	if err := query.Count(&total).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status": "error", "message": "Failed to count posts", "error": err.Error(),
		})
	}
	
	// 5. Apply pagination and order, then find the posts
	err := query.
		Order("created_at DESC"). // Show newest posts first
		Limit(limit).
		Offset(offset).
		Find(&posts).Error

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status": "error", "message": "Failed to retrieve posts", "error": err.Error(),
		})
	}
	
	// 6. Return the response with data and metadata
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status": "success",
		"message": "Posts retrieved successfully",
		"data":    posts,
		"meta": fiber.Map{ // Pagination info for the frontend
			"total":  total,
			"limit":  limit,
			"offset": offset,
		},
	})
}
