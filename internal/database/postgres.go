package database

import (
	"github.com/Rhaqim/savannahtech/config"
	"github.com/Rhaqim/savannahtech/internal/core/repo"
	"github.com/Rhaqim/savannahtech/internal/core/repositories"
	"github.com/Rhaqim/savannahtech/pkg/logger"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormLogger "gorm.io/gorm/logger"
)

var DB *gorm.DB // Global database connection

// InitDB initializes the database connection
func InitDB() {
	dsn := config.Config.DatabaseURL

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: gormLogger.Default.LogMode(gormLogger.Silent),
	})
	if err != nil {
		logger.ErrorLogger.Fatalf("failed to connect to database: %v", err)
	}

	DB = db

	// Initialize repositories
	repositories.RepoStore = repo.NewRepositoryRepo(DB)
	repositories.CommitStore = repo.NewCommitRepo(DB)
}

func Close() {
	sqlDB, err := DB.DB()
	if err != nil {
		logger.ErrorLogger.Fatal("Failed to close database:", err)
	}
	sqlDB.Close()
}
