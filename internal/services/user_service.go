package services

import (
	log "github.com/sirupsen/logrus"

	"github.com/vfa-khuongdv/golang-cms/internal/models"
	"github.com/vfa-khuongdv/golang-cms/internal/repositories"
	"github.com/vfa-khuongdv/golang-cms/internal/utils"
)

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
	return service.repo.Get(id)
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
	return service.repo.FindByEmail(email)
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

// PaginationUser fetches a paginated list of users and the total count.
// Parameters:
//   - page: The page number to retrieve (1-based indexing)
//   - limit: The number of users per page
//
// Returns:
//   - *[]models.User: A pointer to slice of user records for the requested page
//   - int64: The total number of user records across all pages
//
// Example:
//
//	users, total := service.PaginationUser(2, 10) // Gets page 2 with 10 users per page
func (service *UserService) PaginationUser(page int, limit int) (*[]models.User, int64, error) {
	offset, limit := utils.CalculatePagination(page, limit)
	return service.repo.PaginationUser(offset, limit)

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
	return service.repo.FindByToken(token)
}
