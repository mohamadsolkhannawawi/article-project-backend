package middleware

import (
	"log"
	"strings"

	"github.com/mohamadsolkhannawawi/article-backend/config"

	// Import package handlers to access JwtCustomClaims struct
	"github.com/mohamadsolkhannawawi/article-backend/handlers"

	// Import Fiber web framework
	"github.com/gofiber/fiber/v2"
	// Import JWT package
	"github.com/golang-jwt/jwt/v5"
)

// AuthRequired is a middleware to protect routes that require authentication
func AuthRequired() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// 1. Get Authorization header
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"status":  "error",
				"message": "Missing authorization header",
			})
		}

		// 2. Token is usually sent in the format "Bearer <token>"
		// We need to separate "Bearer" from the token itself
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"status":  "error",
				"message": "Invalid authorization header format",
			})
		}

		tokenString := parts[1]

		// 3. Get JWT_SECRET from .env
		// jwtSecret := os.Getenv("JWT_SECRET")
		jwtSecret := config.AppConfig.JWTSecret
		if jwtSecret == "" {
			log.Println("Warning: JWT_SECRET is not set")
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"status":  "error",
				"message": "Server configuration error",
			})
		}

		// 4. Parse and validate token
		claims := &handlers.JwtCustomClaims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			// Make sure the signing method is HS256
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fiber.ErrUnauthorized
			}
			return []byte(jwtSecret), nil
		})

		if err != nil || !token.Valid {
			// This can happen if the token is expired or invalid
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"status":  "error",
				"message": "Invalid or expired token",
				"error":   err.Error(),
			})
		}

		// 5. Token valid!
		// We store user info from the token into Fiber's context
		// so it can be accessed by subsequent handlers.
		c.Locals("userID", claims.UserID)
		c.Locals("userEmail", claims.Email)
		c.Locals("userFullName", claims.FullName)

		// Proceed to the next handler (endpoint)
		return c.Next()
	}
}
