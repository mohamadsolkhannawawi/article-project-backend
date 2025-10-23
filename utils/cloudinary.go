package utils

import (
	"context"
	"io"
	"log"

	"github.com/mohamadsolkhannawawi/article-backend/config"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
)

var cld *cloudinary.Cloudinary

func InitCloudinary() {

	log.Println("Initializing Cloudinary...")

	// Replace os.Getenv with AppConfig
	cloudName := config.AppConfig.CloudinaryCloudName
	apiKey := config.AppConfig.CloudinaryAPIKey
	apiSecret := config.AppConfig.CloudinaryAPISecret

	if cloudName == "" || apiKey == "" || apiSecret == "" {
		log.Println("ERROR: Cloudinary credentials are not set (or using default values). Cloudinary features will be disabled.")
        // REMOVE log.Fatal() to avoid a total crash
		return 
	}

	var err error
	cld, err = cloudinary.NewFromParams(cloudName, apiKey, apiSecret)
	if err != nil {
		log.Printf("ERROR: Failed to initialize Cloudinary: %v", err)
        // REMOVE log.Fatalf()
		return
	}

	log.Println("✓ Cloudinary initialized successfully")
}

func UploadToCloudinary(file io.Reader, folder string) (string, error) {
	if cld == nil {
		log.Println("ERROR: Cloudinary not initialized")
		return "", nil
	}

	ctx := context.Background()
	uploadParams := uploader.UploadParams{
		Folder: folder,
	}

	result, err := cld.Upload.Upload(ctx, file, uploadParams)
	if err != nil {
		return "", err
	}

	return result.SecureURL, nil
}