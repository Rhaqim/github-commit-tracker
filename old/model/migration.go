package model

import (
	"github.com/Rhaqim/savannahtech/internal/database"
	"github.com/Rhaqim/savannahtech/pkg/logger"
)

func Migrations() error {
	logger.InfoLogger.Println("Running database migrations...")

	err := database.DB.AutoMigrate(CommitStore{}, RepositoryStore{})
	if err != nil {
		return err
	}

	return nil
}
