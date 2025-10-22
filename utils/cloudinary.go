package utils

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
)

var Cld *cloudinary.Cloudinary

// InitCloudinary initializes the Cloudinary service
func InitCloudinary() {
	// Get credentials from .env
	cloudName := os.Getenv("CLOUDINARY_CLOUD_NAME")
	apiKey := os.Getenv("CLOUDINARY_API_KEY")
	apiSecret := os.Getenv("CLOUDINARY_API_SECRET")

	if cloudName == "" || apiKey == "" || apiSecret == "" {
		log.Fatal("Cloudinary credentials are not set in .env")
	}

	var err error
	Cld, err = cloudinary.NewFromParams(cloudName, apiKey, apiSecret)
	if err != nil {
		log.Fatalf("Failed to initialize Cloudinary: %v", err)
	}

	log.Println("Cloudinary service initialized successfully.")
}

// UploadToCloudinary uploads a file to Cloudinary
// 'file' can be a path (string) or file content (io.Reader)
func UploadToCloudinary(file interface{}, folder string) (string, error) {
	// Create a context with a timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Upload parameters
	uploadParams := uploader.UploadParams{
		Folder: folder,
	}

	// Perform the upload
	uploadResult, err := Cld.Upload.Upload(ctx, file, uploadParams)
	if err != nil {
		log.Println("Failed to upload file:", err)
		return "", err
	}

	// Return the secure URL
	return uploadResult.SecureURL, nil
}