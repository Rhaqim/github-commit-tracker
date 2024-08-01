package model

import (
	"fmt"
	"savannahtech/src/database"
	"savannahtech/src/types"
	"time"

	"gorm.io/gorm"
)

type CommitStore struct {
	gorm.Model
	SHA             string    `json:"sha"`
	Author          string    `json:"author"`
	Message         string    `json:"message"`
	Date            time.Time `json:"date"`
	URL             string    `json:"url"`
	OwnerRepository string    `json:"owner_repository,omitempty"`
}

func (C *CommitStore) InsertCommit() error {
	var err = database.DB.Create(C).Error
	if err != nil {
		return fmt.Errorf("error inserting commit: %w", err)
	}

	return nil
}

func (C *CommitStore) InsertManyCommits(commits []CommitStore) error {
	var err = database.DB.Create(commits).Error
	if err != nil {
		return fmt.Errorf("error inserting commits: %w", err)
	}

	return nil
}

func (C *CommitStore) GetCommitById(id uint) error {

	var err = database.DB.First(C, id).Error
	if err != nil {
		return fmt.Errorf("error retrieving commit by id: %w", err)
	}

	return nil
}

func (C *CommitStore) GetLastCommitDate() time.Time {
	err := database.DB.Order("date desc").First(C).Error
	if err != nil {
		return time.Time{}
	}
	return C.Date
}

func (C *CommitStore) GetCommits() ([]CommitStore, error) {

	var commits []CommitStore
	err := database.DB.Find(&commits).Error
	if err != nil {
		return nil, fmt.Errorf("error retrieving commits: %w", err)
	}

	return commits, nil
}

func (C *CommitStore) UpdateCommit() error {

	var err = database.DB.Save(C).Error
	if err != nil {
		return fmt.Errorf("error updating commit: %w", err)
	}

	return nil
}

func (C *CommitStore) DeleteCommit() error {

	var err = database.DB.Delete(C).Error
	if err != nil {
		return fmt.Errorf("error deleting commit: %w", err)
	}

	return nil
}

func (C *CommitStore) GetTopCommitAuthors(topN int) ([]types.CommitCount, error) {
	var results []types.CommitCount

	// Perform the query
	err := database.DB.Model(C).
		Select("author, COUNT(*) as commit_count").
		Group("author").
		Order("commit_count DESC").
		Limit(topN).
		Scan(&results).Error

	if err != nil {
		return nil, fmt.Errorf("error retrieving top commit authors: %w", err)
	}

	return results, nil
}

func (C *CommitStore) GetCommitsByAuthor(repoName string) ([]CommitStore, error) {
	var commits []CommitStore

	// Perform a join between the RepositoryStore and CommitStore tables
	err := database.DB.Joins("JOIN repository_stores ON repository_stores.owner_repository = commit_stores.owner_repository").
		Where("repository_stores.name = ?", repoName).
		Find(&commits).Error

	if err != nil {
		return nil, fmt.Errorf("error retrieving commits for repository %s: %w", repoName, err)
	}

	return commits, nil
}
