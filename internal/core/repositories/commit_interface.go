package repositories

import "github.com/Rhaqim/savannahtech/internal/core/entities"

type CommitStoreRepo interface {
	CreateCommits(commit []entities.Commit) error
	GetCommitsByRepository(repoName string, limit, offset int) ([]entities.Commit, error)
	GetTopNCommitAuthors(n int) ([]entities.CommitCount, error)
	GetLastCommitDate(ownerRepo string) string
}

var CommitStore CommitStoreRepo
