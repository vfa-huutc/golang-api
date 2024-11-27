package repositories

import (
	"github.com/vfa-khuongdv/golang-cms/internal/models"
	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

// NewUserRepsitory creates and returns a new UserRepository instance
// Parameters:
//   - db: Pointer to the gorm.DB database connection
//
// Returns:
//   - *UserRepository: Pointer to the newly created UserRepository instance
func NewUserRepsitory(db *gorm.DB) *UserRepository {
	return &UserRepository{
		db: db,
	}
}

// FindByUsername fetches a user by their username from the database
// Parameters:
//   - username: The username to search for
//
// Returns:
//   - *models.User: Pointer to the retrieved User model if found
//   - error: Error if user not found or if there was a database error
func (repo *UserRepository) FindByUsername(username string) (*models.User, error) {
	var user models.User
	if err := repo.db.Where("username = ?", username).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// Register creates a new user in the database
// Parameters:
//   - user: Pointer to the User model to be created
//
// Returns:
//   - error: Error if there was a problem creating the user, nil on success
func (repo *UserRepository) Register(user *models.User) error {
	if err := repo.db.Create(&user).Error; err != nil {
		return err
	}
	return nil
}

// GetUser retrieves a user from the database by their ID
// Parameters:
//   - id: The unique identifier of the user to retrieve
//
// Returns:
//   - *models.User: Pointer to the retrieved User model
//   - error: Error if the user is not found or if there was a database error
func (repo *UserRepository) GetUser(id uint) (*models.User, error) {
	var user models.User
	if err := repo.db.First(&user, id).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// PaginationUser retrieves a paginated list of users from the database
// Parameters:
//   - offset: Number of records to skip
//   - limit: Maximum number of records to return
//
// Returns:
//   - *[]models.User: Pointer to slice of User models containing the paginated results
//   - int64: Total count of all users in the database
//   - error: Error if there was a database error
func (repo *UserRepository) PaginationUser(offset, limit int) (*[]models.User, int64, error) {
	var users []models.User
	var total int64

	// Count total number of records
	if err := repo.db.Model(&models.User{}).Count(&total).Error; err != nil {
		return &[]models.User{}, 0, err
	}

	// Query with limit and offset
	if err := repo.db.Limit(limit).Offset(offset).Find(&users).Error; err != nil {
		return &[]models.User{}, total, err
	}

	return &users, total, nil

}
