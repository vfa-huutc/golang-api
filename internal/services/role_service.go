package services

import (
	"github.com/vfa-khuongdv/golang-cms/internal/models"
	"github.com/vfa-khuongdv/golang-cms/internal/repositories"
	"github.com/vfa-khuongdv/golang-cms/pkg/errors"
)

type IRoleService interface {
	GetByID(id int64) (*models.Role, error)
	Create(role *models.Role) error
	Update(role *models.Role) error
	Delete(id int64) error
}

type RoleService struct {
	repo repositories.IRoleRepository
}

func NewRoleService(repo repositories.IRoleRepository) *RoleService {
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
	data, err := service.repo.GetByID(id)
	if err != nil {
		return nil, errors.New(errors.ErrDatabaseQuery, err.Error())
	}
	return data, nil
}

// Create adds a new role to the repository
// Parameters:
//   - role: The role object to be created
//
// Returns:
//   - error: Any error that occurred during the operation
func (service *RoleService) Create(role *models.Role) error {
	err := service.repo.Create(role)
	if err != nil {
		return errors.New(errors.ErrDatabaseInsert, err.Error())
	}
	return nil
}

// Update modifies an existing role in the repository
// Parameters:
//   - role: The role object containing the updated information
//
// Returns:
//   - error: Any error that occurred during the operation
func (service *RoleService) Update(role *models.Role) error {
	err := service.repo.Update(role)
	if err != nil {
		return errors.New(errors.ErrDatabaseUpdate, err.Error())
	}
	return nil
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
		return errors.New(errors.ErrDatabaseDelete, err.Error())
	}
	return service.repo.Delete(role)
}
