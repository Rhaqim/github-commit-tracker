package services

import (
	"github.com/Rhaqim/savannahtech/internal/core/entities"
	"github.com/Rhaqim/savannahtech/internal/core/repositories"
)

func ProcessRepository(owner, repo string) error {
	var err error
	// return repositories.RepoStore.ProcessRepository(ownerRepo)
	return err
}

func FetchRepositoryByOwnerRepo(ownerRepo string) (entities.Repository, error) {
	return repositories.RepoStore.GetRepositoryByOwnerRepo(ownerRepo)
}
