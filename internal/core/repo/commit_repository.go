package repo

import (
	"fmt"

	"github.com/Rhaqim/savannahtech/internal/core/entities"
	"gorm.io/gorm"
)

type CommitRepo struct {
	db *gorm.DB
}

func NewCommitRepo(db *gorm.DB) *CommitRepo {
	return &CommitRepo{db: db}
}

func (c *CommitRepo) GetCommitsByRepository(repoName string, limit, offset int) ([]entities.Commit, error) {
	var commits []entities.Commit
	if err := c.db.Joins("JOIN repository_stores ON repository_stores.owner_repository = commit_stores.owner_repository").
		Where("repository_stores.name = ?", repoName).
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
	if err := c.db.Model(c).
		Select("author, COUNT(*) as commit_count").
		Group("author").
		Order("commit_count DESC").
		Limit(n).
		Scan(&results).Error; err != nil {
		return nil, fmt.Errorf("error retrieving top commit authors: %w", err)
	}

	return results, nil
}