package repositories

import (
	"github.com/vfa-khuongdv/golang-cms/internal/models"
	"gorm.io/gorm"
)

type RefreshTokenRepository struct {
	db *gorm.DB
}

// NewRefreshTokenRepository creates a new instance of RefreshTokenRepository
// Parameters:
//   - db: pointer to the gorm.DB instance for database operations
//
// Returns:
//   - *RefreshTokenRepository: pointer to the newly created RefreshTokenRepository
func NewRefreshTokenRepository(db *gorm.DB) *RefreshTokenRepository {
	return &RefreshTokenRepository{
		db: db,
	}
}

// Create creates a new refresh token in the database
// Parameters:
//   - token: pointer to the RefreshToken model to be saved
//
// Returns:
//   - error: nil if successful, error otherwise
func (repo *RefreshTokenRepository) Create(token *models.RefreshToken) error {
	if err := repo.db.Save(&token).Error; err != nil {
		return err
	}
	return nil
}

// FindByToken retrieves a refresh token from the database by its token value
// Parameters:
//   - token: string representing the refresh token to search for
//
// Returns:
//   - *models.RefreshToken: pointer to the found RefreshToken model, nil if not found
//   - error: nil if successful, error otherwise
func (repo *RefreshTokenRepository) FindByToken(token string) (*models.RefreshToken, error) {
	var refreshToken models.RefreshToken
	if err := repo.db.Where("refresh_token = ?", token).First(&refreshToken).Error; err != nil {
		return nil, err
	}
	return &refreshToken, nil
}

func (repo *RefreshTokenRepository) UpdateToken(token *models.RefreshToken) error {
	if err := repo.db.Save(token).Error; err != nil {
		return err
	}
	return nil
}
