package models

import (
	"gorm.io/gorm"
)

type Role struct {
	gorm.Model
	ID          uint         `gorm:"primaryKey"`
	Name        string       `gorm:"type:varchar(255);unique;not null"`
	DisplayName string       `gorm:"type:varchar(255);not null"`
	Permissions []Permission `gorm:"many2many:role_permissions;"`
}
