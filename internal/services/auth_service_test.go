package services_test

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"github.com/vfa-khuongdv/golang-cms/internal/models"
	"github.com/vfa-khuongdv/golang-cms/internal/services"
	"github.com/vfa-khuongdv/golang-cms/internal/utils"
	"github.com/vfa-khuongdv/golang-cms/pkg/errors"
	"github.com/vfa-khuongdv/golang-cms/tests/mocks"
	"gorm.io/gorm"
)

// TestSuite is a struct that holds the mock repositories and the service under test
type AuthServiceTestSuite struct {
	suite.Suite
	repo                *mocks.MockUserRepository
	refreshTokenService *mocks.MockRefreshTokenService
	service             services.IAuthService
	bcryptService       *mocks.MockBcryptService
	jwtService          *mocks.MockJWTService
}

func (s *AuthServiceTestSuite) SetupTest() {
	s.repo = new(mocks.MockUserRepository)
	s.refreshTokenService = new(mocks.MockRefreshTokenService)
	s.bcryptService = new(mocks.MockBcryptService)
	s.jwtService = new(mocks.MockJWTService)

	s.service = services.NewAuthService(
		s.repo,
		s.refreshTokenService,
		s.bcryptService,
		s.jwtService,
	)
}

func (s *AuthServiceTestSuite) TestLoginSuccess() {
	// Set up the expected user and mock repository behavior
	email := "test@example.com"
	password := "password123"

	user := &models.User{
		ID:       1,
		Email:    email,
		Password: "hashed_password",
	}
	ip := "127.0.0.1"

	// Mock FindByField to return user
	s.repo.On("FindByField", "email", email).Return(user, nil).Once()
	// Mock BcryptService to validate password and return true
	s.bcryptService.On("CheckPasswordHash", password, user.Password).Return(true).Once()
	// Mock JWTService to generate access token
	s.jwtService.On("GenerateToken", user.ID).Return(&services.JwtResult{
		Token:     "mocked-access-token",
		ExpiresAt: time.Now().Add(1 * time.Hour).Unix(),
	}, nil).Once()

	// Mock Create to return a valid JWT result
	s.refreshTokenService.On("Create", user, ip).Return(&services.JwtResult{
		Token:     "mocked-refresh-token",
		ExpiresAt: time.Now().Add(24 * time.Hour).Unix(),
	}, nil)

	ginCtx, _ := gin.CreateTestContext(nil)
	ginCtx.Request = &http.Request{RemoteAddr: ip + ":12345"}

	// Call the Login method
	resp, err := s.service.Login(email, password, ginCtx)
	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), resp)
	assert.NotEmpty(s.T(), resp.AccessToken.Token)
	assert.Equal(s.T(), "mocked-refresh-token", resp.RefreshToken.Token)

}

func (s *AuthServiceTestSuite) TestLogin_UserNotFound() {
	email := "nonexistent@example.com"
	password := "password123"

	s.repo.On("FindByField", "email", email).Return((*models.User)(nil), gorm.ErrRecordNotFound)

	w := httptest.NewRecorder()
	ginCtx, _ := gin.CreateTestContext(w)

	resp, err := s.service.Login(email, password, ginCtx)

	assert.Error(s.T(), err)
	assert.Nil(s.T(), resp)

	appErr, ok := err.(*errors.AppError)
	s.Require().True(ok, "error should be of type *errors.AppError")
	assert.Equal(s.T(), errors.ErrNotFound, appErr.Code) // Code is 1001 for not found

	s.repo.AssertExpectations(s.T())

}

func (s *AuthServiceTestSuite) TestLogin_InvalidPassword() {
	email := "test@example.com"
	wrongPassword := "wrongpass"
	user := &models.User{
		ID:       1,
		Email:    email,
		Password: "hashed_password", // Assume this is a invalid hashed password
	}

	s.repo.On("FindByField", "email", email).Return(user, nil)
	s.bcryptService.On("CheckPasswordHash", wrongPassword, user.Password).Return(false).Once()

	ginCtx, _ := gin.CreateTestContext(nil)

	resp, err := s.service.Login(email, wrongPassword, ginCtx)
	assert.Error(s.T(), err)
	assert.Nil(s.T(), resp)

	appErr, ok := err.(*errors.AppError)
	s.Require().True(ok, "error should be of type *errors.AppError")
	assert.Equal(s.T(), errors.ErrInvalidPassword, appErr.Code) // Code is 3003 for invalid password

	s.repo.AssertExpectations(s.T())

}

