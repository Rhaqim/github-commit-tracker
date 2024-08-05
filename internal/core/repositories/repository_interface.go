package repositories

import "github.com/Rhaqim/savannahtech/internal/core/entities"

type RepositoryStoreRepo interface {
	StoreRepository(entities.Repository) error
	GetRepositoryByOwnerRepo(ownerRepo string) (entities.Repository, error)
}

var RepoStore RepositoryStoreRepo
