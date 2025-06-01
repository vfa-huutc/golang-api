package mocks

import (
	"github.com/stretchr/testify/mock"
	"github.com/vfa-khuongdv/golang-cms/internal/models"
)

type MockPermissionService struct {
	mock.Mock
}

func (m *MockPermissionService) GetAll() ([]models.Permission, error) {
	args := m.Called()
	if perms, ok := args.Get(0).([]models.Permission); ok {
		return perms, args.Error(1)
	}
	return nil, args.Error(1)
}
