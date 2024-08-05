package services

import (
	"strconv"

	"github.com/Rhaqim/savannahtech/internal/core/entities"
	"github.com/Rhaqim/savannahtech/internal/core/repositories"
	"github.com/Rhaqim/savannahtech/internal/utils"
)

func FetchCommitsByRepository(repoName, pageStr, sizeStr string) ([]entities.Commit, error) {
	size, offset := utils.Paigenation(pageStr, sizeStr)

	return repositories.CommitStore.GetCommitsByRepository(repoName, size, offset)
}

func FetchTopNCommitAuthors(n string) ([]entities.CommitCount, error) {
	nInt, err := strconv.Atoi(n)
	if err != nil {
		return nil, err
	}

	return repositories.CommitStore.GetTopNCommitAuthors(nInt)
}
