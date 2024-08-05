package repo

import (
	"fmt"

	"github.com/Rhaqim/savannahtech/internal/core/entities"
	"github.com/Rhaqim/savannahtech/pkg/logger"
	"gorm.io/gorm"
)

type CommitRepo struct {
	db *gorm.DB
}

func NewCommitRepo(db *gorm.DB) *CommitRepo {
	return &CommitRepo{db: db}
}

func (c *CommitRepo) CreateCommits(commit []entities.Commit) error {
	if err := c.db.Create(&commit).Error; err != nil {
		return fmt.Errorf("error creating commit: %w", err)
	}
	return nil
}

func (c *CommitRepo) GetCommitsByRepository(repoName string, limit, offset int) ([]entities.Commit, error) {
	var commits []entities.Commit

	logger.InfoLogger.Printf("Fetching commits for repository %s\n", repoName)

	if err := c.db.Joins("JOIN repositories ON repositories.owner_repository = commits.owner_repository").
		Where("repositories.name = ?", repoName).
		Limit(limit).
		Offset(offset).
		Find(&commits).Error; err != nil {
		return nil, fmt.Errorf("error retrieving commits for repository %s: %w", repoName, err)
	}
	return commits, nil
}

func (c *CommitRepo) GetTopNCommitAuthors(n int) ([]entities.CommitCount, error) {
	var results []entities.CommitCount

	// Perform the query
	if err := c.db.Model(entities.Commit{}).
		Select("author, COUNT(*) as commit_count").
		Group("author").
		Order("commit_count DESC").
		Limit(n).
		Scan(&results).Error; err != nil {
		return nil, fmt.Errorf("error retrieving top commit authors: %w", err)
	}

	return results, nil
}

func (c *CommitRepo) GetLastCommitDate(ownerRepo string) string {
	var commit entities.Commit
	if err := c.db.Where("owner_repository = ?", ownerRepo).Order("date desc").First(&commit).Error; err != nil {
		return ""
	}
	return commit.Date.Format("2006-01-02T15:04:05Z")
}
