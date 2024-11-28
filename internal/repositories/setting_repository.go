package repositories

import (
	"github.com/vfa-khuongdv/golang-cms/internal/models"
	"gorm.io/gorm"
)

type SettingRepostitory struct {
	db *gorm.DB
}

func NewSettingRepository(db *gorm.DB) *SettingRepostitory {
	return &SettingRepostitory{db: db}
}

// GetAll retrieves all settings from the database
// Parameters:
//   - None
//
// Returns:
//   - *[]models.Setting: Pointer to slice of Setting models containing all settings
//   - error: Error if database operation fails, nil otherwise
func (repo *SettingRepostitory) GetAll() (*[]models.Setting, error) {
	var settings []models.Setting

	if err := repo.db.Find(&settings).Error; err != nil {
		return nil, err
	}
	return &settings, nil

}

// UpdateMany updates multiple settings in the database
// Parameters:
//   - settings: Slice of Setting models containing the settings to update
//
// Returns:
//   - error: Error if database operation fails, nil otherwise
func (repo *SettingRepostitory) UpdateMany(settings []models.Setting) error {
	return repo.db.Transaction(func(tx *gorm.DB) error {
		for _, setting := range settings {
			if err := tx.Model(&models.Setting{}).Where("id = ?", setting.ID).Updates(setting).Error; err != nil {
				return err
			}
		}
		return nil
	})
}
