package models

import (
	"time"

	"gorm.io/gorm"
)

type Role struct {
	ID          uint           `gorm:"column:id;primaryKey" json:"id"`
	Name        string         `gorm:"column:name;type:varchar(255);unique;not null" json:"name"`
	DisplayName string         `gorm:"column:display_name;type:varchar(255);not null" json:"displayName"`
	CreatedAt   time.Time      `gorm:"column:created_at" json:"createdAt"`
	UpdatedAt   time.Time      `gorm:"column:updated_at" json:"updatedAt"`
	DeletedAt   gorm.DeletedAt `gorm:"column:deleted_at;index" json:"deletedAt,omitempty"`

	// Relations
	Permissions []Permission `gorm:"many2many:role_permissions;" json:"permissions,omitempty"`
}
