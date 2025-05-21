package repositories

import (
	"github.com/vfa-khuongdv/golang-cms/internal/models"
	"gorm.io/gorm"
)

type ISettingRepository interface {
	GetAll() ([]models.Setting, error) // We return a slice of Setting models not a pointer to a slice because we don't need to modify the slice itself, just the elements inside it
	GetByKey(key string) (*models.Setting, error)
	Update(setting *models.Setting) error
	Create(setting *models.Setting) error
}

type SettingRepository struct {
	db *gorm.DB
}

func NewSettingRepository(db *gorm.DB) *SettingRepository {
	return &SettingRepository{db: db}
}

// GetAll retrieves all settings from the database
// Parameters:
//   - None
//
// Returns:
//   - []models.Setting: Pointer to slice of Setting models containing all settings
//   - error: Error if database operation fails, nil otherwise
func (repo *SettingRepository) GetAll() ([]models.Setting, error) {
	var settings []models.Setting

	if err := repo.db.Find(&settings).Error; err != nil {
		return nil, err
	}
	return settings, nil

}

// GetByKey retrieves a setting by its key from the database
// Parameters:
//   - key: string - The key to search for
//
// Returns:
//   - *models.Setting: Pointer to Setting model if found, nil if not found
//   - error: Error if database operation fails, nil otherwise
func (repo *SettingRepository) GetByKey(key string) (*models.Setting, error) {
	var setting models.Setting

	if err := repo.db.Model(&models.Setting{}).Where("setting_key = ?", key).First(&setting).Error; err != nil {
		return nil, err
	}
	return &setting, nil
}

// Update saves a setting to the database
// Parameters:
//   - setting: *models.Setting - Pointer to Setting model to be updated
//
// Returns:
//   - *models.Setting: Pointer to updated Setting model
//   - error: Error if database operation fails, nil otherwise
func (repo *SettingRepository) Update(setting *models.Setting) error {
	return repo.db.Save(setting).Error
}

// Create saves a new setting to the database
// Parameters:
//   - setting: *models.Setting - Pointer to Setting model to be created
//
// Returns:
//   - *models.Setting: Pointer to created Setting model
//   - error: Error if database operation fails, nil otherwise
func (repo *SettingRepository) Create(setting *models.Setting) error {
	return repo.db.Create(setting).Error
}
