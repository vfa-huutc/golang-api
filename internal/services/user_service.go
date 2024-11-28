package services

import (
	log "github.com/sirupsen/logrus"

	"github.com/vfa-khuongdv/golang-cms/internal/models"
	"github.com/vfa-khuongdv/golang-cms/internal/repositories"
)

type IUserService interface {
	GetUser(id uint) (*models.User, error)
	GetUserByEmail(email string) (*models.User, error)
	CreateUser(user *models.User) error
	UpdateUser(user *models.User) error
	DeleteUser(id uint) error
	GetUserByToken(token string) (*models.User, error)
	GetProfile(id uint) (*models.User, error)
	UpdateProfile(user *models.User) error
}

type UserService struct {
	repo *repositories.UserRepository
}

// NewUserService creates a new instance of UserService with the provided UserRepository.
// Parameters:
//   - repo: A pointer to the UserRepository that will handle data operations
//
// Returns:
//   - *UserService: A pointer to the newly created UserService instance
//
// Example:
//
//	repo := &repositories.UserRepository{}
//	service := NewUserService(repo)
func NewUserService(repo *repositories.UserRepository) *UserService {
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
	return service.repo.GetByID(id)
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
	return service.repo.FindByField("email", email)
}

// CreateUser creates a new user in the database using the provided user data
// Parameters:
//   - user: Pointer to models.User containing the user information to create
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
func (service *UserService) CreateUser(user *models.User) error {
	if err := service.repo.Create(user); err != nil {
		log.Errorf("Query user error %s\n", err)
		return err
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
	if err := service.repo.Update(user); err != nil {
		log.Errorf("Update user error %s\n", err)
		return err
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
	return service.repo.Delete(id)
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
	return service.repo.FindByField("token", token)
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
	return service.repo.GetProfile(id)
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
	return service.repo.UpdateProfile(user)
}
