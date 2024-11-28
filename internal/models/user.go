package models

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Email    string  `gorm:"type:varchar(45);unique;not null" json:"email" binding:"required,email"`
	Password string  `gorm:"type:varchar(255);not null" json:"password,omitempty" validate:"min=1,max=255"`
	Name     string  `gorm:"type:varchar(45);not null" json:"name" validate:"required,min=1,max=45"`
	Birthday *string `gorm:"type:date;default:null" json:"birthday" validate:"valid_birthday"`
	Address  *string `gorm:"type:varchar(255);default:null" json:"address" validate:"min=1,max=255"`
	Gender   int16   `gorm:"type:smallint;not null" validate:"required,oneof=0 1 2"`
	Token    *string `gorm:"type:varchar(100);default:null;unique" json:"token,omitempty"`
}
