package model

import (
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

	return err
}

func (C *CommitStore) InsertManyCommits(commits []CommitStore) error {
	var err = database.DB.Create(commits).Error

	return err
}

func (C *CommitStore) GetCommitById(id uint) error {

	var err = database.DB.First(C, id).Error

	return err
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
		return nil, err
	}

	return commits, nil
}

func (C *CommitStore) UpdateCommit() error {

	var err = database.DB.Save(C).Error

	return err
}

func (C *CommitStore) DeleteCommit() error {

	var err = database.DB.Delete(C).Error

	return err
}
