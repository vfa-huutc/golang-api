package mocks

import (
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

type MockDB struct {
	mock.Mock
	*gorm.DB
}

func (m *MockDB) Begin() *gorm.DB {
	args := m.Called()
	return args.Get(0).(*gorm.DB)
}

func (m *MockDB) Commit() *gorm.DB {
	m.Called()
	return &gorm.DB{}
}

func (m *MockDB) Rollback() *gorm.DB {
	m.Called()
	return &gorm.DB{}
}

type MockTx struct {
	mock.Mock
	Error error
}

func (m *MockTx) Create(value interface{}) *gorm.DB {
	args := m.Called(value)
	return &gorm.DB{Error: args.Error(0)}
}
