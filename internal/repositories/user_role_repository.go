package repositories

import (
	"github.com/vfa-khuongdv/golang-cms/internal/models"
	"gorm.io/gorm"
)

type IUserRoleRepository interface {
	GetUserRoles(userID uint) ([]models.Role, error)
	GetDB() *gorm.DB
}

type UserRoleRepository struct {
	db *gorm.DB
}

func NewUserRoleRepository(db *gorm.DB) *UserRoleRepository {
	return &UserRoleRepository{db: db}
}

// GetUserRoles retrieves all roles assigned to a user
// Parameters:
//   - userID: The ID of the user
//
// Returns:
//   - []models.Role: Slice of roles assigned to the user
//   - error: nil if successful, error message if failed
func (repo *UserRoleRepository) GetUserRoles(userID uint) ([]models.Role, error) {
	var roles []models.Role

	err := repo.db.Table("roles").
		Joins("JOIN user_roles ON roles.id = user_roles.role_id").
		Where("user_roles.user_id = ?", userID).
		Find(&roles).Error

	return roles, err
}

// GetDB returns the database instance
// Returns:
//   - *gorm.DB: The database instance
func (repo *UserRoleRepository) GetDB() *gorm.DB {
	return repo.db
}
