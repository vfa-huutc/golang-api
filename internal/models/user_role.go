package models

import (
	"time"

	"gorm.io/gorm"
)

type UserRole struct {
	ID        uint           `gorm:"column:id;primaryKey" json:"id"`
	UserID    uint           `gorm:"column:user_id;not null;index" json:"userId"`
	RoleID    uint           `gorm:"column:role_id;not null;index" json:"roleId"`
	CreatedAt time.Time      `gorm:"column:created_at" json:"createdAt"`
	UpdatedAt time.Time      `gorm:"column:updated_at" json:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at;index" json:"deletedAt,omitempty"`

	// Relations
	User  User `gorm:"foreignKey:UserID" json:"user"`
	Roles Role `gorm:"foreignKey:RoleID" json:"role"`
}
