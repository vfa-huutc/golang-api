package repositories_test

import (
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/vfa-khuongdv/golang-cms/internal/models"
	"github.com/vfa-khuongdv/golang-cms/internal/repositories"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type PermissionRepositoryTestSuite struct {
	suite.Suite
	db   *gorm.DB
	repo *repositories.PermissionRepository
}

func (s *PermissionRepositoryTestSuite) SetupTest() {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	s.Require().NoError(err)
	s.Require().NotNil(db)

	// Auto-migrate the models
	err = db.AutoMigrate(&models.Permission{})
	s.Require().NoError(err)
	s.db = db
	s.repo = repositories.NewPermissionRepository(db)
}

func (s *PermissionRepositoryTestSuite) TearDownTest() {
	db, err := s.db.DB()
	if err == nil {
		db.Close()
	}
}

func (s *PermissionRepositoryTestSuite) TestCreate() {
	permission := &models.Permission{
		Resource: "test_resource",
		Action:   "test_action",
	}
	// Test successful creation
	err := s.repo.Create(permission)
	s.NoError(err)
	s.NotEqual(uint(0), permission.ID, "Permission ID should be set after creation")
}

func (s *PermissionRepositoryTestSuite) TestGetAll() {

	mock_permissions := []models.Permission{
		{
			Resource: "resource1",
			Action:   "action1",
		},
		{
			Resource: "resource2",
			Action:   "action2",
		},
	}

	for _, perm := range mock_permissions {
		err := s.repo.Create(&perm)
		s.NoError(err)
	}

	// Now test GetAll
	permissions, err := s.repo.GetAll()
	s.NoError(err)
	s.NotEmpty(permissions, "Permissions should not be empty")
	s.Len(permissions, 2, "There should be exactly two permissions")

	for index, perm := range permissions {
		s.NotEmpty(perm.Resource, "Resource should not be empty")
		s.NotEmpty(perm.Action, "Action should not be empty")
		s.Equal(perm.Resource, mock_permissions[index].Resource, "Resource should match")
		s.Equal(perm.Action, mock_permissions[index].Action, "Action should match")
	}
}

func TestPermissionRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(PermissionRepositoryTestSuite))
}
