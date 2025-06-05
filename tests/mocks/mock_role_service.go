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

func (m *MockRoleService) Create(role *models.Role) error {
	args := m.Called(role)
	return args.Error(0)
}

func (m *MockRoleService) GetByID(id int64) (*models.Role, error) {
	args := m.Called(id)
	var role *models.Role
	if r := args.Get(0); r != nil {
		role = r.(*models.Role)
	}
	return role, args.Error(1)
}

func (m *MockRoleService) Update(role *models.Role) error {
	args := m.Called(role)
	return args.Error(0)
}

func (m *MockRoleService) Delete(id int64) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockRoleService) AssignPermissions(roleID uint, permissionIDs []uint) error {
	args := m.Called(roleID, permissionIDs)
	return args.Error(0)
}

func (m *MockRoleService) GetRolePermissions(roleID uint) ([]models.Permission, error) {
	args := m.Called(roleID)
	return args.Get(0).([]models.Permission), args.Error(1)
}
