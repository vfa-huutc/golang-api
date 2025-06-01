package mocks

import (
	"github.com/stretchr/testify/mock"
	"github.com/vfa-khuongdv/golang-cms/internal/models"
	"github.com/vfa-khuongdv/golang-cms/internal/services"
)

type MockRoleService struct {
	mock.Mock
	services.IRoleService
}

// Create mocks the Create method of the role service interface.
// It records the call with the provided role model and returns the error
// configured in the test setup.
func (m *MockRoleService) Create(role *models.Role) error {
	args := m.Called(role)
	return args.Error(0)
}

// GetByID mocks the retrieval of a role by its ID.
// It records the call with the provided ID and returns the configured role model and error.
func (m *MockRoleService) GetByID(id int64) (*models.Role, error) {
	args := m.Called(id)
	var role *models.Role
	if r := args.Get(0); r != nil {
		role = r.(*models.Role)
	}
	return role, args.Error(1)
}

// Update mocks the Update method of the role service interface.
// It records the call with the provided role model and returns the error
// configured in the test setup.
func (m *MockRoleService) Update(role *models.Role) error {
	args := m.Called(role)
	return args.Error(0)
}

// Delete mocks the Delete method of the role service interface.
// It records the call with the provided role ID and returns the error
// configured in the test setup.
func (m *MockRoleService) Delete(id int64) error {
	args := m.Called(id)
	return args.Error(0)
}

// AssignPermissions mocks the AssignPermissions method of the role service interface.
// It records the call with the provided role ID and permission IDs array, then returns
// the error configured in the test setup.
func (m *MockRoleService) AssignPermissions(roleID uint, permissionIDs []uint) error {
	args := m.Called(roleID, permissionIDs)
	return args.Error(0)
}

// GetRolePermissions mocks the GetRolePermissions method of the role service interface.
// It records the call with the provided role ID and returns the configured list of
// permission models and error from the test setup.
func (m *MockRoleService) GetRolePermissions(roleID uint) ([]models.Permission, error) {
	args := m.Called(roleID)
	return args.Get(0).([]models.Permission), args.Error(1)
}
