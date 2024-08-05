package config

import (
	"os"
	"sync"

	"github.com/Rhaqim/savannahtech/pkg/logger"
	"github.com/joho/godotenv"
)

type AppConfig struct {
	ServerPort       string
	DatabaseURL      string
	DefaultOwner     string
	DefaultRepo      string
	DefaultStartDate string
	RefetchInterval  string
}

var (
	Config *AppConfig
	once   sync.Once
)

func LoadConfig() {
	once.Do(func() {
		err := godotenv.Load()
		if err != nil {
			logger.ErrorLogger.Printf("Error loading .env file")
		}

		Config = &AppConfig{
			ServerPort:       os.Getenv("SERVER_PORT"),
			DatabaseURL:      os.Getenv("DATABASE_URL"),
			DefaultOwner:     os.Getenv("DEFAULT_OWNER"),
			DefaultRepo:      os.Getenv("DEFAULT_REPO"),
			DefaultStartDate: os.Getenv("DEFAULT_START_DATE"),
			RefetchInterval:  os.Getenv("REFETCH_INTERVAL"),
		}
	})

}
