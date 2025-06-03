package services

import (
	"github.com/gin-gonic/gin"
	"github.com/vfa-khuongdv/golang-cms/internal/repositories"
	"github.com/vfa-khuongdv/golang-cms/pkg/errors"
)

type IAuthService interface {
	Login(email, password string, ctx *gin.Context) (*LoginResponse, error)
	RefreshToken(token string, ctx *gin.Context) (*LoginResponse, error)
}

type AuthService struct {
	repo                repositories.IUserRepository
	refreshTokenService IRefreshTokenService
	bcryptService       IBcryptService
	jwtService          IJWTService
}

type LoginResponse struct {
	AccessToken  JwtResult `json:"accessToken"`
	RefreshToken JwtResult `json:"refreshToken"`
}

// NewAuthService creates and returns a new instance of AuthService
// Parameters:
//   - repo: User repository for database operations
//   - tokenService: Service for handling refresh token operations
//
// Returns:
//   - *AuthService: New AuthService instance initialized with the provided dependencies
func NewAuthService(repo repositories.IUserRepository, refreshTokenService IRefreshTokenService, bcryptService IBcryptService, jwtService IJWTService) *AuthService {
	return &AuthService{
		repo:                repo,
		refreshTokenService: refreshTokenService,
		bcryptService:       bcryptService,
		jwtService:          jwtService,
	}
}

// Login authenticates a user with their username and password
// Parameters:
//   - username: The username of the user trying to log in
//   - password: The password provided by the user
//   - ctx: Gin context containing request information
//
// Returns:
//   - *LoginResponse: Contains access token and refresh token if login successful
//   - error: Returns error if login fails (user not found, invalid password, token generation fails)
func (service *AuthService) Login(email, password string, ctx *gin.Context) (*LoginResponse, error) {
	user, err := service.repo.FindByField("email", email)
	if err != nil {
		return nil, errors.New(errors.ErrNotFound, err.Error())
	}

	// Validate password
	if isValid := service.bcryptService.CheckPasswordHash(password, user.Password); !isValid {
		return nil, errors.New(errors.ErrInvalidPassword, "Invalid credentials")
	}

	// Generate access token
	accessToken, err := service.jwtService.GenerateToken(user.ID)
	if err != nil {
		return nil, errors.New(errors.ErrInternal, err.Error())
	}

	// Create new refresh token
	ipAddress := ctx.ClientIP()
	refreshToken, err := service.refreshTokenService.Create(user, ipAddress)
	if err != nil {
		return nil, err
	}

	res := &LoginResponse{
		AccessToken: JwtResult{
			Token:     accessToken.Token,
			ExpiresAt: accessToken.ExpiresAt,
		},
		RefreshToken: JwtResult{
			Token:     refreshToken.Token,
			ExpiresAt: refreshToken.ExpiresAt,
		},
	}

	return res, nil
}

// RefreshToken generates new access and refresh tokens using an existing refresh token
// Parameters:
//   - token: The existing refresh token string
//   - ctx: Gin context containing request information
//
// Returns:
//   - *LoginResponse: Contains new access token and refresh token if successful
//   - error: Returns error if token refresh fails (invalid token, user not found, token generation fails)
func (service *AuthService) RefreshToken(token string, ctx *gin.Context) (*LoginResponse, error) {
	// Get the client's IP address from the request context
	ipAddress := ctx.ClientIP()
	// Create new refresh token using the refresh token service
	res, err := service.refreshTokenService.Update(token, ipAddress)
	if err != nil {
		return nil, err // error is already wrapped by the service, so we can return it directly
	}

	// Get user details from the database using the user ID from refresh token
	user, err := service.repo.GetByID(res.UserId)
	if err != nil {
		return nil, errors.New(errors.ErrDBQuery, err.Error())
	}
	// Generate new access token for the user
	resultToken, err := service.jwtService.GenerateToken(user.ID)
	if err != nil {
		return nil, errors.New(errors.ErrInternal, err.Error())
	}

	// Return new access and refresh tokens
	return &LoginResponse{
		AccessToken: JwtResult{
			Token:     resultToken.Token,
			ExpiresAt: resultToken.ExpiresAt,
		},
		RefreshToken: *res.Token,
	}, nil

}
