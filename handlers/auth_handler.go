package handlers

import (
	"log"
	"os"
	"time"

	// Import package database dan models
	"github.com/mohamadsolkhannawawi/article-backend/database"
	"github.com/mohamadsolkhannawawi/article-backend/models"

	// Import package validator for input validation
	"github.com/go-playground/validator/v10"
	// Import Fiber web framework
	"github.com/gofiber/fiber/v2"
	// Import package for UUID generation
	"github.com/google/uuid"
	// Import bcrypt for password hashing
	"golang.org/x/crypto/bcrypt"
	// Import JWT for token generation
	"github.com/golang-jwt/jwt/v5"
	// Import GORM for database operations
	"gorm.io/gorm"
)

// Initialize validator instance
var validate = validator.New()

// RegisterRequest is struct for parsing and validating user registration request body
type RegisterRequest struct {
	FullName string `json:"full_name" validate:"required,min=3"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
}

// hashPassword is a helper function to hash passwords using bcrypt
func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14) // 14 is the cost factor
	return string(bytes), err
}

// RegisterUser is a handler for the POST /api/register endpoint
func RegisterUser(c *fiber.Ctx) error {
	// 1. Parse request body to RegisterRequest struct
	req := new(RegisterRequest)
	if err := c.BodyParser(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid request body",
			"error":   err.Error(),
		})
	}

	// 2. Validate input using validator
	if err := validate.Struct(req); err != nil {
		// Return more descriptive validation errors
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Validation failed",
			"error":   err.Error(),
		})
	}

	// 3. Hash password
	hashedPassword, err := hashPassword(req.Password)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to hash password",
			"error":   err.Error(),
		})
	}

	// 4. Create a new User object
	newUser := models.User{
		ID:           uuid.New(), // Generate new UUID
		FullName:     req.FullName,
		Email:        req.Email,
		PasswordHash: hashedPassword,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	// 5. Save user to database
	// We use Create().Error to handle potential errors (e.g., duplicate email)
	if err := database.DB.Create(&newUser).Error; err != nil {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{
			"status":  "error",
			"message": "Email already exists",
			"error":   err.Error(),
		})
	}

	// 6. Return success response
	// We do not return the password hash
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"status":  "success",
		"message": "User registered successfully",
		"data": fiber.Map{
			"id":         newUser.ID,
			"full_name":  newUser.FullName,
			"email":      newUser.Email,
			"created_at": newUser.CreatedAt,
		},
	})
}

// LoginRequest is struct for parsing and validating user login request body
type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

// JwtCustomClaims defines the custom claims for JWT
type JwtCustomClaims struct {
	UserID   string `json:"user_id"`
	FullName string `json:"full_name"`
	Email    string `json:"email"`
	jwt.RegisteredClaims
}

// checkPasswordHash compares the raw password with the hash
func checkPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil // true if match, false if not
}

// generateJWT creates a new JWT token for the user
func generateJWT(user *models.User) (string, error) {
	// Get JWT_SECRET from .env
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		log.Println("Warning: JWT_SECRET environment variable not set. Using default.")
		jwtSecret = "default_secret_please_change" // Fallback (not safe for production)
	}

	// Set Claims is data where will be stored in the token
	claims := &JwtCustomClaims{
		UserID:   user.ID.String(),
		FullName: user.FullName,
		Email:    user.Email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 72)), // Token berlaku 72 jam
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	// Create token with claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign the token with your secret
	t, err := token.SignedString([]byte(jwtSecret))
	if err != nil {
		return "", err
	}

	return t, nil
}


// LoginUser is a handler for the POST /api/login endpoint
func LoginUser(c *fiber.Ctx) error {
	// 1. Parse request body to LoginRequest struct
	req := new(LoginRequest)
	if err := c.BodyParser(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status": "error", "message": "Invalid request body", "error": err.Error(),
		})
	}

	// 2. Validate input
	if err := validate.Struct(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status": "error", "message": "Validation failed", "error": err.Error(),
		})
	}

	var user models.User

	// 3. Find user by email
	// We use First() to get a single record.
	err := database.DB.Where("email = ?", req.Email).First(&user).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			// Email not found
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"status": "error", "message": "Invalid credentials",
			})
		}
		// Database error and other errors
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status": "error", "message": "Database error", "error": err.Error(),
		})
	}

	// 4. Check password for hash match
	if !checkPasswordHash(req.Password, user.PasswordHash) {
		// Password is incorrect
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status": "error", "message": "Invalid credentials",
		})
	}

	// 5. Create JWT
	token, err := generateJWT(&user)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status": "error", "message": "Failed to generate token", "error": err.Error(),
		})
	}

	// 6. Return success response with token
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "Login successful",
		"data": fiber.Map{
			"token": token,
		},
	})
}