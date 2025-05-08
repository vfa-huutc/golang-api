package models

import (
	"gorm.io/gorm"
)

type UserRole struct {
	gorm.Model
	UserID uint `gorm:"not null;index"`
	RoleID uint `gorm:"not null;index"`
	User   User `gorm:"foreignKey:UserID"`
	Roles  Role `gorm:"foreignKey:RoleID"`
}
