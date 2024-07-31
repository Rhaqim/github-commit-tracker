package model

import (
	"fmt"
	"log"
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

func GetLastCommitSHA(db *gorm.DB) string {
	var lastCommit CommitStore
	db.Order("date desc").First(&lastCommit)
	return lastCommit.SHA
}

func SaveCommitToDB(db *gorm.DB, commit CommitStore) error {
	return db.Create(&commit).Error
}

func PeriodicFetch(repoURL string, interval time.Duration, db *gorm.DB) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			lastCommitSHA := GetLastCommitSHA(db)
			_ = fmt.Sprintf("%s/commits?since=%s", repoURL, lastCommitSHA)

			var commits []CommitStore

			// commits, err := FetchCommits(commitsURL, MakeRequest)
			// if err != nil {
			// 	log.Printf("Error fetching commits: %v", err)
			// 	continue
			// }

			// Save new commits to the database
			for _, commit := range commits {
				if err := SaveCommitToDB(db, commit); err != nil {
					log.Printf("Error saving commit to database: %v", err)
				}
			}
		}
	}
}
