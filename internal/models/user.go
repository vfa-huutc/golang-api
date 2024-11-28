package models

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Email     string  `gorm:"type:varchar(45);unique;not null"`
	Password  string  `gorm:"type:varchar(255);not null" json:"-"`
	Name      string  `gorm:"type:varchar(45);not null"`
	Birthday  *string `gorm:"type:date;default:null"`
	Address   *string `gorm:"type:varchar(255);default:null"`
	Gender    int16   `gorm:"type:smallint;not null"`
	Token     *string `gorm:"type:varchar(100);default:null;unique" json:"-"`
	ExpiredAt *int64  `gorm:"type:bigint;default:null"`
}
