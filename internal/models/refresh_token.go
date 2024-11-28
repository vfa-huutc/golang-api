package models

import (
	"gorm.io/gorm"
)

type RefreshToken struct {
	gorm.Model
	RefreshToken string `gorm:"type:varchar(60);not null" json:"refresh_token"`
	IpAddress    string `gorm:"type:varchar(45);not null" json:"ip_address"`
	UsedCount    int64  `gorm:"default:0" json:"used_count"`
	ExpiredAt    int64  `gorm:"not null" json:"expired_at"`
	UserID       uint   `gorm:"not null" json:"user_id"`
	User         User   `gorm:"constraint:OnDelete:CASCADE;foreignKey:UserID"`
}
