package config

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Config struct {
	DB            *gorm.DB
	Port          string
	JWTKey        string
	JWTExpiration int
	UploadStorage string
	AppURL        string
}

func LoadConfig() *Config {
	// Try loading from nested directory first, fallback to root if running from somewhere else
	err := godotenv.Load("internal/config/env/.env")
	if err != nil {
		err = godotenv.Load(".env")
		if err != nil {
			log.Println("No .env file found in internal/config/env/ or root, using OS env vars")
		}
	}

	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbUser := os.Getenv("DB_USER")
	dbPass := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Jakarta",
		dbHost, dbUser, dbPass, dbName, dbPort)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	jwtKey := os.Getenv("JWT_SECRET")
	if jwtKey == "" {
		jwtKey = "secret"
	}

	uploadStorage := os.Getenv("UPLOAD_STORAGE")
	if uploadStorage == "" {
		uploadStorage = "local" // Default to local storage
	}

	appURL := os.Getenv("APP_URL")
	if appURL == "" {
		appURL = "http://localhost:" + port // Default to localhost
	}

	jwtExpStr := os.Getenv("JWT_EXPIRATION_SECONDS")
	jwtExpSeconds, err := strconv.Atoi(jwtExpStr)
	if err != nil || jwtExpSeconds <= 0 {
		jwtExpSeconds = 86400 // default 24 hours in seconds
	}

	return &Config{
		DB:            db,
		Port:          port,
		JWTKey:        jwtKey,
		JWTExpiration: jwtExpSeconds,
		UploadStorage: uploadStorage,
		AppURL:        appURL,
	}
}
