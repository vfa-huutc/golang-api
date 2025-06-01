package mocks

import (
	"github.com/stretchr/testify/mock"
	"github.com/vfa-khuongdv/golang-cms/internal/models"
)

type MockRoleRepository struct {
	mock.Mock
}

func (m *MockRoleRepository) GetByID(id int64) (*models.Role, error) {
	args := m.Called(id)
	return args.Get(0).(*models.Role), args.Error(1)
}

func (m *MockRoleRepository) Create(role *models.Role) error {
	args := m.Called(role)
	return args.Error(0)
}

func (m *MockRoleRepository) Update(role *models.Role) error {
	args := m.Called(role)
	return args.Error(0)
}

func (m *MockRoleRepository) Delete(role *models.Role) error {
	args := m.Called(role)
	return args.Error(0)
}
