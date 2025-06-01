package services_test

import (
	"testing"

	originErrors "errors"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"github.com/vfa-khuongdv/golang-cms/internal/models"
	"github.com/vfa-khuongdv/golang-cms/internal/services"
	"github.com/vfa-khuongdv/golang-cms/tests/mocks"
)

type RefreshTokenServiceTestSuite struct {
	suite.Suite
	repo                *mocks.MockRefreshTokenRepository
	refreshTokenService *services.RefreshTokenService
}

func (s *RefreshTokenServiceTestSuite) SetupTest() {
	s.repo = new(mocks.MockRefreshTokenRepository)
	s.refreshTokenService = services.NewRefreshTokenService(s.repo)
}

func (s *RefreshTokenServiceTestSuite) TestCreate_Success() {
	user := &models.User{
		ID:    1,
		Email: "test@example.com",
	}

	ipAddress := "127.0.0.1"

	s.repo.On("Create", mock.MatchedBy(func(token *models.RefreshToken) bool {
		return token.UserID == user.ID && token.IpAddress == ipAddress
	})).Return(nil)

	result, err := s.refreshTokenService.Create(user, ipAddress)

	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), result)
	assert.Len(s.T(), result.Token, 60)
	assert.Greater(s.T(), result.ExpiresAt, int64(0))

	s.repo.AssertExpectations(s.T())
}

func (s *RefreshTokenServiceTestSuite) TestCreate_Error() {
	user := &models.User{
		ID:    1,
		Email: "test@example.com",
	}
	ipAddress := "127.0.0.1"
	s.repo.On("Create", mock.Anything).Return(originErrors.New("database error"))
	_, err := s.refreshTokenService.Create(user, ipAddress)
	assert.Error(s.T(), err)
	s.repo.AssertExpectations(s.T())
}

func TestRefreshTokenServiceTestSuite(t *testing.T) {
	suite.Run(t, new(RefreshTokenServiceTestSuite))
}

func (s *RefreshTokenServiceTestSuite) TestUpdate_Success() {
	originalToken := &models.RefreshToken{
		RefreshToken: "existing_token",
		IpAddress:    "",
		UsedCount:    0,
		ExpiredAt:    0,
		UserID:       1,
	}

	s.repo.On("FindByToken", "existing_token").Return(originalToken, nil).Once()
	s.repo.On("Update", mock.AnythingOfType("*models.RefreshToken")).Return(nil).Once()

	result, err := s.refreshTokenService.Update("existing_token", "127.0.0.2")

	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), result)
	assert.Equal(s.T(), originalToken.UserID, result.UserId)
	assert.Len(s.T(), result.Token.Token, 60)
	assert.Greater(s.T(), result.Token.ExpiresAt, int64(0))

	s.repo.AssertExpectations(s.T())
}

func (s *RefreshTokenServiceTestSuite) TestUpdate_TokenNotFound() {
	s.repo.On("FindByToken", "missing_token").Return((*models.RefreshToken)(nil), assert.AnError).Once()

	result, err := s.refreshTokenService.Update("missing_token", "127.0.0.1")

	assert.Error(s.T(), err)
	assert.Nil(s.T(), result)

	s.repo.AssertExpectations(s.T())
}

func (s *RefreshTokenServiceTestSuite) TestUpdate_Error() {
	originalToken := &models.RefreshToken{
		RefreshToken: "existing_token",
		IpAddress:    "",
		UsedCount:    0,
		ExpiredAt:    0,
		UserID:       1,
	}

	s.repo.On("FindByToken", "existing_token").Return(originalToken, nil).Once()
	s.repo.On("Update", mock.AnythingOfType("*models.RefreshToken")).Return(originErrors.New("Update item error")).Once()

	result, err := s.refreshTokenService.Update("existing_token", "127.0.0.1")

	assert.Error(s.T(), err)
	assert.Nil(s.T(), result)

	s.repo.AssertExpectations(s.T())
}
