package services

import (
	"github.com/vfa-khuongdv/golang-cms/internal/models"
	"github.com/vfa-khuongdv/golang-cms/internal/repositories"
	"github.com/vfa-khuongdv/golang-cms/pkg/errors"
)

type IUserService interface {
	GetUser(id uint) (*models.User, error)
	GetUserByEmail(email string) (*models.User, error)
	CreateUser(user *models.User, roleIds []uint) error
	UpdateUser(user *models.User) error
	DeleteUser(id uint) error
	GetUserByToken(token string) (*models.User, error)
	GetProfile(id uint) (*models.User, error)
	UpdateProfile(user *models.User) error
}

type UserService struct {
	repo repositories.IUserRepository
}

func NewUserService(repo repositories.IUserRepository) *UserService {
	return &UserService{
		repo: repo,
	}
}

// GetUser retrieves a user by their ID from the database.
// Parameters:
//   - id: The unique identifier of the user to retrieve
//
// Returns:
//   - *models.User: A pointer to the user record if found
//   - error: nil if successful, otherwise returns the error that occurred
//
// Example:
//
//	user, err := service.GetUser(1) // Gets user with ID 1
func (service *UserService) GetUser(id uint) (*models.User, error) {
	data, err := service.repo.GetByID(id)
	if err != nil {
		return nil, errors.New(errors.ErrDatabaseQuery, err.Error())
	}
	return data, nil
}

// GetUserByEmail retrieves a user by their email address from the database.
// Parameters:
//   - email: The email address of the user to retrieve
//
// Returns:
//   - *models.User: A pointer to the user record if found
//   - error: nil if successful, otherwise returns the error that occurred
//
// Example:
//
//	user, err := service.GetUserByEmail("john@example.com")
func (service *UserService) GetUserByEmail(email string) (*models.User, error) {
	data, err := service.repo.FindByField("email", email)
	if err != nil {
		return nil, errors.New(errors.ErrDatabaseQuery, err.Error())
	}
	return data, nil
}

// CreateUser creates a new user in the database using the provided user data
// Parameters:
//   - user: Pointer to models.User containing the user information to create
//   - roleIds: Slice of role IDs to assign to the user
//
// Returns:
//   - error: nil if successful, otherwise returns the error that occurred
//
// Example:
//
//	user := &models.User{
//	    Name: "John Doe",
//	    Email: "john@example.com",
//	}
//	err := service.CreateUser(user)
func (service *UserService) CreateUser(user *models.User, roleIds []uint) error {
	// Start a transaction to ensure both user creation and role assignments succeed or fail together
	tx := service.repo.GetDB().Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Create the user within the transaction
	createdUser, err := service.repo.Create(user)
	if err != nil {
		tx.Rollback()
		return errors.New(errors.ErrDatabaseInsert, err.Error())
	}

	// Create user-role associations for each role ID
	for _, roleId := range roleIds {
		userRole := models.UserRole{
			UserID: createdUser.ID,
			RoleID: roleId,
		}
		if err := tx.Create(&userRole).Error; err != nil {
			tx.Rollback()
			return errors.New(errors.ErrDatabaseInsert, "Failed to assign roles: "+err.Error())
		}
	}

	// Commit the transaction
	if err := tx.Commit().Error; err != nil {
		return errors.New(errors.ErrDatabaseInsert, "Transaction failed: "+err.Error())
	}

	return nil
}

// UpdateUser updates an existing user's information in the database.
// Parameters:
//   - user: Pointer to models.User containing the updated user information
//
// Returns:
//   - error: nil if successful, otherwise returns the error that occurred
//
// Example:
//
//	user := &models.User{
//	    ID: 1,
//	    Name: "Updated Name",
//	    Email: "updated@example.com",
//	}
//	err := service.UpdateUser(user)
func (service *UserService) UpdateUser(user *models.User) error {
	err := service.repo.Update(user)
	if err != nil {
		return errors.New(errors.ErrDatabaseUpdate, err.Error())
	}
	return nil
}

// DeleteUser removes a user from the database by their ID.
// Parameters:
//   - id: The unique identifier of the user to delete
//
// Returns:
//   - error: nil if successful, otherwise returns the error that occurred
//
// Example:
//
//	err := service.DeleteUser(1) // Deletes user with ID 1
func (service *UserService) DeleteUser(id uint) error {
	err := service.repo.Delete(id)
	if err != nil {
		return errors.New(errors.ErrDatabaseDelete, err.Error())
	}
	return nil
}

// GetUserByToken retrieves a user by their authentication token from the database.
// Parameters:
//   - token: The authentication token string associated with the user
//
// Returns:
//   - *models.User: A pointer to the user record if found
//   - error: nil if successful, otherwise returns the error that occurred
//
// Example:
//
//	user, err := service.GetUserByToken("abc123token")
func (service *UserService) GetUserByToken(token string) (*models.User, error) {
	data, err := service.repo.FindByField("token", token)
	if err != nil {
		return nil, errors.New(errors.ErrDatabaseQuery, err.Error())
	}
	return data, nil
}

// GetProfile retrieves a user's profile information by their ID from the database.
// Parameters:
//   - id: The unique identifier of the user whose profile to retrieve
//
// Returns:
//   - *models.User: A pointer to the user profile record if found
//   - error: nil if successful, otherwise returns the error that occurred
//
// Example:
//
//	profile, err := service.GetProfile(1) // Gets profile for user with ID 1
func (service *UserService) GetProfile(id uint) (*models.User, error) {
	data, err := service.repo.GetProfile(id)
	if err != nil {
		return nil, errors.New(errors.ErrDatabaseQuery, err.Error())
	}
	return data, nil
}

// UpdateProfile updates a user's profile information in the database.
// Parameters:
//   - user: Pointer to models.User containing the updated profile information
//
// Returns:
//   - error: nil if successful, otherwise returns the error that occurred
//
// Example:
//
//	user := &models.User{
//	    ID: 1,
//	    Name: "Updated Name",
//	    Bio: "Updated bio"
//	}
//	err := service.UpdateProfile(user)
func (service *UserService) UpdateProfile(user *models.User) error {
	err := service.repo.UpdateProfile(user)
	if err != nil {
		return errors.New(errors.ErrDatabaseUpdate, err.Error())
	}
	return nil
}
