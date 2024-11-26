package models

import "gorm.io/gorm"

type RefreshToken struct {
	*gorm.Model
	RefreshToken string `gorm:"not null" json:"refres_token"`
	IpAddress    string `gorm:"not null" json:"ip_address"`
	UsedCount    int64  `gorm:"default:0" json:"used_count"`
	ExpiredAt    int64  `gorm:"no null" json:"expired_at"`
	UserID       uint   `gorm:"not null" json:"user_id"`
	User         User   `gorm:"constraint:OnDelete:CASCADE;"`
}
