package seeders

import (
	"fmt"

	"github.com/vfa-khuongdv/golang-cms/internal/models"
	"github.com/vfa-khuongdv/golang-cms/internal/utils"
	"gorm.io/gorm"
)

func SeedUsers(db *gorm.DB) error {
	users := []models.User{
		{
			Name:     "John Doe",
			Email:    "john@example.com",
			Password: utils.HashPassword("password123"),
		},
		{
			Name:     "Jane Smith",
			Email:    "jane@example.com",
			Password: utils.HashPassword("password123"),
		},
		{
			Name:     "Uncle Bob",
			Email:    "unclebob@example.com",
			Password: utils.HashPassword("password123"),
		},
	}

	for _, user := range users {
		if err := db.Create(&user).Error; err != nil {
			fmt.Printf("The user %v was run before\n", user.Name)
		}
	}

	return nil
}
