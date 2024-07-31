package model

import (
	"fmt"
	"savannahtech/src/database"
	"time"

	"gorm.io/gorm"
)

type CommitStore struct {
	gorm.Model
	Author  string    `json:"author"`
	Message string    `json:"message"`
	Date    time.Time `json:"date"`
	URL     string    `json:"url"`
}

func (C *CommitStore) InsertCommit() error {
	fmt.Println("Storing Commit", C)
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

func (C *CommitStore) GetCommits() error {

	var err = database.DB.Find(C).Error

	return err
}
