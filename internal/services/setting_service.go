package services

import (
	"github.com/vfa-khuongdv/golang-cms/internal/models"
	"github.com/vfa-khuongdv/golang-cms/internal/repositories"
)

type SettingService struct {
	repo *repositories.SettingRepostitory
}

func NewSettingService(repo *repositories.SettingRepostitory) *SettingService {
	return &SettingService{repo: repo}
}

// GetSetting retrieves all settings from the repository
// Returns a pointer to a slice of Setting models and any error encountered
func (service *SettingService) GetSetting() (*[]models.Setting, error) {
	return service.repo.GetAll()
}
