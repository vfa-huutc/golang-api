package repositories_test

import (
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/vfa-khuongdv/golang-cms/internal/models"
	"github.com/vfa-khuongdv/golang-cms/internal/repositories"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// SettingRepositoryTestSuite defines a test suite for the setting repository
type SettingRepositoryTestSuite struct {
	suite.Suite
	db   *gorm.DB
	repo *repositories.SettingRepository
}

// SetupTest prepares the database and repository for each test
func (s *SettingRepositoryTestSuite) SetupTest() {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	s.Require().NoError(err)
	s.Require().NotNil(db)

	// Auto-migrate the models
	err = db.AutoMigrate(&models.Setting{})
	s.Require().NoError(err)

	s.db = db
	s.repo = repositories.NewSettingRepository(db)
}

func (s *SettingRepositoryTestSuite) TearDownTest() {
	// Cleanup: drop tables
	db, err := s.db.DB()
	if err == nil {
		db.Close()
	}
}

func (s *SettingRepositoryTestSuite) TestGetAll() {
	mockSettings := []models.Setting{
		{SettingKey: "site_name", Value: "Test Site"},
		{SettingKey: "site_url", Value: "https://testsite.com"},
	}

	// Insert mock data into the database
	for _, setting := range mockSettings {
		err := s.repo.Create(&setting)
		s.NoError(err)
	}

	// Test GetAll method
	settings, err := s.repo.GetAll()
	s.NoError(err)
	s.Len(settings, len(mockSettings), "Should return all settings")
	for i, setting := range settings {
		s.Equal(mockSettings[i].SettingKey, setting.SettingKey, "Setting key should match")
		s.Equal(mockSettings[i].Value, setting.Value, "Setting value should match")
	}
}

func (s *SettingRepositoryTestSuite) TestGetAllError() {
	// Close the underlying DB connection to simulate error on DB access
	sqlDB, err := s.db.DB()
	s.Require().NoError(err)
	err = sqlDB.Close()
	s.Require().NoError(err)

	// Now s.db is still there but the connection is closed, this should cause errors

	// No need to create new repo, because it uses s.db which now has closed connection
	settings, err := s.repo.GetAll()
	s.Error(err, "Should return an error when DB connection is closed")
	s.Nil(settings, "Settings should be nil on error")
}
func (s *SettingRepositoryTestSuite) TestGetByKey() {
	mockSetting := models.Setting{
		SettingKey: "site_name",
		Value:      "Test Site",
	}

	// Insert mock data into the database
	err := s.repo.Create(&mockSetting)
	s.NoError(err)

	// Test GetByKey method
	setting, err := s.repo.GetByKey("site_name")
	s.NoError(err)
	s.NotNil(setting, "Should return a setting")
	s.Equal(mockSetting.SettingKey, setting.SettingKey, "Setting key should match")
	s.Equal(mockSetting.Value, setting.Value, "Setting value should match")

	// Test non-existing key
	setting, err = s.repo.GetByKey("non_existing_key")
	s.Error(err, "Should return an error for non-existing key")
	s.Nil(setting, "Should return nil for non-existing key")
}

func (s *SettingRepositoryTestSuite) TestGetByKeyError() {
	// Close the underlying DB connection to simulate error on DB access
	sqlDB, err := s.db.DB()
	s.Require().NoError(err)
	err = sqlDB.Close()
	s.Require().NoError(err)

	// Now s.db is still there but the connection is closed, this should cause errors

	// No need to create new repo, because it uses s.db which now has closed connection
	setting, err := s.repo.GetByKey("site_name")
	s.Error(err, "Should return an error when DB connection is closed")
	s.Nil(setting, "Setting should be nil on error")
}

func (s *SettingRepositoryTestSuite) TestUpdate() {
	mockSetting := &models.Setting{
		SettingKey: "site_name",
		Value:      "Test Site",
	}

	// Insert mock data into the database
	err := s.repo.Create(mockSetting)
	s.NoError(err)

	// Update the setting
	mockSetting.Value = "Updated Site Name"
	err = s.repo.Update(mockSetting)
	s.NoError(err)

	// Retrieve the updated setting
	updatedSetting, err := s.repo.GetByKey("site_name")
	s.NoError(err)
	s.NotNil(updatedSetting, "Should return an updated setting")
	s.Equal("Updated Site Name", updatedSetting.Value, "Updated value should match")
}

func (s *SettingRepositoryTestSuite) TestUpdateError() {
	// Close the underlying DB connection to simulate error on DB access
	sqlDB, err := s.db.DB()
	s.Require().NoError(err)
	err = sqlDB.Close()
	s.Require().NoError(err)

	// Now s.db is still there but the connection is closed, this should cause errors

	// No need to create new repo, because it uses s.db which now has closed connection
	mockSetting := &models.Setting{
		SettingKey: "site_name",
		Value:      "Test Site",
	}
	err = s.repo.Update(mockSetting)
	s.Error(err, "Should return an error when DB connection is closed")
}

func (s *SettingRepositoryTestSuite) TestCreate() {
	mockSetting := &models.Setting{
		SettingKey: "site_name",
		Value:      "Test Site",
	}

	// Test successful creation
	err := s.repo.Create(mockSetting)
	s.NoError(err)
	s.NotEqual(uint(0), mockSetting.ID, "Setting ID should be set after creation")

	// Test duplicate key creation
	err = s.repo.Create(mockSetting)
	s.Error(err, "Should return an error for duplicate setting key")
}

func (s *SettingRepositoryTestSuite) TestCreateError() {
	// Close the underlying DB connection to simulate error on DB access
	sqlDB, err := s.db.DB()
	s.Require().NoError(err)
	err = sqlDB.Close()
	s.Require().NoError(err)

	// Now s.db is still there but the connection is closed, this should cause errors

	// No need to create new repo, because it uses s.db which now has closed connection
	mockSetting := &models.Setting{
		SettingKey: "site_name",
		Value:      "Test Site",
	}
	err = s.repo.Create(mockSetting)
	s.Error(err, "Should return an error when DB connection is closed")
}

func TestSettingRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(SettingRepositoryTestSuite))
}
