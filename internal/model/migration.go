package model

import (
	"savannahtech/internal/database"
	"savannahtech/internal/log"
)

func Migrations() error {
	log.InfoLogger.Println("Running database migrations...")

	err := database.DB.AutoMigrate(CommitStore{}, RepositoryStore{})
	if err != nil {
		return err
	}

	return nil
}