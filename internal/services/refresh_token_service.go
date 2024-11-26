package services

import (
	"time"

	"github.com/vfa-khuongdv/golang-cms/configs"
	"github.com/vfa-khuongdv/golang-cms/internal/models"
	"github.com/vfa-khuongdv/golang-cms/internal/repositories"
	"github.com/vfa-khuongdv/golang-cms/internal/utils"
)

type RefreshTokenService struct {
	repo *repositories.RefreshTokenRepository
}

func NewRefreshTokenService(repo *repositories.RefreshTokenRepository) *RefreshTokenService {
	return &RefreshTokenService{
		repo: repo,
	}
}

// Create new token
func (service *RefreshTokenService) Create(user models.User, ipAddress string) (*configs.JwtResult, error) {

	tokenString := utils.GenerateRandomString(60)
	expiredAt := time.Now().Add(time.Hour * 24 * 30).Unix()
	token := models.RefreshToken{
		RefreshToken: tokenString,
		IpAddress:    ipAddress, // ipaddress of user
		UsedCount:    0,         // init is zero
		ExpiredAt:    expiredAt, // 30 days
		UserID:       user.ID,   // userId
	}

	err := service.repo.Create(&token)
	if err != nil {
		return nil, err
	}

	return &configs.JwtResult{
		Token:     tokenString,
		ExpiresAt: expiredAt,
	}, nil
}
