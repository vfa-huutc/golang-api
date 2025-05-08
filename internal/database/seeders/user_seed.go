package seeders

import (
	"log"

	"github.com/vfa-khuongdv/golang-cms/internal/models"
	"github.com/vfa-khuongdv/golang-cms/internal/utils"
	"github.com/vfa-khuongdv/golang-cms/pkg/logger"
	"gorm.io/gorm"
)

type UserSeeder struct {
	User    *models.User
	RoleIds *[]uint
}

func SeedUsers(db *gorm.DB) error {
	users := []UserSeeder{
		{
			User: &models.User{
				Name:     "John Doe",
				Email:    "john@example.com",
				Password: utils.HashPassword("password123"),
			},
			RoleIds: &[]uint{1},
		},
		{
			User: &models.User{
				Name:     "Jane Smith",
				Email:    "jane@example.com",
				Password: utils.HashPassword("password123"),
			},
			RoleIds: &[]uint{2},
		},
	}

	for _, userData := range users {
		// Create new user
		if err := db.Create(&userData.User).Error; err != nil {
			logger.Errorf("Error creating user %s: %v", userData.User.Name, err)
			continue
		}
		// Create user-role associations
		for _, roleID := range *userData.RoleIds {
			userRole := models.UserRole{
				UserID: userData.User.ID,
				RoleID: roleID,
			}
			if err := db.Create(&userRole).Error; err != nil {
				logger.Errorf("Error adding role %d to user %s: %v", roleID, userData.User.Name, err)
			}
			log.Printf("Created user %s with %d roles", userData.User.Name, len(*userData.RoleIds))
		}
	}

	return nil
}
