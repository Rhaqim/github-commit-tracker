package database

import (
	"github.com/Rhaqim/savannahtech/pkg/logger"

	"github.com/Rhaqim/savannahtech/internal/core/entities"
)

// RunMigrations runs the database migrations
func RunMigrations() {
	err := DB.AutoMigrate(&entities.Repository{}, &entities.Commit{})
	if err != nil {
		logger.ErrorLogger.Fatalf("failed to run database migrations: %v", err)
	}
}
