package models

import "gorm.io/gorm"

// RolePermission represents the many-to-many relationship between roles and permissions
type RolePermission struct {
	gorm.Model
	RoleID       uint `gorm:"uniqueIndex:idx_role_permission"`
	PermissionID uint `gorm:"uniqueIndex:idx_role_permission"`
}
