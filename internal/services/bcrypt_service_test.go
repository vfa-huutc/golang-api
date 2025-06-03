package services_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vfa-khuongdv/golang-cms/internal/services"
	"golang.org/x/crypto/bcrypt"
)

func TestBcryptService_HashAndCheckPassword(t *testing.T) {
	service := services.NewBcryptService()

	password := "securepassword123"
	hashedPassword, err := service.HashPassword(password)

	assert.NoError(t, err, "HashPassword should not return an error")
	assert.NotEmpty(t, hashedPassword, "Hashed password should not be empty")
	assert.True(t, service.CheckPasswordHash(password, hashedPassword), "CheckPasswordHash should return true for valid password")
	assert.False(t, service.CheckPasswordHash("wrongpassword", hashedPassword), "CheckPasswordHash should return false for invalid password")
}

func TestBcryptService_HashPasswordWithCost(t *testing.T) {
	service := services.NewBcryptService()

	password := "anotherpassword456"
	cost := bcrypt.MinCost

	hashedPassword, err := service.HashPasswordWithCost(password, cost)
	assert.NoError(t, err, "HashPasswordWithCost should not return an error")
	assert.NotEmpty(t, hashedPassword, "Hashed password should not be empty")
	assert.True(t, service.CheckPasswordHash(password, hashedPassword), "CheckPasswordHash should work with hashed password created using custom cost")
}

func TestBcryptService_HashPasswordWithInvalidCost(t *testing.T) {
	service := services.NewBcryptService()

	password := "invalidcost"
	invalidCost := 1000 // invalid bcrypt cost

	_, err := service.HashPasswordWithCost(password, invalidCost)
	assert.Error(t, err, "HashPasswordWithCost should return error for invalid cost")
}
