package tests_internal_repositories

import (
	"testing"

	"github.com/stretchr/testify/assert"
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

// TestAssignPermissions tests assigning permissions to a role
func (s *RoleRepositoryTestSuite) TestAssignPermissions() {
	// Create a role
	role := &models.Role{
		Name:        "permission_tester",
		DisplayName: "Permission Tester",
	}
	err := s.repo.Create(role)
	s.NoError(err)

	// Create some permissions
	permissions := []models.Permission{
		{Resource: "users", Action: "create"},
		{Resource: "users", Action: "read"},
		{Resource: "users", Action: "update"},
		{Resource: "users", Action: "delete"},
	}

	for i := range permissions {
		err := s.db.Create(&permissions[i]).Error
		s.NoError(err)
		s.NotEqual(uint(0), permissions[i].ID)
	}

	// Get permission IDs
	permissionIDs := []uint{
		permissions[0].ID,
		permissions[1].ID,
		permissions[2].ID,
	}

	// Assign permissions to the role
	err = s.repo.AssignPermissions(role.ID, permissionIDs)
	s.NoError(err)

	// Verify assigned permissions
	rolePermissions, err := s.repo.GetRolePermissions(role.ID)
	s.NoError(err)
	s.Equal(3, len(rolePermissions), "Should have 3 permissions assigned")

	// Update permissions (remove one, add another)
	newPermissionIDs := []uint{
		permissions[0].ID,
		permissions[2].ID,
		permissions[3].ID, // New permission
	}

	err = s.repo.AssignPermissions(role.ID, newPermissionIDs)
	s.NoError(err)

	// Verify updated permissions
	rolePermissions, err = s.repo.GetRolePermissions(role.ID)
	s.NoError(err)
	s.Equal(3, len(rolePermissions), "Should still have 3 permissions but with different composition")

	// Verify correct permissions are assigned
	hasPermission := func(perms []models.Permission, resource, action string) bool {
		for _, p := range perms {
			if p.Resource == resource && p.Action == action {
				return true
			}
		}
		return false
	}

	s.True(hasPermission(rolePermissions, "users", "create"))
	s.True(hasPermission(rolePermissions, "users", "update"))
	s.True(hasPermission(rolePermissions, "users", "delete"))
	s.False(hasPermission(rolePermissions, "users", "read"), "Read permission should have been removed")
}

// TestGetRolePermissions tests retrieving permissions assigned to a role
func (s *RoleRepositoryTestSuite) TestGetRolePermissions() {
	// Create a role
	role := &models.Role{
		Name:        "permission_viewer",
		DisplayName: "Permission Viewer",
	}
	err := s.repo.Create(role)
	s.NoError(err)

	// Create some permissions
	permissions := []models.Permission{
		{Resource: "posts", Action: "create"},
		{Resource: "posts", Action: "read"},
	}

	for i := range permissions {
		err := s.db.Create(&permissions[i]).Error
		s.NoError(err)
	}

	// Initially, role should have no permissions
	rolePermissions, err := s.repo.GetRolePermissions(role.ID)
	s.NoError(err)
	s.Equal(0, len(rolePermissions), "New role should have no permissions")

	// Assign permissions
	permissionIDs := []uint{permissions[0].ID, permissions[1].ID}
	err = s.repo.AssignPermissions(role.ID, permissionIDs)
	s.NoError(err)

	// Now role should have 2 permissions
	rolePermissions, err = s.repo.GetRolePermissions(role.ID)
	s.NoError(err)
	s.Equal(2, len(rolePermissions), "Role should have 2 permissions")

	// Verify correct permissions are retrieved
	s.Equal("posts", rolePermissions[0].Resource)
	s.Equal("create", rolePermissions[0].Action)
	s.Equal("posts", rolePermissions[1].Resource)
	s.Equal("read", rolePermissions[1].Action)
}

// TestRoleRepository runs the role repository test suite
func TestRoleRepository(t *testing.T) {
	suite.Run(t, new(RoleRepositoryTestSuite))
}

// Additional tests for edge cases

// TestRoleRepository_Calls tests basic repository calls without assertions
// This is similar to the existing role_repository2_test.go
func TestRoleRepository_Calls(t *testing.T) {
	db, _ := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	_ = db.AutoMigrate(&models.Role{}, &models.Permission{}, &models.RolePermission{})
	repo := repositories.NewRoleRepository(db)

	// Test Create
	role := &models.Role{Name: "TestRole", DisplayName: "Test Role"}
	err := repo.Create(role)
	assert.NoError(t, err)

	// Test GetByID
	retrievedRole, err := repo.GetByID(int64(role.ID))
	assert.NoError(t, err)
	assert.Equal(t, role.Name, retrievedRole.Name)

	// Test Update
	role.Name = "UpdatedRole"
	err = repo.Update(role)
	assert.NoError(t, err)

	// Get the updated role and verify
	updatedRole, err := repo.GetByID(int64(role.ID))
	assert.NoError(t, err)
	assert.Equal(t, "UpdatedRole", updatedRole.Name)

	// Test GetRolePermissions (should return empty, no panic)
	permissions, err := repo.GetRolePermissions(role.ID)
	assert.NoError(t, err)
	assert.Equal(t, 0, len(permissions))

	// Test AssignPermissions with empty permission IDs
	err = repo.AssignPermissions(role.ID, []uint{})
	assert.NoError(t, err)

	// Create a permission to test assigning
	permission := &models.Permission{Resource: "test", Action: "test"}
	err = db.Create(permission).Error
	assert.NoError(t, err)

	// Test AssignPermissions with actual permission
	err = repo.AssignPermissions(role.ID, []uint{permission.ID})
	assert.NoError(t, err)

	// Verify permission was assigned
	permissions, err = repo.GetRolePermissions(role.ID)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(permissions))
	assert.Equal(t, "test", permissions[0].Resource)

	// Test Delete
	err = repo.Delete(role)
	assert.NoError(t, err)

	// Verify role was deleted
	_, err = repo.GetByID(int64(role.ID))
	assert.Error(t, err)
}
