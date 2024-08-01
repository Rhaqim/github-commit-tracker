package model

import (
	"savannahtech/src/database"
	"savannahtech/src/log"
)

func Migrations() error {
	log.InfoLogger.Println("Running database migrations...")

	err := database.DB.AutoMigrate(CommitStore{}, RepositoryStore{})
	if err != nil {
		return err
	}

	return nil
}
