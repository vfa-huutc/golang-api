package services_test

import (
	originErrors "errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/vfa-khuongdv/golang-cms/internal/models"
	"github.com/vfa-khuongdv/golang-cms/internal/services"
	"github.com/vfa-khuongdv/golang-cms/pkg/errors"
	"github.com/vfa-khuongdv/golang-cms/tests/mocks"
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
	s.repo.On("GetAll").Return(([]models.Permission)(nil), originErrors.New("query failed")).Once()

	permissions, err := s.permissionService.GetAll()

	s.Error(err)
	s.Contains(err.Error(), fmt.Sprintf("code: %d", errors.ErrDBQuery))
	s.Nil(permissions)
}
func TestPermissionServiceTestSuite(t *testing.T) {
	suite.Run(t, new(PermissionServiceTestSuite))
}
