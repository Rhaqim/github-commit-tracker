package config

import (
	"log"
	"os"
	"sync"

	"github.com/joho/godotenv"
)

type AppConfig struct {
	ServerAddress    string
	DatabaseURL      string
	DefaultOwner     string
	DefaultRepo      string
	DefaultStartDate string
	RefetchInterval  string
}

var (
	config *AppConfig
	once   sync.Once
)

func LoadConfig() {
	once.Do(func() {
		err := godotenv.Load()
		if err != nil {
			log.Fatal("Error loading .env file")
		}

		config = &AppConfig{
			ServerAddress:    os.Getenv("SERVER_ADDRESS"),
			DatabaseURL:      os.Getenv("DATABASE_URL"),
			DefaultOwner:     os.Getenv("DEFAULT_OWNER"),
			DefaultRepo:      os.Getenv("DEFAULT_REPO"),
			DefaultStartDate: os.Getenv("DEFAULT_START_DATE"),
			RefetchInterval:  os.Getenv("REFETCH_INTERVAL"),
		}
	})
}

func Config() *AppConfig {
	return config
}
