package services_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/vfa-khuongdv/golang-cms/internal/models"
	"github.com/vfa-khuongdv/golang-cms/internal/services"
	"github.com/vfa-khuongdv/golang-cms/pkg/apperror"
	"github.com/vfa-khuongdv/golang-cms/tests/mocks"
	"gorm.io/gorm"
)

type PermissionServiceTestSuite struct {
	suite.Suite
	repo              *mocks.MockPermissionRepository
	permissionService *services.PermissionService
}

func (s *PermissionServiceTestSuite) SetupTest() {
	s.repo = new(mocks.MockPermissionRepository)
	s.permissionService = services.NewPermissionService(s.repo)
}
func (s *PermissionServiceTestSuite) TearDownTest() {
	s.repo.AssertExpectations(s.T())
}

func (s *PermissionServiceTestSuite) TestGetAll_Success() {
	expected := []models.Permission{
		{ID: 1, Resource: "user", Action: "view"},
		{ID: 2, Resource: "user", Action: "edit"},
	}

	s.repo.On("GetAll").Return(expected, nil).Once()

	permissions, err := s.permissionService.GetAll()

	s.NoError(err)
	s.Equal(expected, permissions)
}
func (s *PermissionServiceTestSuite) TestGetAll_Error() {
	s.repo.On("GetAll").Return(([]models.Permission)(nil), gorm.ErrInvalidDB).Once()

	permissions, err := s.permissionService.GetAll()

	// Check if the error is of type apperror.InternalError
	if appErr, ok := err.(*apperror.AppError); ok {
		s.Equal(apperror.ErrInternal, appErr.Code)
		s.Equal(gorm.ErrInvalidDB.Error(), appErr.Message)
	} else {
		s.Fail("Expected apperror.InternalError, got", fmt.Sprintf("%T", err))
	}

	s.Error(err)
	s.Nil(permissions)
}
func TestPermissionServiceTestSuite(t *testing.T) {
	suite.Run(t, new(PermissionServiceTestSuite))
}
