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

func (repo *PermissionService) GetAll() (*[]models.Permission, error) {
	return repo.repo.GetAll()
}
