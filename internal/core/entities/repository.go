package entities

import "gorm.io/gorm"

// Repository represents the application repository entity.
type Repository struct {
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
	OwnerRepository string `json:"owner_repository,omitempty" gorm:"index;unique"`
	Indexed         bool   `json:"indexed" default:"false"`
}
