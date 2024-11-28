package models

import "gorm.io/gorm"

type Permission struct {
	gorm.Model
	Resource string `gorm:"type:varchar(60);not null;uniqueIndex:idx_resource_action"`
	Action   string `gorm:"type:varchar(60);not null;uniqueIndex:idx_resource_action"`
}
