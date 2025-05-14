package models

import (
	"time"

	"gorm.io/gorm"
)

type Setting struct {
	ID         uint           `gorm:"column:id;primaryKey" json:"id"`
	SettingKey string         `gorm:"column:setting_key;type:varchar(25);not null" json:"settingKey"`
	Value      string         `gorm:"column:value;type:varchar(255);not null" json:"value"`
	CreatedAt  time.Time      `gorm:"column:created_at" json:"createdAt"`
	UpdatedAt  time.Time      `gorm:"column:updated_at" json:"updatedAt"`
	DeletedAt  gorm.DeletedAt `gorm:"column:deleted_at;index" json:"deletedAt,omitempty"`
}
