package mocks

import (
	"github.com/stretchr/testify/mock"
	"github.com/vfa-khuongdv/golang-cms/internal/models"
)

type MockSettingRepository struct {
	mock.Mock
}

func (m *MockSettingRepository) GetAll() ([]models.Setting, error) {
	args := m.Called()
	return args.Get(0).([]models.Setting), args.Error(1)
}

func (m *MockSettingRepository) GetByKey(key string) (*models.Setting, error) {
	args := m.Called(key)
	return args.Get(0).(*models.Setting), args.Error(1)
}

func (m *MockSettingRepository) Update(setting *models.Setting) error {
	args := m.Called(setting)
	return args.Error(0)
}

func (m *MockSettingRepository) Create(setting *models.Setting) error {
	args := m.Called(setting)
	return args.Error(0)
}