func (s *AuthServiceTestSuite) TestLogin_CreateTokenError() {
	email := "test@example.com"
	password := "password123"
	hashedPassword := utils.HashPassword(password)
	user := &models.User{
		ID:       1,
		Email:    email,
		Password: hashedPassword,
	}
	ipAddress := "127.0.0.1"

	// Mock user repository and bcrypt service
	s.repo.On("FindByField", "email", email).Return(user, nil)
	s.bcryptService.On("CheckPasswordHash", password, user.Password).Return(true).Once()
	s.jwtService.On("GenerateToken", user.ID).
		Return(&services.JwtResult{
			Token:     "mocked-access-token",
			ExpiresAt: time.Now().Add(1 * time.Hour).Unix(),
		}, nil).Once()
	s.refreshTokenService.On("Create", user, ipAddress).
		Return(nil, errors.New(errors.ErrInvalidData, "token generation failed")).
		Once()

	// Create a proper gin.Context with ResponseWriter
	w := httptest.NewRecorder()
	ginCtx, _ := gin.CreateTestContext(w)
	ginCtx.Request = &http.Request{
		RemoteAddr: ipAddress + ":12345",
	}

	resp, err := s.service.Login(email, password, ginCtx)

	assert.Error(s.T(), err)
	assert.Nil(s.T(), resp)

	appErr, ok := err.(*errors.AppError)
	s.Require().True(ok, "error should be of type *errors.AppError")

	assert.Equal(s.T(), errors.ErrInvalidData, appErr.Code)

	s.repo.AssertExpectations(s.T())
	s.refreshTokenService.AssertExpectations(s.T())
}

func (s *AuthServiceTestSuite) TestLogin_JwtError() {
	email := "test@example.com"
	password := "password123"
	hashedPassword := utils.HashPassword(password)
	user := &models.User{
		ID:       1,
		Email:    email,
		Password: hashedPassword,
	}
	ipAddress := "127.0.0.1"

	// Mock user repository and bcrypt service
	s.repo.On("FindByField", "email", email).Return(user, nil)
	s.bcryptService.On("CheckPasswordHash", password, user.Password).Return(true).Once()
	s.jwtService.On("GenerateToken", user.ID).
		Return(&services.JwtResult{}, errors.New(errors.ErrInternal, "jwt generation failed")).Once()

	w := httptest.NewRecorder()
	ginCtx, _ := gin.CreateTestContext(w)
	ginCtx.Request = &http.Request{
		RemoteAddr: ipAddress + ":12345",
	}

	resp, err := s.service.Login(email, password, ginCtx)

	assert.Error(s.T(), err)
	assert.Nil(s.T(), resp)

	appErr, ok := err.(*errors.AppError)
	s.Require().True(ok, "error should be of type *errors.AppError")

	assert.Equal(s.T(), errors.ErrInternal, appErr.Code)

	s.repo.AssertExpectations(s.T())
	s.refreshTokenService.AssertExpectations(s.T())
}

func (s *AuthServiceTestSuite) TestRefreshToken_Success() {
	// Test input values
	oldRefreshToken := "valid-refresh-token"
	ipAddress := "127.0.0.1"
	userID := uint(1)

	// Mock new refresh token that would be returned by refresh token service
	mockRefreshToken := &services.JwtResult{
		Token:     "new-refresh-token",
		ExpiresAt: time.Now().Add(24 * time.Hour * 30).Unix(), // 30 days
	}
	mockRes := &services.RefreshTokenResult{
		UserId: userID,
		Token:  mockRefreshToken,
	}

	// Mock user that would be returned by user repository
	mockUser := &models.User{
		ID:    userID,
		Email: "user@example.com",
	}

	// Should update refresh token with correct old token and IP
	s.refreshTokenService.On("Update", oldRefreshToken, ipAddress).Return(mockRes, nil).Once()
	s.repo.On("GetByID", mockRes.UserId).Return(mockUser, nil).Once()
	s.jwtService.On("GenerateToken", mockUser.ID).Return(&services.JwtResult{
		Token:     "new-access-token",
		ExpiresAt: time.Now().Add(1 * time.Hour).Unix(),
	}, nil).Once()

	// Setup gin test context with IP
	w := httptest.NewRecorder()
	ginCtx, _ := gin.CreateTestContext(w)
	ginCtx.Request = &http.Request{RemoteAddr: ipAddress + ":12345"}

	// Execute the refresh token flow
	result, err := s.service.RefreshToken(oldRefreshToken, ginCtx)

	// Assert no errors occurred
	s.NoError(err, "Expected no error from RefreshToken")
	s.NotNil(result, "Expected result not to be nil")

	// Verify response structure and values
	s.NotEmpty(result.AccessToken.Token, "Expected access token to be set")
	s.True(result.AccessToken.ExpiresAt > time.Now().Unix(), "Expected access token to expire in the future")

	// Verify refresh token matches mock
	s.Equal(mockRefreshToken.Token, result.RefreshToken.Token, "Refresh token should match mock")
	s.Equal(mockRefreshToken.ExpiresAt, result.RefreshToken.ExpiresAt, "Refresh token expiry should match mock")

	// Validate mock expectations
	s.refreshTokenService.AssertExpectations(s.T())
	s.repo.AssertExpectations(s.T())
}

