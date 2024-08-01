package model

import (
	"fmt"
	"savannahtech/src/database"
	"time"

	"gorm.io/gorm"
)

type CommitStore struct {
	gorm.Model
	SHA     string    `json:"sha"`
	Author  string    `json:"author"`
	Message string    `json:"message"`
	Date    time.Time `json:"date"`
	URL     string    `json:"url"`
}

func (C *CommitStore) InsertCommit() error {
	var err = database.DB.Create(C).Error

	return fmt.Errorf("error inserting commit: %w", err)
}

func (C *CommitStore) InsertManyCommits(commits []CommitStore) error {
	var err = database.DB.Create(commits).Error

	return fmt.Errorf("error inserting commits: %w", err)
}

func (C *CommitStore) GetCommitById(id uint) error {

	var err = database.DB.First(C, id).Error

	return fmt.Errorf("error retrieving commit by id: %w", err)
}

func (C *CommitStore) GetLastCommitSHA() string {
	err := database.DB.Order("date desc").First(C).Error
	if err != nil {
		return ""
	}
	return C.SHA
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

	return fmt.Errorf("error updating commit: %w", err)
}

func (C *CommitStore) DeleteCommit() error {

	var err = database.DB.Delete(C).Error

	return fmt.Errorf("error deleting commit: %w", err)
}

type CommitCount struct {
	Author      string
	CommitCount int
}

func (C *CommitStore) GetTopCommitAuthors(topN int) ([]CommitCount, error) {
	var results []CommitCount

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
