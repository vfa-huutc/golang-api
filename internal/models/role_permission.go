package models

import (
	"time"

	"gorm.io/gorm"
)

type RolePermission struct {
	ID           uint           `gorm:"column:id;primaryKey" json:"id"`
	RoleID       uint           `gorm:"column:role_id;uniqueIndex:idx_role_permission" json:"roleId"`
	PermissionID uint           `gorm:"column:permission_id;uniqueIndex:idx_role_permission" json:"permissionId"`
	CreatedAt    time.Time      `gorm:"column:created_at" json:"createdAt"`
	UpdatedAt    time.Time      `gorm:"column:updated_at" json:"updatedAt"`
	DeletedAt    gorm.DeletedAt `gorm:"column:deleted_at;index" json:"deletedAt,omitempty"`
}
