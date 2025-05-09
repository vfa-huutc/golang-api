package seeders

import (
	"github.com/vfa-khuongdv/golang-cms/internal/models"
	"github.com/vfa-khuongdv/golang-cms/pkg/logger"
	"gorm.io/gorm"
)

type RoleSeeder struct {
	Role          *models.Role
	PermissionIds []uint
}

func SeedRoles(db *gorm.DB) error {
	roles := []RoleSeeder{
		{
			Role: &models.Role{
				ID:          1,
				Name:        "Admin",
				DisplayName: "Administrator",
			},
			PermissionIds: []uint{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12}, // Full access to all permissions
		},
		{
			Role: &models.Role{
				ID:          2,
				Name:        "User",
				DisplayName: "Regular User",
			},
			PermissionIds: []uint{1, 4, 6, 9, 11}, // Can view users, roles, settings but not modify them
		},
	}

	for _, roleData := range roles {
		// Create new role
		if err := db.Create(&roleData.Role).Error; err != nil {
			logger.Errorf("Error creating role %s: %v", roleData.Role.Name, err)
			continue
		}

		// Create role-permission associations
		for _, permID := range roleData.PermissionIds {
			rolePermission := models.RolePermission{
				RoleID:       roleData.Role.ID,
				PermissionID: permID,
			}
			if err := db.Create(&rolePermission).Error; err != nil {
				logger.Errorf("Error adding permission %d to role %s: %v", permID, roleData.Role.Name, err)
			}
		}

		logger.Infof("Created role %s with %d permissions", roleData.Role.Name, len(roleData.PermissionIds))
	}

	return nil
}
