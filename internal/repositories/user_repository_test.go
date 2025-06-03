package repositories_test

import (
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/vfa-khuongdv/golang-cms/internal/models"
	"github.com/vfa-khuongdv/golang-cms/internal/repositories"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type UserRepositoryTestSuite struct {
	suite.Suite
	db   *gorm.DB
	repo *repositories.UserRepository
}

func (s *UserRepositoryTestSuite) SetupTest() {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	s.Require().NoError(err)
	s.Require().NotNil(db)

	// Auto-migrate the models
	err = db.AutoMigrate(
		&models.User{},
		&models.Role{},
		&models.Permission{},
		&models.UserRole{},
	)
	s.Require().NoError(err)
	s.db = db
	s.repo = repositories.NewUserRepository(db)
}

func (s *UserRepositoryTestSuite) TearDownTest() {
	db, err := s.db.DB()
	if err == nil {
		db.Close()
	}
}

func (s *UserRepositoryTestSuite) TestPaginateUser() {
	mockUsers := []*models.User{
		{Name: "User1", Email: "email1@example.com", Password: "password1", Gender: 1},
		{Name: "User2", Email: "email2@example.com", Password: "password2", Gender: 1},
		{Name: "User3", Email: "email3@example.com", Password: "password3", Gender: 1},
		{Name: "User4", Email: "email4@example.com", Password: "password4", Gender: 1},
		{Name: "User5", Email: "email5@example.com", Password: "password5", Gender: 1},
		{Name: "User6", Email: "email6@example.com", Password: "password6", Gender: 1},
		{Name: "User7", Email: "email7@example.com", Password: "password7", Gender: 1},
		{Name: "User8", Email: "email8@example.com", Password: "password8", Gender: 1},
		{Name: "User9", Email: "email9@example.com", Password: "password9", Gender: 1},
		{Name: "User10", Email: "email10@example.com", Password: "password10", Gender: 1},
	}
	for _, user := range mockUsers {
		_, err := s.repo.Create(user)
		s.NoError(err, "Expected no error when creating mock users")
	}
	pagination, err := s.repo.PaginateUser(1, 5)
	s.NoError(err, "Expected no error when paginating users")
	s.NotNil(pagination, "Expected pagination to be not nil")
	s.Equal(1, pagination.Page, "Expected page number to be 1")
	s.Equal(5, pagination.Limit, "Expected limit to be 5")
	s.Equal(10, pagination.TotalItems, "Expected total items to be 10")
	s.Len(pagination.Data, 5, "Expected 5 items in the first page of pagination")

	// Parse the data to []models.User
	users, ok := pagination.Data.([]models.User)
	s.True(ok, "Expected pagination data to be of type []models.User")

	//s.T().Logf("Users in first page: %+v", users)
	// Expect the first 5 users at page 1 to be returned with correct names
	s.Equal("User10", users[0].Name)
	s.Equal("User9", users[1].Name)
	s.Equal("User8", users[2].Name)
	s.Equal("User7", users[3].Name)
	s.Equal("User6", users[4].Name)

}

func (s *UserRepositoryTestSuite) TestGetAll() {
	mockUsers := []*models.User{
		{Name: "User1", Email: "email1@example.com", Password: "password1", Gender: 1},
		{Name: "User2", Email: "email2@example.com", Password: "password2", Gender: 1},
	}
	for _, user := range mockUsers {
		_, err := s.repo.Create(user)
		s.NoError(err, "Expected no error when creating mock users")
	}
	users, err := s.repo.GetAll()
	s.NoError(err, "Expected no error when getting all users")
	s.Len(users, 2, "Expected 2 users to be returned")
}

func (s *UserRepositoryTestSuite) TestGetAllError() {
	// Close the underlying DB connection to simulate error on DB access
	sqlDB, err := s.db.DB()
	s.Require().NoError(err)
	err = sqlDB.Close()
	s.Require().NoError(err)

	// Test GetAll method after closing the DB
	users, err := s.repo.GetAll()
	s.Error(err, "Expected error when getting all users after closing DB")
	s.Nil(users, "Expected users to be nil after error")
}

func (s *UserRepositoryTestSuite) TestGetByID() {
	mockUsers := []*models.User{
		{ID: 1, Name: "User1", Email: "email1@example.com", Password: "password1", Gender: 1},
		{ID: 2, Name: "User2", Email: "email2@example.com", Password: "password2", Gender: 1},
	}
	for _, user := range mockUsers {
		_, err := s.repo.Create(user)
		s.NoError(err, "Expected no error when creating mock user")
	}
	// Test getting user by ID
	user, err := s.repo.GetByID(1)
	s.NoError(err, "Expected no error when getting user by ID")
	s.NotNil(user, "Expected user to be not nil")
	s.Equal("User1", user.Name, "Expected user name to be 'User1'")
}

func (s *UserRepositoryTestSuite) TestGetByIDError() {
	// Test getting user by non-existing ID
	user, err := s.repo.GetByID(999)
	s.Error(err, "Expected error when getting user by non-existing ID")
	s.Nil(user, "Expected user to be nil when not found")
}

func (s *UserRepositoryTestSuite) TestCreate() {
	mockUser := &models.User{
		Name:     "New User",
		Email:    "email@example.com",
		Password: "password",
		Gender:   1,
	}
	createdUser, err := s.repo.Create(mockUser)
	s.NoError(err, "Expected no error when creating user")
	s.NotNil(createdUser, "Expected created user to be not nil")
	s.Equal("New User", createdUser.Name, "Expected user name to be 'New User'")
}

func (s *UserRepositoryTestSuite) TestCreate_Error_Duplicate() {
	user1 := &models.User{
		Email: "test@example.com",
		Name:  "testuser",
		// other required fields
	}

	user2 := &models.User{
		Email: "test@example.com", // duplicate email to cause constraint violation
		Name:  "anotheruser",
	}

	// Create first user successfully
	createdUser, err := s.repo.Create(user1)
	s.Require().NoError(err)
	s.NotNil(createdUser)

	// Try to create second user with duplicate email
	createdUser2, err := s.repo.Create(user2)
	s.Error(err, "Should return error on duplicate user creation")
	s.Nil(createdUser2, "Created user should be nil on error")
}

func (s *UserRepositoryTestSuite) TestDeleteError() {
	// Close the underlying DB connection to simulate error on DB access
	sqlDB, err := s.db.DB()
	s.Require().NoError(err)
	err = sqlDB.Close()
	s.Require().NoError(err)

	// Now s.db is still there but the connection is closed, this should cause errors

	// No need to create new repo, because it uses s.db which now has closed connection
	err = s.repo.Delete(999)
	s.Error(err, "Expected error when deleting user with non-existing ID")
}

func (s *UserRepositoryTestSuite) TestFindByField() {
	mockUsers := []*models.User{
		{Name: "Find User", Email: "email@example.com", Password: "password", Gender: 1},
		{Name: "Another User", Email: "another@example.com", Password: "password", Gender: 1},
	}

	for _, user := range mockUsers {
		_, err := s.repo.Create(user)
		s.NoError(err, "Expected no error when creating mock users")
	}

	// find by field email
	foundUser, err := s.repo.FindByField("email", "email@example.com")
	s.NoError(err, "Expected no error when finding user by email")
	s.NotNil(foundUser, "Expected found user to be not nil")
	s.Equal("Find User", foundUser.Name, "Expected user name to be 'Find User'")
	// find by field name
	foundUserByName, err := s.repo.FindByField("name", "Another User")
	s.NoError(err, "Expected no error when finding user by name")
	s.NotNil(foundUserByName, "Expected found user by name to be not nil")
	s.Equal("Another User", foundUserByName.Name, "Expected user name to be 'Another User'")
	// Test finding user by non-existing field
	nonExistentUser, err := s.repo.FindByField("email", "notfound@example.com")
	s.Error(err, "Expected error when finding user by non-existing email")
	s.Nil(nonExistentUser, "Expected non-existent user to be nil")
}

func (s *UserRepositoryTestSuite) TestFindByFieldError() {
	// Test finding user by non-existing field
	nonExistentUser, err := s.repo.FindByField("email", "notfound@example.com")
	s.Error(err, "Expected error when finding user by non-existing email")
	s.Nil(nonExistentUser, "Expected non-existent user to be nil")
}

func (s *UserRepositoryTestSuite) TestGetProfile() {
	type MockUser struct {
		users *models.User
		roles []*models.Role
	}
	mockUsers := []MockUser{
		{
			users: &models.User{
				Name:     "Profile User",
				Email:    "profile@example.com",
				Password: "password",
				Gender:   1,
			},
			roles: []*models.Role{
				{Name: "Admin", DisplayName: "Administrator"},
				{Name: "Editor", DisplayName: "Content Editor"},
			},
		},
	}

	// Create mock users and roles
	for _, mock := range mockUsers {
		createdUser, err := s.repo.Create(mock.users)
		s.NoError(err, "Expected no error when creating mock user")
		s.NotNil(createdUser, "Expected created user to be not nil")

		for _, role := range mock.roles {
			role.ID = 0 // Reset ID for new role creation
			err = s.repo.GetDB().Create(role).Error
			s.NoError(err, "Expected no error when creating mock role")
		}
		// Assign roles to user
		s.repo.GetDB().Model(createdUser).Association("Roles").Append(mock.roles)
	}
	profile, err := s.repo.GetProfile(mockUsers[0].users.ID)

	s.NoError(err, "Expected no error when getting user profile")
	s.NotNil(profile, "Expected user profile to be not nil")
	s.Equal("Profile User", profile.Name, "Expected user name to be 'Profile User'")

	// Check if roles are assigned correctly
	roleNames := make([]string, len(profile.Roles))
	for i, role := range profile.Roles {
		roleNames[i] = role.Name
	}

	s.ElementsMatch([]string{"Admin", "Editor"}, roleNames, "Expected user roles to match")

}

func (s *UserRepositoryTestSuite) TestGetProfileError() {
	// Test getting profile for non-existing user
	profile, err := s.repo.GetProfile(999)
	s.Error(err, "Expected error when getting profile for non-existing user")
	s.Nil(profile, "Expected profile to be nil when user does not exist")
}

func (s *UserRepositoryTestSuite) TestUpdateProfile() {
	mockUser := &models.User{
		ID:       1,
		Name:     "Profile User",
		Email:    "email@example.com",
		Password: "password",
		Gender:   1,
	}
	// 1. Create a mock user
	createdUser, err := s.repo.Create(mockUser)
	s.NoError(err, "Expected no error when creating mock user")
	s.NotNil(createdUser, "Expected created user to be not nil")

	// 2. Update the user profile
	createdUser.Name = "Updated Profile User"
	err = s.repo.UpdateProfile(createdUser)
	s.NoError(err, "Expected no error when updating user profile")

	// 3. Retrieve the updated user profile
	updatedUser, err := s.repo.GetByID(createdUser.ID)
	s.NoError(err, "Expected no error when getting updated user by ID")
	s.NotNil(updatedUser, "Expected updated user to be not nil")
	s.Equal("Updated Profile User", updatedUser.Name, "Expected user name to be 'Updated Profile User'")

}

func (s *UserRepositoryTestSuite) TestGetUserPermissions() {

	// Mock data for roles and permissions
	type MasterWithRole struct {
		Roles       *models.Role
		Permissions []models.Permission
	}
	// Mock user with roles and permissions
	type MockUser struct {
		user      *models.User
		userRoles []models.UserRole
	}
	// Create mock roles and permissions
	mockRoles := []MasterWithRole{
		{
			Roles: &models.Role{ID: 1, Name: "Admin", DisplayName: "Administrator"},
			Permissions: []models.Permission{
				{Resource: "create_user", Action: "Create User"},
				{Resource: "delete_user", Action: "Delete User"},
				{Resource: "update_user", Action: "Update User"},
				{Resource: "view_user", Action: "View User"},
			},
		},
	}
	// Create mock users with roles
	mockUsers := []MockUser{
		{
			user: &models.User{Name: "Bob", Email: "bob@example.com", Password: "password", Gender: 1},
			userRoles: []models.UserRole{
				{ID: 1, RoleID: 1, UserID: 1},
			},
		},
		{
			user:      &models.User{Name: "John", Email: "john@example.com", Password: "password", Gender: 1},
			userRoles: []models.UserRole{},
		},
	}

	// Create mock roles and permissions in the database
	for _, mock := range mockRoles {
		err := s.repo.GetDB().Create(mock.Roles).Error
		s.NoError(err, "Expected no error when creating mock role")
		for _, perm := range mock.Permissions {
			perm.ID = 0 // Reset ID for new permission creation
			err = s.repo.GetDB().Create(&perm).Error
			s.NoError(err, "Expected no error when creating mock permission")
			// Associate permission with role
			err = s.repo.GetDB().Model(mock.Roles).Association("Permissions").Append(&perm)
			s.NoError(err, "Expected no error when associating permission with role")
		}
	}
	// Create mock users and assign roles
	for _, mock := range mockUsers {
		createdUser, err := s.repo.Create(mock.user)
		s.NoError(err, "Expected no error when creating mock user")
		s.NotNil(createdUser, "Expected created user to be not nil")

		for _, userRole := range mock.userRoles {
			err = s.repo.GetDB().Create(&userRole).Error
			s.NoError(err, "Expected no error when creating mock user role")
		}
	}
	// Test bob's with with roles and permissions
	permissions, err := s.repo.GetUserPermissions(mockUsers[0].user.ID)
	s.NoError(err, "Expected no error when getting user permissions")
	s.Len(permissions, 4, "Expected 4 permissions for the user")

	// Test John's with no roles and permissions
	permissions, err = s.repo.GetUserPermissions(mockUsers[1].user.ID)
	s.NoError(err, "Expected no error when getting user permissions")
	s.Len(permissions, 0, "Expected no permissions for user without roles")
}

func (s *UserRepositoryTestSuite) TestGetUserPermissionsError() {
	// Test with invalid user ID
	permissions, err := s.repo.GetUserPermissions(999)
	s.Error(err, "Expected error when getting permissions for non-existent user")
	s.Nil(permissions, "Expected nil permissions for non-existent user")
}

func (s *UserRepositoryTestSuite) TestCreateWithTx() {
	mockUser := &models.User{
		Name:  "Transaction User",
		Email: "email@example.com",
	}
	// Start a transaction
	tx := s.repo.GetDB().Begin()
	s.Require().NoError(tx.Error, "Expected no error when starting transaction")
	// Create user within the transaction
	createdUser, err := s.repo.CreateWithTx(tx, mockUser)
	s.NoError(err, "Expected no error when creating user with transaction")
	s.NotNil(createdUser, "Expected created user to be not nil")
	// Commit the transaction
	err = tx.Commit().Error
	s.NoError(err, "Expected no error when committing transaction")
	// Verify the user was created
	retrievedUser, err := s.repo.GetByID(createdUser.ID)
	s.NoError(err, "Expected no error when retrieving user by ID")
	s.NotNil(retrievedUser, "Expected retrieved user to be not nil")
	s.Equal("Transaction User", retrievedUser.Name, "Expected user name to be 'Transaction User'")
}

func (s *UserRepositoryTestSuite) TestCreateWithTx_Error_DuplicateEmail() {
	// Assume Email is unique
	user1 := &models.User{
		Email:    "duplicate@example.com",
		Name:     "user1",
		Password: "pass",
	}

	user2 := &models.User{
		Email:    "duplicate@example.com", // same email
		Name:     "user2",
		Password: "pass",
	}

	err := s.db.Create(user1).Error
	s.Require().NoError(err)

	tx := s.db.Begin()
	s.Require().NoError(tx.Error)

	createdUser, err := s.repo.CreateWithTx(tx, user2)
	s.Error(err, "Should return error due to duplicate email")
	s.Nil(createdUser, "Created user should be nil on duplicate constraint")

	tx.Rollback()
}

func (s *UserRepositoryTestSuite) TestGetDB() {
	db := s.repo.GetDB()
	s.NotNil(db, "Expected database connection to be not nil")
}

func TestUserRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(UserRepositoryTestSuite))
}
