package model

import "savannahtech/src/database"

func Migrations() error {
	err := database.DB.AutoMigrate(CommitStore{}, RepositoryStore{})
	if err != nil {
		return err
	}
	return nil
}
