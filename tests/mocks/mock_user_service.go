package mocks

import (
	"github.com/stretchr/testify/mock"
	"github.com/vfa-khuongdv/golang-cms/internal/models"
	"github.com/vfa-khuongdv/golang-cms/internal/utils"
)

type MockUserService struct {
	mock.Mock
}

func (m *MockUserService) PaginateUser(page, limit int) (*utils.Pagination, error) {
	args := m.Called(page, limit)
	return args.Get(0).(*utils.Pagination), args.Error(1)
}

func (m *MockUserService) GetUser(id uint) (*models.User, error) {
	args := m.Called(id)
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserService) GetUserByEmail(email string) (*models.User, error) {
	args := m.Called(email)
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserService) CreateUser(user *models.User, roleIds []uint) error {
	args := m.Called(user, roleIds)
	return args.Error(0)
}

func (m *MockUserService) UpdateUser(user *models.User) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *MockUserService) DeleteUser(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockUserService) GetUserByToken(token string) (*models.User, error) {
	args := m.Called(token)
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserService) GetProfile(id uint) (*models.User, error) {
	args := m.Called(id)
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserService) UpdateProfile(user *models.User) error {
	args := m.Called(user)
	return args.Error(0)
}
