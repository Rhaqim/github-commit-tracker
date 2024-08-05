package repositories

import "github.com/Rhaqim/savannahtech/internal/core/entities"

type RepositoryStoreRepo interface {
	ProcessRepository(entities.Repository) error
	GetRepositoryByOwnerRepo(ownerRepo string) (entities.Repository, error)
}

var RepoStore RepositoryStoreRepo
