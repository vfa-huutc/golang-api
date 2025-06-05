package services_test

import (
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/vfa-khuongdv/golang-cms/internal/models"
	"github.com/vfa-khuongdv/golang-cms/internal/services"
	"github.com/vfa-khuongdv/golang-cms/pkg/apperror"
	"github.com/vfa-khuongdv/golang-cms/tests/mocks"
	"gorm.io/gorm"
)

type RoleServiceTestSuite struct {
	suite.Suite
	repo        *mocks.MockRoleRepository
	roleService *services.RoleService
}

func (s *RoleServiceTestSuite) SetupTest() {
	s.repo = new(mocks.MockRoleRepository)
	s.roleService = services.NewRoleService(s.repo)
}

func (s *RoleServiceTestSuite) TestGetByID_Success() {
	expected := &models.Role{Name: "admin", DisplayName: "Administrator"}
	s.repo.On("GetByID", int64(1)).Return(expected, nil).Once()

	role, err := s.roleService.GetByID(1)

	s.NoError(err)
	s.Equal(expected, role)
	s.repo.AssertExpectations(s.T())
}

func (s *RoleServiceTestSuite) TestGetByID_NotFound() {
	s.repo.On("GetByID", int64(999)).Return((*models.Role)(nil), apperror.NewDBQueryError("Not found resource")).Once()

	role, err := s.roleService.GetByID(999)

	s.Error(err)
	s.Nil(role)
	s.Contains(err.Error(), "code: 2001")
	s.repo.AssertExpectations(s.T())
}

func (s *RoleServiceTestSuite) TestCreate_Success() {
	role := &models.Role{Name: "editor", DisplayName: "Content Editor"}
	s.repo.On("Create", role).Return(nil).Once()

	err := s.roleService.Create(role)

	s.NoError(err)
	s.repo.AssertExpectations(s.T())
}

func (s *RoleServiceTestSuite) TestCreate_Error() {
	role := &models.Role{Name: "existing_role", DisplayName: "Existing Role"}
	s.repo.On("Create", role).Return(apperror.NewDBInsertError("duplicate entry")).Once()

	err := s.roleService.Create(role)

	s.Error(err)
	s.Contains(err.Error(), "code: 2002")
	s.repo.AssertExpectations(s.T())
}

func (s *RoleServiceTestSuite) TestUpdate_Success() {
	role := &models.Role{Name: "moderator", DisplayName: "Content Moderator"}
	s.repo.On("Update", role).Return(nil).Once()

	err := s.roleService.Update(role)

	s.NoError(err)
	s.repo.AssertExpectations(s.T())
}

func (s *RoleServiceTestSuite) TestUpdate_Error() {
	role := &models.Role{Name: "invalid_role", DisplayName: "Invalid Role"}
	s.repo.On("Update", role).Return(apperror.NewDBUpdateError("record not found")).Once()

	err := s.roleService.Update(role)

	s.Error(err)
	s.Contains(err.Error(), "code: 2003")
	s.repo.AssertExpectations(s.T())
}

func (s *RoleServiceTestSuite) TestDelete_Success() {
	roleID := int64(1)
	role := &models.Role{Name: "guest", DisplayName: "Guest User"}
	s.repo.On("GetByID", roleID).Return(role, nil).Once()
	s.repo.On("Delete", role).Return(nil).Once()

	err := s.roleService.Delete(roleID)

	s.NoError(err)
	s.repo.AssertExpectations(s.T())
}

func (s *RoleServiceTestSuite) TestDelete_RoleNotFound() {
	roleID := int64(999)
	s.repo.On("GetByID", roleID).Return((*models.Role)(nil), gorm.ErrRecordNotFound).Once()

	err := s.roleService.Delete(roleID)

	s.Error(err)
	if appErr, ok := err.(*apperror.AppError); ok {
		s.Equal(apperror.ErrDBQuery, appErr.Code)
		s.Equal(gorm.ErrRecordNotFound.Error(), appErr.Message)
	} else {
		s.Fail("Expected apperror.AppError, got", err)
	}
	s.repo.AssertExpectations(s.T())
}

func (s *RoleServiceTestSuite) TestDelete_CannotDelete() {
	roleID := int64(2)
	role := &models.Role{Name: "manager", DisplayName: "Manager"}
	s.repo.On("GetByID", roleID).Return(role, nil).Once()
	s.repo.On("Delete", role).Return(gorm.ErrForeignKeyViolated).Once()

	err := s.roleService.Delete(roleID)

	if appErr, ok := err.(*apperror.AppError); ok {
		s.Equal(apperror.ErrDBDelete, appErr.Code)
		s.Equal(gorm.ErrForeignKeyViolated.Error(), appErr.Message)
	} else {
		s.Fail("Expected apperror.AppError, got", err)
	}
	s.repo.AssertExpectations(s.T())
}

func (s *RoleServiceTestSuite) TestDelete_DeleteError() {
	roleID := int64(2)
	role := &models.Role{Name: "manager", DisplayName: "Manager"}
	s.repo.On("GetByID", roleID).Return(role, nil).Once()
	s.repo.On("Delete", role).Return(apperror.NewDBDeleteError("foreign key constraint")).Once()

	err := s.roleService.Delete(roleID)

	s.Error(err)
	s.repo.AssertExpectations(s.T())
}

func TestRoleServiceTestSuite(t *testing.T) {
	suite.Run(t, new(RoleServiceTestSuite))
}
