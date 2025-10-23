package config

import (
	"log"
	"os"
)

type Config struct {
	DatabaseURL         string
	CloudinaryCloudName string
	CloudinaryAPIKey    string
	CloudinaryAPISecret string
	JWTSecret           string
}

var AppConfig *Config

func init() {
	log.Println("Initializing config...")
	LoadConfig()
}

func LoadConfig() {
	AppConfig = &Config{
		DatabaseURL:         getEnvOrDefault("DATABASE_URL", "secret_default_db_url"),
		CloudinaryCloudName: getEnvOrDefault("CLOUDINARY_CLOUD_NAME", "secret_default_cloud_name"),
		CloudinaryAPIKey:    getEnvOrDefault("CLOUDINARY_API_KEY", "secret_default_api_key"),
		CloudinaryAPISecret: getEnvOrDefault("CLOUDINARY_API_SECRET", "secret_default_api_secret"),
		JWTSecret:           getEnvOrDefault("JWT_SECRET", "secret_default_jwt_secret"),
	}

	log.Println("✓ Configuration loaded successfully")
	log.Printf("  Database: %v", AppConfig.DatabaseURL != "")
	log.Printf("  Cloudinary: %v", AppConfig.CloudinaryCloudName != "")
	log.Printf("  JWT Secret: %v", AppConfig.JWTSecret != "")
}

func getEnvOrDefault(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		log.Printf("  Using default value for %s", key)
		return defaultValue
	}
	log.Printf("  Using env var for %s", key)
	return value
}
