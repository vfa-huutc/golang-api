package services

import (
	"log"

	"github.com/vfa-khuongdv/golang-cms/configs"
	"github.com/vfa-khuongdv/golang-cms/internal/models"
)

func GetUser(id uint) *models.User {
	var user models.User

	if err := configs.DB.First(&user, id).Error; err != nil {
		log.Printf("Query user error %s\n", err.Error())
	}

	return &user
}

func PaginationUser(page int, limit int) (*[]models.User, int64) {
	var users []models.User
	var total int64

	// Count total number of records
	if err := configs.DB.Model(&models.User{}).Count(&total).Error; err != nil {
		log.Printf("Failed to count users: %s\n", err.Error())
		return &[]models.User{}, 0
	}

	// Calculate offset
	offset := (page - 1) * limit

	// Query with limit and offset
	if err := configs.DB.Limit(limit).Offset(offset).Find(&users).Error; err != nil {
		log.Printf("Query paginated list user error: %s\n", err.Error())
		return &[]models.User{}, total
	}

	return &users, total
}
