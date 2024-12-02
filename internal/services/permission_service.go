package services

import (
	"github.com/vfa-khuongdv/golang-cms/internal/models"
	"github.com/vfa-khuongdv/golang-cms/internal/repositories"
)

type IPermissionService interface {
	GetAll() (*[]models.Permission, error)
}

type PermissionService struct {
	repo *repositories.PermissionRepository
}

func NewPermissionService(repo *repositories.PermissionRepository) *PermissionService {
	return &PermissionService{repo: repo}
}

// GetAll retrieves all permissions from the repository
// Returns:
//   - *[]models.Permission: Pointer to slice of Permission models containing all permissions
//   - error: Error if any occurred during the operation
func (repo *PermissionService) GetAll() (*[]models.Permission, error) {
	return repo.repo.GetAll()
}
