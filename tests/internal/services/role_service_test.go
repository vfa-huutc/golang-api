package tests_internal_services

import (
	"testing"

	"github.com/vfa-khuongdv/golang-cms/pkg/errors"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/vfa-khuongdv/golang-cms/internal/models"
	"github.com/vfa-khuongdv/golang-cms/internal/services"
)

// MockRoleRepository is a mock of IRoleRepository interface
type MockRoleRepository struct {
	mock.Mock
}

func (m *MockRoleRepository) GetByID(id int64) (*models.Role, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
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

func TestGetByID(t *testing.T) {
	mockRepo := new(MockRoleRepository)
	roleService := services.NewRoleService(mockRepo)

	t.Run("Success", func(t *testing.T) {
		expectedRole := &models.Role{
			Name:        "admin",
			DisplayName: "Administrator",
		}

		// Setup expectations
		mockRepo.On("GetByID", int64(1)).Return(expectedRole, nil).Once()

		// Call service method
		role, err := roleService.GetByID(1)

		// Assert expectations
		assert.NoError(t, err)
		assert.Equal(t, expectedRole, role)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Not found", func(t *testing.T) {
		// Setup expectations
		mockRepo.On("GetByID", int64(999)).Return(
			nil,
			errors.New(errors.ErrDatabaseQuery, "record not found"),
		).Once()

		// Call service method
		role, err := roleService.GetByID(999)

		// Assert expectations
		assert.Error(t, err)
		assert.Nil(t, role)
		assert.Contains(t, err.Error(), "code: 2001")
		mockRepo.AssertExpectations(t)
	})
}

func TestCreate(t *testing.T) {
	mockRepo := new(MockRoleRepository)
	roleService := services.NewRoleService(mockRepo)

	t.Run("Success", func(t *testing.T) {
		role := &models.Role{
			Name:        "editor",
			DisplayName: "Content Editor",
		}

		// Setup expectations
		mockRepo.On("Create", role).Return(nil).Once()

		// Call service method
		err := roleService.Create(role)

		// Assert expectations
		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Error", func(t *testing.T) {
		role := &models.Role{
			Name:        "existing_role",
			DisplayName: "Existing Role",
		}

		// Setup expectations - simulate a database error
		mockRepo.On("Create", role).Return(
			errors.New(errors.ErrDatabaseInsert, "duplicate entry"),
		).Once()

		// Call service method
		err := roleService.Create(role)

		// Assert expectations
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "code: 2002")
		mockRepo.AssertExpectations(t)
	})
}

func TestUpdate(t *testing.T) {
	mockRepo := new(MockRoleRepository)
	roleService := services.NewRoleService(mockRepo)

	t.Run("Success", func(t *testing.T) {
		role := &models.Role{
			Name:        "moderator",
			DisplayName: "Content Moderator",
		}

		// Setup expectations
		mockRepo.On("Update", role).Return(nil).Once()

		// Call service method
		err := roleService.Update(role)

		// Assert expectations
		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Error", func(t *testing.T) {
		role := &models.Role{
			Name:        "invalid_role",
			DisplayName: "Invalid Role",
		}

		// Setup expectations - simulate a database error
		mockRepo.On("Update", role).Return(
			errors.New(errors.ErrDatabaseUpdate, "record not found"),
		).Once()

		// Call service method
		err := roleService.Update(role)

		// Assert expectations
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "code: 2003")
		mockRepo.AssertExpectations(t)
	})
}

func TestDelete(t *testing.T) {
	mockRepo := new(MockRoleRepository)
	roleService := services.NewRoleService(mockRepo)

	t.Run("Success", func(t *testing.T) {
		roleID := int64(1)
		role := &models.Role{
			Name:        "guest",
			DisplayName: "Guest User",
		}

		// Setup expectations
		mockRepo.On("GetByID", roleID).Return(role, nil).Once()
		mockRepo.On("Delete", role).Return(nil).Once()

		// Call service method
		err := roleService.Delete(roleID)

		// Assert expectations
		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Role not found", func(t *testing.T) {
		roleID := int64(999)

		// Setup expectations
		mockRepo.On("GetByID", roleID).Return(
			nil,
			errors.New(errors.ErrDatabaseDelete, "record not found"),
		).Once()

		// Call service method
		err := roleService.Delete(roleID)

		// Assert expectations
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "code: 2004")
		mockRepo.AssertExpectations(t)
	})

	t.Run("Delete error", func(t *testing.T) {
		roleID := int64(2)
		role := &models.Role{
			Name:        "manager",
			DisplayName: "Manager",
		}

		// Setup expectations
		mockRepo.On("GetByID", roleID).Return(role, nil).Once()
		mockRepo.On("Delete", role).Return(
			errors.New(errors.ErrDatabaseDelete, "foreign key constraint"),
		).Once()

		// Call service method
		err := roleService.Delete(roleID)

		// Assert expectations
		assert.Error(t, err)
		mockRepo.AssertExpectations(t)
	})
}
