package models

import (
	"time"

	"gorm.io/gorm"
)

type Permission struct {
	ID        uint           `gorm:"column:id;primaryKey" json:"id"`
	Resource  string         `gorm:"column:resource;type:varchar(60);not null;uniqueIndex:idx_resource_action" json:"resource"`
	Action    string         `gorm:"column:action;type:varchar(60);not null;uniqueIndex:idx_resource_action" json:"action"`
	CreatedAt time.Time      `gorm:"column:created_at" json:"created_at"`
	UpdatedAt time.Time      `gorm:"column:updated_at" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at;index" json:"deleted_at,omitempty"`
}
