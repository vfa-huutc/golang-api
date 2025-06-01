package mocks

import (
	"github.com/stretchr/testify/mock"
	"github.com/vfa-khuongdv/golang-cms/internal/models"
)

type MockPermissionRepository struct {
	mock.Mock
}

func (m *MockPermissionRepository) Create(item *models.Permission) error {
	args := m.Called(item)
	return args.Error(0)
}

func (m *MockPermissionRepository) GetAll() ([]models.Permission, error) {
	args := m.Called()
	return args.Get(0).([]models.Permission), args.Error(1)
}
