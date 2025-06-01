package repositories_test

import (
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/vfa-khuongdv/golang-cms/internal/models"
	"github.com/vfa-khuongdv/golang-cms/internal/repositories"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// RoleRepositoryTestSuite defines a test suite for the role repository
type RoleRepositoryTestSuite struct {
	suite.Suite
	db   *gorm.DB
	repo *repositories.RoleRepository
}

// SetupTest prepares the database and repository for each test
func (s *RoleRepositoryTestSuite) SetupTest() {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	s.Require().NoError(err)
	s.Require().NotNil(db)

	// Auto-migrate the models
	err = db.AutoMigrate(&models.Role{}, &models.Permission{}, &models.RolePermission{})
	s.Require().NoError(err)

	s.db = db
	s.repo = repositories.NewRoleRepository(db)
}

// TearDownTest cleans up after each test
func (s *RoleRepositoryTestSuite) TearDownTest() {
	// Cleanup: drop tables
	db, err := s.db.DB()
	if err == nil {
		db.Close()
	}
}

// TestCreate tests the role creation functionality
func (s *RoleRepositoryTestSuite) TestCreate() {
	role := &models.Role{
		Name:        "admin",
		DisplayName: "Administrator",
	}

	// Test successful creation
	err := s.repo.Create(role)
	s.NoError(err)
	s.NotEqual(uint(0), role.ID, "Role ID should be set after creation")

	// Test unique constraint violation
	duplicateRole := &models.Role{
		Name:        "admin", // Same name as the previously created role
		DisplayName: "Admin Role",
	}
	err = s.repo.Create(duplicateRole)
	s.Error(err, "Expected error when creating a role with duplicate name")
}

// TestGetByID tests retrieving a role by ID
func (s *RoleRepositoryTestSuite) TestGetByID() {
	// Create a role first
	role := &models.Role{
		Name:        "editor",
		DisplayName: "Editor Role",
	}
	err := s.repo.Create(role)
	s.NoError(err)
	s.NotEqual(uint(0), role.ID)

	// Test retrieving the role
	retrievedRole, err := s.repo.GetByID(int64(role.ID))
	s.NoError(err)
	s.NotNil(retrievedRole)
	s.Equal(role.ID, retrievedRole.ID)
	s.Equal("editor", retrievedRole.Name)
	s.Equal("Editor Role", retrievedRole.DisplayName)

	// Test retrieving non-existent role
	_, err = s.repo.GetByID(9999)
	s.Error(err, "Expected error when retrieving non-existent role")
}

// TestUpdate tests updating a role
func (s *RoleRepositoryTestSuite) TestUpdate() {
	// Create a role first
	role := &models.Role{
		Name:        "moderator",
		DisplayName: "Moderator",
	}
	err := s.repo.Create(role)
	s.NoError(err)

	// Update the role
	role.DisplayName = "Updated Moderator"
	role.Name = "updated_moderator"
	err = s.repo.Update(role)
	s.NoError(err)

	// Verify the update
	updatedRole, err := s.repo.GetByID(int64(role.ID))
	s.NoError(err)
	s.Equal("updated_moderator", updatedRole.Name)
	s.Equal("Updated Moderator", updatedRole.DisplayName)
}

// TestDelete tests deleting a role
func (s *RoleRepositoryTestSuite) TestDelete() {
	// Create a role first
	role := &models.Role{
		Name:        "guest",
		DisplayName: "Guest User",
	}
	err := s.repo.Create(role)
	s.NoError(err)

	// Delete the role
	err = s.repo.Delete(role)
	s.NoError(err)

	// Verify deletion
	_, err = s.repo.GetByID(int64(role.ID))
	s.Error(err, "Expected error when retrieving deleted role")
}

func TestRoleRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(RoleRepositoryTestSuite))
}
