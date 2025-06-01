package mocks

import (
	"github.com/stretchr/testify/mock"
	"github.com/vfa-khuongdv/golang-cms/internal/models"
)

type MockSettingService struct {
	mock.Mock
}

func (m *MockSettingService) GetSetting() ([]models.Setting, error) {
	args := m.Called()
	return args.Get(0).([]models.Setting), args.Error(1)
}

func (m *MockSettingService) GetSettingByKey(key string) (*models.Setting, error) {
	args := m.Called(key)
	return args.Get(0).(*models.Setting), args.Error(1)
}

func (m *MockSettingService) Update(setting *models.Setting) error {
	args := m.Called(setting)
	return args.Error(0)
}

func (m *MockSettingService) Create(setting *models.Setting) error {
	args := m.Called(setting)
	return args.Error(0)
}