func (s *AuthServiceTestSuite) TestRefreshToken_UpdateError() {
	// Test input values
	invalidToken := "invalid-refresh-token"
	ipAddress := "127.0.0.1"

	// Mock refresh token service to return error for invalid token
	mockError := errors.New(errors.ErrDBQuery, "token not found")
	s.refreshTokenService.On("Update", invalidToken, ipAddress).Return(nil, mockError).Once()

	// Setup gin test context with IP
	w := httptest.NewRecorder()
	ginCtx, _ := gin.CreateTestContext(w)
	ginCtx.Request = &http.Request{RemoteAddr: ipAddress + ":12345"}

	// Execute the refresh token flow
	result, err := s.service.RefreshToken(invalidToken, ginCtx) // Assert error was returned
	s.Error(err, "Expected error for invalid refresh token")
	s.Nil(result, "Expected nil result for error case")
	s.Contains(err.Error(), mockError.Error(), "Expected database query error message")

	// Validate mock expectations
	s.refreshTokenService.AssertExpectations(s.T())
	// User repo should not be called when token is invalid
	s.repo.AssertNotCalled(s.T(), "GetByID")
}

func (s *AuthServiceTestSuite) TestRefreshToken_GetByIDError() {
	oldRefreshToken := "old-refresh-token"

	ipAddress := "127.0.0.1"
	// Mock new refresh token that would be returned by refresh token service
	mockRefreshToken := &services.JwtResult{
		Token:     "new-refresh-token",
		ExpiresAt: time.Now().Add(24 * time.Hour * 30).Unix(), // 30 days
	}
	mockRes := &services.RefreshTokenResult{
		UserId: 1,
		Token:  mockRefreshToken,
	}

	// Should update refresh token with correct old token and IP
	s.refreshTokenService.On("Update", oldRefreshToken, ipAddress).Return(mockRes, nil).Once()
	// Should fetch user with ID from refresh token
	s.repo.On("GetByID", mockRes.UserId).Return((*models.User)(nil), gorm.ErrInvalidData).Once()

	// Setup gin test context with IP
	w := httptest.NewRecorder()
	ginCtx, _ := gin.CreateTestContext(w)
	ginCtx.Request = &http.Request{RemoteAddr: ipAddress + ":12345"}

	// Execute the refresh token flow
	result, err := s.service.RefreshToken(oldRefreshToken, ginCtx)
	s.Error(err, "Expected error for valid refresh token")
	s.Nil(result, "Expected nil result for error case")

	appErr, ok := err.(*errors.AppError)
	s.Require().True(ok, "error should be of type *errors.AppError")
	assert.Equal(s.T(), errors.ErrDBQuery, appErr.Code) // Code is 2001 for database query error

	s.T().Logf("Error message: %s", err.Error())

	s.refreshTokenService.AssertExpectations(s.T())
}

func (s *AuthServiceTestSuite) TestRefreshToken_JwtError() {
	oldRefreshToken := "old-refresh-token"

	user := &models.User{
		ID:    1,
		Email: "email@example.com",
	}

	ipAddress := "127.0.0.1"
	// Mock new refresh token that would be returned by refresh token service
	mockRefreshToken := &services.JwtResult{
		Token:     "new-refresh-token",
		ExpiresAt: time.Now().Add(24 * time.Hour * 30).Unix(), // 30 days
	}
	mockRes := &services.RefreshTokenResult{
		UserId: 1,
		Token:  mockRefreshToken,
	}

	// Should update refresh token with correct old token and IP
	s.refreshTokenService.On("Update", oldRefreshToken, ipAddress).Return(mockRes, nil).Once()
	// Should fetch user with ID from refresh token
	s.repo.On("GetByID", mockRes.UserId).Return(user, nil).Once()
	// Should generate new access token for user
	s.jwtService.On("GenerateToken", user.ID).Return(&services.JwtResult{}, errors.New(errors.ErrInternal, "jwt generation failed")).Once()

	// Setup gin test context with IP
	w := httptest.NewRecorder()
	ginCtx, _ := gin.CreateTestContext(w)
	ginCtx.Request = &http.Request{RemoteAddr: ipAddress + ":12345"}

	// Execute the refresh token flow
	result, err := s.service.RefreshToken(oldRefreshToken, ginCtx)
	s.Error(err, "Expected error for valid refresh token")
	s.Nil(result, "Expected nil result for error case")

	appErr, ok := err.(*errors.AppError)
	s.Require().True(ok, "error should be of type *errors.AppError")
	assert.Equal(s.T(), errors.ErrInternal, appErr.Code) // Code is 2001 for database query error

	s.T().Logf("Error message: %s", err.Error())

	s.refreshTokenService.AssertExpectations(s.T())
}

func TestAuthServiceTestSuite(t *testing.T) {
	gin.SetMode(gin.TestMode)
	suite.Run(t, new(AuthServiceTestSuite))
}
