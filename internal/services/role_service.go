package services

import (
	"github.com/vfa-khuongdv/golang-cms/internal/models"
	"github.com/vfa-khuongdv/golang-cms/internal/repositories"
)

type IRoleService interface {
	GetByID(id int64) (*models.Role, error)
	Create(role *models.Role) error
	Update(role *models.Role) error
	Delete(id int64) error
}

type RoleService struct {
	repo *repositories.RoleRepository
}

func NewRoleService(repo *repositories.RoleRepository) *RoleService {
	return &RoleService{
		repo: repo,
	}
}

// Get retrieves a role by its ID from the repository
// Parameters:
//   - id: The unique identifier of the role to retrieve
//
// Returns:
//   - *models.Role: The role object if found
//   - error: Any error that occurred during the operation
func (service *RoleService) GetByID(id int64) (*models.Role, error) {
	return service.repo.GetByID(id)
}

// Create adds a new role to the repository
// Parameters:
//   - role: The role object to be created
//
// Returns:
//   - error: Any error that occurred during the operation
func (service *RoleService) Create(role *models.Role) error {
	return service.repo.Create(role)
}

// Update modifies an existing role in the repository
// Parameters:
//   - role: The role object containing the updated information
//
// Returns:
//   - error: Any error that occurred during the operation
func (service *RoleService) Update(role *models.Role) error {
	return service.repo.Update(role)
}

// Delete removes a role from the repository by its ID
// Parameters:
//   - id: The unique identifier of the role to delete
//
// Returns:
//   - error: Any error that occurred during the operation, including if the role is not found
func (service *RoleService) Delete(id int64) error {
	role, err := service.repo.GetByID(id)
	if err != nil {
		return err
	}
	return service.repo.Delete(role)
}
