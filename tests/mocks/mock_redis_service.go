package mocks

import (
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/mock"
)

// MockRedisService is a testify mock for IRedisService
type MockRedisService struct {
	mock.Mock
}

func (m *MockRedisService) Set(key string, value any, ttl time.Duration) error {
	args := m.Called(key, value, ttl)
	return args.Error(0)
}

func (m *MockRedisService) Get(key string) (string, error) {
	args := m.Called(key)
	return args.String(0), args.Error(1)
}

func (m *MockRedisService) Delete(key string) error {
	args := m.Called(key)
	return args.Error(0)
}

func (m *MockRedisService) Exists(key string) (bool, error) {
	args := m.Called(key)
	return args.Bool(0), args.Error(1)
}

func (m *MockRedisService) GetClient() redis.Cmdable {
	args := m.Called()
	val := args.Get(0)
	if val == nil {
		return nil
	}
	return val.(redis.Cmdable)
}
