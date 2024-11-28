package models

import (
	"gorm.io/gorm"
)

type Setting struct {
	gorm.Model
	SettingKey string `gorm:"type:varchar(25);not null;column:setting_key"`
	Value      string `gorm:"type:varchar(255);not null"`
}
