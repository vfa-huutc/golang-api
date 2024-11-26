package services

import (
	"errors"

	"github.com/vfa-khuongdv/golang-cms/configs"
	"github.com/vfa-khuongdv/golang-cms/internal/repositories"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	repo *repositories.UserRepository
}

func NewAuthService(repo *repositories.UserRepository) *AuthService {
	return &AuthService{
		repo: repo,
	}
}

func (service *AuthService) Login(username, password string) (string, error) {
	user, err := service.repo.FindByUsername(username)
	if err != nil {
		return "", errors.New("not found user")
	}

	// Validate password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return "", errors.New("invalid password")
	}

	return configs.GenerateToken(user.Username)
}
