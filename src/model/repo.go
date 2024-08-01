package model

import (
	"fmt"
	"savannahtech/src/database"

	"gorm.io/gorm"
)

type RepositoryStore struct {
	gorm.Model
	Name            string `json:"name" gorm:"unique"`
	Description     string `json:"description"`
	URL             string `json:"url" gorm:"unique"`
	Language        string `json:"language"`
	ForksCount      int    `json:"forks_count,omitempty"`
	StargazersCount int    `json:"stargazers_count,omitempty"`
	WatchersCount   int    `json:"watchers_count,omitempty"`
	RepoCreatedAt   string `json:"repo_created_at,omitempty"`
	RepoUpdatedAt   string `json:"repo_updated_at,omitempty"`
}

func (R *RepositoryStore) InsertRepository() error {
	var err = database.DB.Create(R).Error

	return fmt.Errorf("error inserting repository: %w", err)
}

func (R *RepositoryStore) GetRepositoryById(id uint) error {

	var err = database.DB.First(R, id).Error

	return fmt.Errorf("error retrieving repository by id: %w", err)
}

func (R *RepositoryStore) GetRepositoriesByOwnerAndRepo(owner, repo string) error {

	var err = database.DB.Where("owner = ? AND repo = ?", owner, repo).Find(R).Error

	return fmt.Errorf("error retrieving repositories by owner and repo: %w", err)
}

func (R *RepositoryStore) GetRepositoriesByOwner(owner string) error {

	var err = database.DB.Where("owner = ?", owner).Find(R).Error

	return fmt.Errorf("error retrieving repositories by owner: %w", err)
}

func (R *RepositoryStore) GetRepositoriesByRepo(repo string) error {

	var err = database.DB.Where("repo = ?", repo).Find(R).Error

	return fmt.Errorf("error retrieving repositories by repo: %w", err)
}

func (R *RepositoryStore) GetRepositories() error {

	var err = database.DB.Find(R).Error

	return fmt.Errorf("error retrieving repositories: %w", err)
}

func (R *RepositoryStore) UpdateRepository() error {

	var err = database.DB.Save(R).Error

	return fmt.Errorf("error updating repository: %w", err)
}

func (R *RepositoryStore) DeleteRepository() error {

	var err = database.DB.Delete(R).Error

	return fmt.Errorf("error deleting repository: %w", err)
}
