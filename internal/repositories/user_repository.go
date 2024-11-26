package repositories

import (
	"log"

	"github.com/vfa-khuongdv/golang-cms/internal/models"
	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepsitory(db *gorm.DB) *UserRepository {
	return &UserRepository{
		db: db,
	}
}

// FindByUsername fetches a user by their username
func (repo *UserRepository) FindByUsername(username string) (*models.User, error) {
	var user models.User
	if err := repo.db.Where("username = ?", username).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// Create user
func (repo *UserRepository) Register(user *models.User) error {
	if err := repo.db.Create(&user).Error; err != nil {
		return err
	}
	return nil
}

// Get user
func (repo *UserRepository) GetUser(id uint) (*models.User, error) {
	var user models.User
	if err := repo.db.First(&user, id).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// Paginate user
func (repo *UserRepository) PaginationUser(offset, limit int) (*[]models.User, int64) {
	var users []models.User
	var total int64

	// Count total number of records
	if err := repo.db.Model(&models.User{}).Count(&total).Error; err != nil {
		log.Printf("Failed to count users: %s\n", err.Error())
		return &[]models.User{}, 0
	}

	// Query with limit and offset
	if err := repo.db.Limit(limit).Offset(offset).Find(&users).Error; err != nil {
		log.Printf("Query paginated list user error: %s\n", err.Error())
		return &[]models.User{}, total
	}

	return &users, total

}
