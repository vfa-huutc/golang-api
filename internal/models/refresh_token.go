package models

import (
	"gorm.io/gorm"
)

type RefreshToken struct {
	gorm.Model
	RefreshToken string `gorm:"type:varchar(60);not null"`
	IpAddress    string `gorm:"type:varchar(45);not null"`
	UsedCount    int64  `gorm:"default:0"`
	ExpiredAt    int64  `gorm:"not null"`
	UserID       uint   `gorm:"not null"`
	User         User   `gorm:"constraint:OnDelete:CASCADE;foreignKey:UserID"`
}
