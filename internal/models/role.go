package models

import (
	"gorm.io/gorm"
)

type Role struct {
	gorm.Model
	Name        string `gorm:"type:varchar(255);unique;not null"`
	DisplayName string `gorm:"type:varchar(255);not null"`
}
