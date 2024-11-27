package models

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username string  `gorm:"unique;not null" json:"username"`
	Password string  `gorm:"not null" json:"password"`
	Name     string  `gorm:"not null" json:"name" validate:"required,min=2,max=255"`
	Birthday *string `gorm:"type:date;default:null" json:"birthday" validate:"required,valid_birthday"`
	Address  *string `gorm:"default:null" json:"address" validate:"required,min=1,max=255"`
	Gender   int16   `gorm:"type:smallint; not null" validate:"required,oneof=0 1 2"`
	Token    *string `gorm:"default:null" json:"token"`
}
