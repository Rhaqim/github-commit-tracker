package utils

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

// Environment variables
func Env(key, fallback string) string {
	_ = []string{".env", ".env.local", ".env.dev", ".env.prod", ".dockerenv"}
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found error: ", err)
	}

	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
