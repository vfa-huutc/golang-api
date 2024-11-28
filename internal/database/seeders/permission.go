package seeders

import (
	"fmt"

	"github.com/vfa-khuongdv/golang-cms/internal/models"
	"gorm.io/gorm"
)

func SeedPermissions(db *gorm.DB) error {
	permissions := []models.Permission{
		// User resource permissions
		{
			Resource: "users", // Resource: User management
			Action:   "index", // Action: List all users
		},
		{
			Resource: "users",  // Resource: User management
			Action:   "create", // Action: Create new user
		},
		{
			Resource: "users",  // Resource: User management
			Action:   "update", // Action: Update existing user
		},
		{
			Resource: "users", // Resource: User management
			Action:   "view",  // Action: View user details
		},
		{
			Resource: "users",  // Resource: User management
			Action:   "delete", // Action: Delete user
		},
		// Role resource permissions
		{
			Resource: "roles", // Resource: Role management
			Action:   "index", // Action: List all roles
		},
		{
			Resource: "roles",  // Resource: Role management
			Action:   "create", // Action: Create new role
		},
		{
			Resource: "roles",  // Resource: Role management
			Action:   "update", // Action: Update existing role
		},
		{
			Resource: "roles", // Resource: Role management
			Action:   "view",  // Action: View role details
		},
		{
			Resource: "roles",  // Resource: Role management
			Action:   "delete", // Action: Delete role
		},
		// Settings resource permissions
		{
			Resource: "settings", // Resource: System settings
			Action:   "view",     // Action: View settings
		},
		{
			Resource: "settings", // Resource: System settings
			Action:   "update",   // Action: Update settings
		},
	}

	for _, permission := range permissions {
		if err := db.Create(&permission).Error; err != nil {
			fmt.Printf("The permission %v, action %v was run before\n", permission.Resource, permission.Action)
		}
	}

	return nil
}
