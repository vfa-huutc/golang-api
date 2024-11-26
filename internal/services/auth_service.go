package services

import (
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/vfa-khuongdv/golang-cms/configs"
	"github.com/vfa-khuongdv/golang-cms/internal/repositories"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	repo         *repositories.UserRepository
	tokenService *RefreshTokenService
}

type LoginResponse struct {
	AccessToken  configs.JwtResult
	RefreshToken configs.JwtResult
}

func NewAuthService(repo *repositories.UserRepository, tokenService *RefreshTokenService) *AuthService {
	return &AuthService{
		repo:         repo,
		tokenService: tokenService,
	}
}

// Login with username and password
func (service *AuthService) Login(username, password string, ctx *gin.Context) (*LoginResponse, error) {
	user, err := service.repo.FindByUsername(username)
	if err != nil {
		return nil, errors.New("not found user")
	}

	// Validate password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return nil, errors.New("invalid password")
	}

	// Generate refresh token
	token, err := configs.GenerateToken(user.Username)
	if err != nil {
		return nil, err
	}

	// Create new refresh token
	ipAddress := ctx.ClientIP()
	tokenResult, err := service.tokenService.Create(*user, ipAddress)
	if err != nil {
		return nil, err
	}

	res := &LoginResponse{
		AccessToken: configs.JwtResult{
			Token:     token.Token,
			ExpiresAt: token.ExpiresAt,
		},
		RefreshToken: configs.JwtResult{
			Token:     tokenResult.Token,
			ExpiresAt: tokenResult.ExpiresAt,
		},
	}

	return res, nil
}
