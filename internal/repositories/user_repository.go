package repositories

import (
	"github.com/vfa-khuongdv/golang-cms/internal/models"
	"gorm.io/gorm"
)

type IUserRepository interface {
	GetAll() (*[]models.User, error)
	GetByID(id uint) (*models.User, error)
	Create(user *models.User) error
	Update(user *models.User) error
	Delete(user *models.User) error
	FindByField(field string, value interface{}) (*models.User, error)
	GetProfile(id uint) (*models.User, error)
	UpdateProfile(user *models.User) error
}

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

// GetAll retrieves all users from the database
// Parameters:
//   - None
//
// Returns:
//   - []models.User: Slice containing all User models in the database
//   - error: Error if there was a database error, nil on success
func (repo *UserRepository) GetAll() (*[]models.User, error) {
	var users []models.User
	if err := repo.db.Find(&users).Error; err != nil {
		return nil, err
	}
	return &users, nil
}

// GetByID retrieves a user from the database by their ID
// Parameters:
//   - id: The unique identifier of the user to retrieve
//
// Returns:
//   - *models.User: Pointer to the retrieved User model
//   - error: Error if the user is not found or if there was a database error
func (repo *UserRepository) GetByID(id uint) (*models.User, error) {
	var user models.User
	if err := repo.db.First(&user, id).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// Create creates a new user in the database
// Parameters:
//   - user: Pointer to the User model to be created
//
// Returns:
//   - error: Error if there was a problem creating the user, nil on success
func (repo *UserRepository) Create(user *models.User) error {
	if err := repo.db.Create(&user).Error; err != nil {
		return err
	}
	return nil
}

// Update updates an existing user in the database
// Parameters:
//   - user: Pointer to the User model to be updated
//
// Returns:
//   - error: Error if there was a problem updating the user, nil on success
func (repo *UserRepository) Update(user *models.User) error {
	if err := repo.db.Save(&user).Error; err != nil {
		return err
	}
	return nil
}

// Delete removes a user from the database
// Parameters:
//   - id: userId to be deleted
//
// Returns:
//   - error: Error if there was a problem deleting the user, nil on success
func (repo *UserRepository) Delete(userId uint) error {
	var user models.User
	return repo.db.Delete(&user, userId).Error
}

// FindByField retrieves a user from the database by a specified field and value
// Parameters:
//   - field: The database field name to search on
//   - value: The value to match against the specified field
//
// Returns:
//   - *models.User: Pointer to the retrieved User model if found
//   - error: Error if user not found or if there was a database error
func (repo *UserRepository) FindByField(field string, value interface{}) (*models.User, error) {
	var user models.User
	if err := repo.db.Where(field+" = ?", value).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// GetProfile retrieves a user's profile from the database by their ID
// Parameters:
//   - id: The unique identifier of the user whose profile is to be retrieved
//
// Returns:
//   - *models.User: Pointer to the retrieved User model containing profile information
//   - error: Error if the profile is not found or if there was a database error
func (repo *UserRepository) GetProfile(id uint) (*models.User, error) {
	var user models.User
	if err := repo.db.First(&user, id).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// UpdateProfile updates a user's profile information in the database
// Parameters:
//   - user: Pointer to the User model containing updated profile information
//
// Returns:
//   - error: Error if there was a problem updating the profile, nil on success
func (repo *UserRepository) UpdateProfile(user *models.User) error {
	return repo.db.Save(&user).Error
}
