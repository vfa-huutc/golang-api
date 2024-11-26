package repositories

import (
	"github.com/vfa-khuongdv/golang-cms/internal/models"
	"gorm.io/gorm"
)

type RefreshTokenRepository struct {
	db *gorm.DB
}

func NewRefreshTokenRepository(db *gorm.DB) *RefreshTokenRepository {
	return &RefreshTokenRepository{
		db: db,
	}
}

// Create refresh token
func (repo *RefreshTokenRepository) Create(token *models.RefreshToken) error {
	if err := repo.db.Save(&token).Error; err != nil {
		return err
	}
	return nil
}
