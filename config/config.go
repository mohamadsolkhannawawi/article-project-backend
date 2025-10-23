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
		DatabaseURL:         getEnvOrDefault("DATABASE_URL", "postgresql://neondb_owner:npg_Ai1ZmSbfC4pn@ep-broad-paper-addvxhil-pooler.c-2.us-east-1.aws.neon.tech/neondb?sslmode=require&channel_binding=require"),
		CloudinaryCloudName: getEnvOrDefault("CLOUDINARY_CLOUD_NAME", "dl5agg7km"),
		CloudinaryAPIKey:    getEnvOrDefault("CLOUDINARY_API_KEY", "699414449319264"),
		CloudinaryAPISecret: getEnvOrDefault("CLOUDINARY_API_SECRET", "Wv_ubBCzd7ct-wT03MQTgfI3xEA"),
		JWTSecret:           getEnvOrDefault("JWT_SECRET", "R39ZKnV5RKqJNVInJWnjxJjCeouk048I758uZyUHgobTabvdxb7mYCnp42tbFOa1"),
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
