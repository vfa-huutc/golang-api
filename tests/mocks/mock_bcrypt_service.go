package mocks

import (
	"github.com/stretchr/testify/mock"
)

// MockBcryptService is a mock implementation of IBcryptService.
type MockBcryptService struct {
	mock.Mock
}

func (m *MockBcryptService) HashPassword(password string) (string, error) {
	args := m.Called(password)
	return args.String(0), args.Error(1)
}

func (m *MockBcryptService) CheckPasswordHash(password, hashPassword string) bool {
	args := m.Called(password, hashPassword)
	return args.Bool(0)
}

func (m *MockBcryptService) HashPasswordWithCost(password string, cost int) (string, error) {
	args := m.Called(password, cost)
	return args.String(0), args.Error(1)
}
