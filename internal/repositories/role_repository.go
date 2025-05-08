package repositories

import (
	"github.com/vfa-khuongdv/golang-cms/internal/models"
	"gorm.io/gorm"
)

type IRoleRepository interface {
	GetByID(id int64) (*models.Role, error)
	Create(role *models.Role) error
	Update(role *models.Role) error
	Delete(role *models.Role) error
	AssignPermissions(roleID uint, permissionIDs []uint) error
	GetRolePermissions(roleID uint) ([]models.Permission, error)
	GetDB() *gorm.DB
}

type RoleRepository struct {
	db *gorm.DB
}

func NewRoleRepository(db *gorm.DB) *RoleRepository {
	return &RoleRepository{db: db}
}

// GetByID retrieves a role by its ID from the database
// Parameters:
//   - id: The unique identifier of the role to retrieve
//
// Returns:
//   - *models.Role: Pointer to the retrieved role if found
//   - error: nil if successful, error message if failed
func (repo *RoleRepository) GetByID(id int64) (*models.Role, error) {
	var role models.Role
	if err := repo.db.First(&role, id).Error; err != nil {
		return nil, err
	}
	return &role, nil
}

// Create inserts a new role record into the database
// Parameters:
//   - role: Pointer to the role model to be created
//
// Returns:
//   - error: nil if successful, error message if failed
func (repo *RoleRepository) Create(role *models.Role) error {
	return repo.db.Create(role).Error
}

// Update modifies an existing role record in the database
// Parameters:
//   - role: Pointer to the role model containing updated data
//
// Returns:
//   - error: nil if successful, error message if failed
func (repo *RoleRepository) Update(role *models.Role) error {
	return repo.db.Save(role).Error
}

// Delete removes a role record from the database
// Parameters:
//   - role: Pointer to the role model to be deleted
//
// Returns:
//   - error: nil if successful, error message if failed
func (repo *RoleRepository) Delete(role *models.Role) error {
	return repo.db.Delete(role).Error
}

// AssignPermissions assigns a list of permissions to a role
// This implementation replaces the existing permissions with the new set
// Parameters:
//   - roleID: The ID of the role to assign permissions to
//   - permissionIDs: Slice of permission IDs to assign to the role
//
// Returns:
//   - error: nil if successful, error message if failed
func (repo *RoleRepository) AssignPermissions(roleID uint, permissionIDs []uint) error {
	// Start a transaction
	tx := repo.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Delete existing role-permission relationships for this role
	if err := tx.Unscoped().Delete(&models.RolePermission{}, "role_id = ?", roleID).Error; err != nil {
		tx.Rollback()
		return err
	}

	// Create new role-permission relationships
	for _, permID := range permissionIDs {
		rolePermission := models.RolePermission{
			RoleID:       roleID,
			PermissionID: permID,
		}
		if err := tx.Create(&rolePermission).Error; err != nil {
			tx.Rollback()
			return err
		}
	}

	// Commit the transaction
	return tx.Commit().Error
}

// GetRolePermissions retrieves all permission objects assigned to a role
// Parameters:
//   - roleID: The ID of the role to get permissions for
//
// Returns:
//   - []models.Permission: Slice of permission objects assigned to the role
//   - error: nil if successful, error message if failed
func (repo *RoleRepository) GetRolePermissions(roleID uint) ([]models.Permission, error) {
	var role models.Role
	repo.db.Preload("Permissions").First(&role, roleID)
	if role.ID == 0 {
		return nil, gorm.ErrRecordNotFound
	}

	return role.Permissions, nil
}

// GetDB returns the database instance
// Returns:
//   - *gorm.DB: The database instance
func (repo *RoleRepository) GetDB() *gorm.DB {
	return repo.db
}
