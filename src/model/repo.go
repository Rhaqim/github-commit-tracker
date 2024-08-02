package model

import (
	"fmt"
	"savannahtech/src/database"

	"gorm.io/gorm"
)

type RepositoryStore struct {
	gorm.Model
	Name            string `json:"name"`
	Description     string `json:"description"`
	URL             string `json:"url"`
	Language        string `json:"language"`
	ForksCount      int    `json:"forks_count,omitempty"`
	StargazersCount int    `json:"stargazers_count,omitempty"`
	WatchersCount   int    `json:"watchers_count,omitempty"`
	RepoCreatedAt   string `json:"repo_created_at,omitempty"`
	RepoUpdatedAt   string `json:"repo_updated_at,omitempty"`
	OwnerRepository string `json:"owner_repository,omitempty" gorm:"unique"`
	Indexed         bool   `json:"indexed" default:"false"`
}

func (R *RepositoryStore) InsertRepository() error {
	var err = database.DB.Create(R).Error
	if err != nil {
		return fmt.Errorf("error inserting repository: %w", err)
	}

	return nil
}

func (R *RepositoryStore) GetRepositoryById(id uint) error {
	var err = database.DB.First(R, id).Error
	if err != nil {
		return fmt.Errorf("error retrieving repository by id: %w", err)
	}

	return nil
}

func (R *RepositoryStore) GetRepositoryByOwnerRepo(ownerRepo string) error {
	var err = database.DB.Where("owner_repository = ?", ownerRepo).First(R).Error
	if err != nil {
		return err
	}

	return nil
}

func (R *RepositoryStore) GetRepositoriesByOwner(owner string) error {
	var err = database.DB.Where("owner = ?", owner).Find(R).Error
	if err != nil {
		return fmt.Errorf("error retrieving repositories by owner: %w", err)
	}

	return nil
}

func (R *RepositoryStore) GetRepositoriesByRepo(repo string) error {
	var err = database.DB.Where("repo = ?", repo).Find(R).Error
	if err != nil {
		return fmt.Errorf("error retrieving repositories by repo: %w", err)
	}

	return nil
}

func (R *RepositoryStore) GetRepositories() error {
	var err = database.DB.Find(R).Error
	if err != nil {
		return fmt.Errorf("error retrieving repositories: %w", err)
	}

	return nil
}

func (R *RepositoryStore) UpdateRepository() error {
	var err = database.DB.Save(R).Error
	if err != nil {
		return fmt.Errorf("error updating repository: %w", err)
	}

	return nil
}

func (R *RepositoryStore) DeleteRepository() error {
	var err = database.DB.Delete(R).Error
	if err != nil {
		return fmt.Errorf("error deleting repository: %w", err)
	}

	return nil
}
