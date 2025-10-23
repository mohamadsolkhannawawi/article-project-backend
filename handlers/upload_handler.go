package handlers

import (
	"github.com/mohamadsolkhannawawi/article-backend/utils"

	"github.com/gofiber/fiber/v2"
)

// UploadImage is the handler for POST /api/upload
func UploadImage(c *fiber.Ctx) error {
	// 1. Get the file from the form-data
	// "image" is the key name we expect from the frontend
	file, err := c.FormFile("image")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status": "error", "message": "Failed to get image from form", "error": err.Error(),
		})
	}

	// 2. Open the file
	src, err := file.Open()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status": "error", "message": "Failed to open file", "error": err.Error(),
		})
	}
	defer src.Close()

	// 3. Upload the file to Cloudinary
	// We pass the file source (src) and a folder name
	secureURL, err := utils.UploadToCloudinary(src, "article-project")
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status": "error", "message": "Failed to upload file to Cloudinary", "error": err.Error(),
		})
	}

	// 4. Return the secure URL
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "Image uploaded successfully",
		"data": fiber.Map{
			"url": secureURL,
		},
	})
}
