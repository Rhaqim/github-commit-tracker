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

func (C *CommitStore) UpdateCommit() error {

	var err = database.DB.Save(C).Error

	return err
}

func (C *CommitStore) DeleteCommit() error {

	var err = database.DB.Delete(C).Error

	return err
}

func (C *CommitStore) GetCommitById(id uint) error {

	var err = database.DB.First(C, id).Error

	return err
}

func (C *CommitStore) GetCommitsByOwnerAndRepo(owner, repo string) error {

	var err = database.DB.Where("owner = ? AND repo = ?", owner, repo).Find(C).Error

	return err
}

func (C *CommitStore) GetCommitsByOwner(owner string) error {

	var err = database.DB.Where("owner = ?", owner).Find(C).Error

	return err
}

func (C *CommitStore) GetCommitsByRepo(repo string) error {

	var err = database.DB.Where("repo = ?", repo).Find(C).Error

	return err
}

func (C *CommitStore) GetCommits() ([]CommitStore, error) {

	var commits []CommitStore
	err := database.DB.Find(&commits).Error
	if err != nil {
		return nil, err
	}

	return commits, nil
}
