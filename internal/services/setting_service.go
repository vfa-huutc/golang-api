package services

import (
	"github.com/vfa-khuongdv/golang-cms/internal/models"
	"github.com/vfa-khuongdv/golang-cms/internal/repositories"
	"github.com/vfa-khuongdv/golang-cms/pkg/errors"
)

type ISettingService interface {
	GetSetting() (*[]models.Setting, error)
	UpdateMany(settings *[]models.Setting) error
	GetSettingByKey(key string) (*models.Setting, error)
	Update(setting *models.Setting) error
	Create(setting *models.Setting) error
}

type SettingService struct {
	repo *repositories.SettingRepostitory
}

func NewSettingService(repo *repositories.SettingRepostitory) *SettingService {
	return &SettingService{repo: repo}
}

// GetSetting retrieves all settings from the repository
// Returns:
//   - *[]models.Setting: pointer to a slice of Setting models containing all settings
//   - error: any error encountered during the retrieval operation
func (service *SettingService) GetSetting() (*[]models.Setting, error) {
	data, err := service.repo.GetAll()
	if err != nil {
		return nil, errors.New(errors.ErrDatabaseQuery, err.Error())
	}
	return data, nil
}

// UpdateSetting updates multiple settings in the repository
// Parameters:
//   - settings: pointer to a slice of Setting models to be updated
//
// Returns:
//   - error: any error encountered during the update operation
func (service *SettingService) UpdateMany(settings *[]models.Setting) error {
	err := service.repo.UpdateMany(settings)
	if err != nil {
		return errors.New(errors.ErrDatabaseUpdate, err.Error())
	}
	return nil
}

// GetSettingByKey retrieves a specific setting from the repository by its key
// Parameters:
//   - key: string representing the unique identifier of the setting
//
// Returns:
//   - *models.Setting: pointer to the Setting model if found
//   - error: any error encountered during the retrieval operation
func (service *SettingService) GetSettingByKey(key string) (*models.Setting, error) {
	data, err := service.repo.GetByKey(key)
	if err != nil {
		return nil, errors.New(errors.ErrDatabaseQuery, err.Error())
	}
	return data, nil
}

// Update updates a single setting in the repository
// Parameters:
//   - setting: pointer to the Setting model to be updated
//
// Returns:
//   - *models.Setting: pointer to the updated Setting model
//   - error: any error encountered during the update operation
func (service *SettingService) Update(setting *models.Setting) error {
	err := service.repo.Update(setting)
	if err != nil {
		return errors.New(errors.ErrDatabaseUpdate, err.Error())
	}
	return nil
}

// Create creates a new setting in the repository
// Parameters:
//   - setting: pointer to the Setting model to be created
//
// Returns:
//   - *models.Setting: pointer to the created Setting model
//   - error: any error encountered during the creation operation
func (service *SettingService) Create(setting *models.Setting) error {
	err := service.repo.Create(setting)
	if err != nil {
		return errors.New(errors.ErrDatabaseInsert, err.Error())
	}
	return nil
}
