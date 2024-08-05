package utils

import (
	"os"

	"github.com/Rhaqim/savannahtech/pkg/logger"
	"github.com/joho/godotenv"
)

// Environment variables
func Env(key, fallback string) string {
	if err := godotenv.Load(); err != nil {
		logger.ErrorLogger.Println("No .env file found error: ", err)
	}

	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
