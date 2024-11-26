package services

import (
	"log"

	"github.com/vfa-khuongdv/golang-cms/internal/models"
	"github.com/vfa-khuongdv/golang-cms/internal/repositories"
	"github.com/vfa-khuongdv/golang-cms/internal/utils"
)

type UserService struct {
	repo *repositories.UserRepository
}

// NewUserService creates a new instance of UserService with the provided UserRepository.
// Example: repo := &repositories.UserRepository{}; service := NewUserService(repo).
func NewUserService(repo *repositories.UserRepository) *UserService {
	return &UserService{
		repo: repo,
	}
}

// GetUser retrieves a user by their ID.
// Example: GetUser(1) retrieves the user with ID 1.
func (service *UserService) GetUser(id uint) (*models.User, error) {
	return service.repo.GetUser(id)
}

// CreateUser registers a new user in the system.
// Example: CreateUser(&models.User{Name: "John Doe", Email: "john@example.com"}).
func (service *UserService) CreateUser(user *models.User) error {
	if err := service.repo.Register(user); err != nil {
		log.Printf("Query user error %s\n", err)
		return err
	}
	return nil
}

// PaginationUser fetches a paginated list of users and the total count.
// Returns a slice of users and the total record count.
// Example: PaginationUser(2, 10) retrieves page 2 with 10 users per page.
func (service *UserService) PaginationUser(page int, limit int) (*[]models.User, int64) {
	offset, limit := utils.CalculatePagination(page, limit)
	return service.repo.PaginationUser(offset, limit)
}
