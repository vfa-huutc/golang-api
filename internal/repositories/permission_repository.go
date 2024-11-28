package repositories

import (
	"github.com/vfa-khuongdv/golang-cms/internal/models"
	"gorm.io/gorm"
)

type IPermission interface {
	Create(item *models.Permission) error
	GetAll() (*[]models.Permission, error)
}

type PermissionRepository struct {
	db *gorm.DB
}

func NewPermissionRepository(db *gorm.DB) *PermissionRepository {
	return &PermissionRepository{db}
}

func (repo *PermissionRepository) GetAll() (*[]models.Permission, error) {
	var permissions []models.Permission
	err := repo.db.Find(&permissions).Error
	if err != nil {
		return nil, err
	}
	return &permissions, nil
}

func (repo *PermissionRepository) Create(item *models.Permission) error {
	return repo.db.Create(item).Error
}
