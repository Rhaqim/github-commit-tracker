package repositories

import "github.com/Rhaqim/savannahtech/internal/core/entities"

type CommitStoreRepo interface {
	GetCommitsByRepository(repoName string, limit, offset int) ([]entities.Commit, error)
	GetTopNCommitAuthors(n int) ([]entities.CommitCount, error)
}

var CommitStore CommitStoreRepo
