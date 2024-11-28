package models

import (
	"gorm.io/gorm"
)

type Setting struct {
	gorm.Model
	Key   string `gorm:"type:varchar(25); not null" json:"key"`
	Value string `gorm:"type:varchar(255); not null" json:"value"`
}
