package model

import (
	"log"
	"savannahtech/src/database"
)

func Migrations() error {
	log.Println("Running database migrations...")

	err := database.DB.AutoMigrate(CommitStore{}, RepositoryStore{})
	if err != nil {
		return err
	}
	return nil
}
