package repositories

import (
	"github.com/vfa-khuongdv/golang-cms/internal/models"
	"gorm.io/gorm"
)

type IPermissionRepository interface {
	Create(item *models.Permission) error
	GetAll() ([]models.Permission, error)
}

type PermissionRepository struct {
	db *gorm.DB
}

func NewPermissionRepository(db *gorm.DB) *PermissionRepository {
	return &PermissionRepository{db: db}
}

// GetAll retrieves all permission records from the database
// Returns:
//   - *[]models.Permission: pointer to slice containing all permissions
//   - error: nil if successful, error if the database operation fails
func (repo *PermissionRepository) GetAll() ([]models.Permission, error) {
	var permissions []models.Permission

	if err := repo.db.Find(&permissions).Error; err != nil {
		return nil, err
	}
	return permissions, nil
}

// Create adds a new permission record to the database
// Parameters:
//   - item: pointer to the Permission model to be created
//
// Returns:
//   - error: nil if successful, error if the database operation fails
func (repo *PermissionRepository) Create(item *models.Permission) error {
	return repo.db.Create(item).Error
}
