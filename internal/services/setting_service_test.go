package services_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/vfa-khuongdv/golang-cms/internal/models"
	"github.com/vfa-khuongdv/golang-cms/internal/services"
	"github.com/vfa-khuongdv/golang-cms/pkg/apperror"
	"github.com/vfa-khuongdv/golang-cms/tests/mocks"
)

type SettingServiceTestSuite struct {
	suite.Suite
	repo           *mocks.MockSettingRepository
	settingService *services.SettingService
}

func (s *SettingServiceTestSuite) SetupTest() {
	s.repo = new(mocks.MockSettingRepository)
	s.settingService = services.NewSettingService(s.repo)
}

func (s *SettingServiceTestSuite) TearDownTest() {
	s.repo.AssertExpectations(s.T())
}

func (s *SettingServiceTestSuite) TestGetSetting_Success() {
	expected := []models.Setting{
		{SettingKey: "site_name", Value: "My Site"},
		{SettingKey: "site_url", Value: "https://example.com"},
	}

	s.repo.On("GetAll").Return(expected, nil).Once()

	settings, err := s.settingService.GetSetting()

	s.NoError(err)
	s.Equal(expected, settings)
}

func (s *SettingServiceTestSuite) TestGetSetting_Error() {
	s.repo.On("GetAll").Return(([]models.Setting)(nil), apperror.NewDBQueryError("query failed")).Once()

	settings, err := s.settingService.GetSetting()

	s.Error(err)
	s.Contains(err.Error(), fmt.Sprintf("code: %d", apperror.ErrDBQuery))
	s.Nil(settings)
}

func (s *SettingServiceTestSuite) TestGetSettingByKey_Success() {
	expected := &models.Setting{SettingKey: "site_name", Value: "My Site"}

	s.repo.On("GetByKey", "site_name").Return(expected, nil).Once()

	setting, err := s.settingService.GetSettingByKey("site_name")

	s.NoError(err)
	s.Equal(expected, setting)
}

func (s *SettingServiceTestSuite) TestGetSettingByKey_Error() {
	s.repo.On("GetByKey", "non_existent_key").Return((*models.Setting)(nil), apperror.NewNotFoundError("not found")).Once()

	setting, err := s.settingService.GetSettingByKey("non_existent_key")

	s.Error(err)
	s.Nil(setting)
	s.Contains(err.Error(), fmt.Sprintf("code: %d", apperror.ErrNotFound))
}

func (s *SettingServiceTestSuite) TestUpdate_Success() {
	setting := &models.Setting{SettingKey: "site_name", Value: "Updated Site"}

	s.repo.On("Update", setting).Return(nil).Once()

	err := s.settingService.Update(setting)

	s.NoError(err)
}

func (s *SettingServiceTestSuite) TestUpdate_Error() {
	setting := &models.Setting{SettingKey: "site_name", Value: "Updated Site"}

	s.repo.On("Update", setting).Return(apperror.NewDBUpdateError("update failed")).Once()

	err := s.settingService.Update(setting)

	s.Error(err)
	s.Contains(err.Error(), fmt.Sprintf("code: %d", apperror.ErrDBUpdate))

}

func (s *SettingServiceTestSuite) TestCreate_Success() {
	setting := &models.Setting{SettingKey: "site_name", Value: "New Site"}

	s.repo.On("Create", setting).Return(nil).Once()

	err := s.settingService.Create(setting)

	s.NoError(err)
}

func (s *SettingServiceTestSuite) TestCreate_Error() {
	setting := &models.Setting{SettingKey: "site_name", Value: "New Site"}

	s.repo.On("Create", setting).Return(apperror.NewDBInsertError("insert failed")).Once()

	err := s.settingService.Create(setting)

	s.Error(err)
	s.Contains(err.Error(), fmt.Sprintf("code: %d", apperror.ErrDBInsert))
}

func TestSettingServiceTestSuite(t *testing.T) {
	suite.Run(t, new(SettingServiceTestSuite))
}
