package repo

import (
	"github.com/Rhaqim/savannahtech/internal/core/entities"
	"gorm.io/gorm"
)

type RepositoryRepo struct {
	db *gorm.DB
}

func NewRepositoryRepo(db *gorm.DB) *RepositoryRepo {
	return &RepositoryRepo{db: db}
}

func (r *RepositoryRepo) StoreRepository(repo entities.Repository) error {
	if err := r.db.Create(&repo).Error; err != nil {
		return err
	}
	return nil
}

func (r *RepositoryRepo) GetRepositoryByOwnerRepo(ownerRepo string) (entities.Repository, error) {
	var repo entities.Repository
	if err := r.db.Where("owner_repository = ?", ownerRepo).First(&repo).Error; err != nil {
		return entities.Repository{}, err
	}
	return repo, nil
}

func (r *RepositoryRepo) UpdateRepositoryIndexed(ownerRepo string) error {
	if err := r.db.Model(&entities.Repository{}).Where("owner_repository = ?", ownerRepo).Update("indexed", true).Error; err != nil {
		return err
	}
	return nil
}
