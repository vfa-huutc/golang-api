package services_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/vfa-khuongdv/golang-cms/internal/models"
	"github.com/vfa-khuongdv/golang-cms/internal/services"
	"github.com/vfa-khuongdv/golang-cms/internal/utils"
	"github.com/vfa-khuongdv/golang-cms/tests/mocks"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type UserServiceTestSuite struct {
	suite.Suite
	db      *gorm.DB
	repo    *mocks.MockUserRepository
	service *services.UserService
}

func (s *UserServiceTestSuite) SetupTest() {
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	s.Require().NoError(err)
	s.Require().NotNil(db)

	err = db.AutoMigrate(&models.User{}, &models.Role{}, &models.UserRole{})
	s.Require().NoError(err)
	s.db = db
	s.repo = new(mocks.MockUserRepository)
	s.service = services.NewUserService(s.repo)

}

func (s *UserServiceTestSuite) TearDownTest() {
	s.repo.AssertExpectations(s.T())
}

func (s *UserServiceTestSuite) TestPaginateUser() {
	s.Run("Success", func() {
		// Mock repo
		expectedPagination := &utils.Pagination{
			Page:       1,
			Limit:      10,
			TotalItems: 2,
			TotalPages: 1,
			Data: []models.User{
				{ID: 1, Email: ""},
				{ID: 2, Email: ""},
			},
		}
		s.repo.On("PaginateUser", 1, 10).Return(expectedPagination, nil).Once()

		// Call service
		pagination, err := s.service.PaginateUser(1, 10)
		s.NoError(err)
		s.Equal(expectedPagination, pagination)
	})
	s.Run("Error", func() {
		// Mock repo
		s.repo.On("PaginateUser", 1, 10).Return(&utils.Pagination{}, errors.New("db error")).Once()

		// Call service
		pagination, err := s.service.PaginateUser(1, 10)
		s.Error(err)
		s.Nil(pagination)
	})
}

func (s *UserServiceTestSuite) TestGetUser() {
	s.Run("Success", func() {
		// Mock repo
		expectedUser := &models.User{ID: 1, Email: "example@gmail.com", Password: "password123"}
		s.repo.On("GetByID", uint(1)).Return(expectedUser, nil).Once()
		// Call service
		user, err := s.service.GetUser(1)
		s.NoError(err)
		s.Equal(expectedUser, user)
	})
	s.Run("Error", func() {
		// Mock repo
		s.repo.On("GetByID", uint(999)).Return(&models.User{}, errors.New("user not found")).Once()

		// Call service
		user, err := s.service.GetUser(999)
		s.Error(err)
		s.Nil(user)
	})
}

func (s *UserServiceTestSuite) TestGetUserByEmail() {
	s.Run("Success", func() {
		// Mock repo
		expectedUser := &models.User{ID: 1, Email: "example@gmail.com", Password: "password123"}
		s.repo.On("FindByField", "email", "example@gmail.com").Return(expectedUser, nil).Once()

		// Call service
		user, err := s.service.GetUserByEmail("example@gmail.com")
		s.NoError(err)
		s.Equal(expectedUser, user)
	})
	s.Run("Error", func() {
		// Mock repo
		s.repo.On("FindByField", "email", "notfound@gmail.com").Return(&models.User{}, errors.New("user not found")).Once()

		// Call service
		user, err := s.service.GetUserByEmail("notfound@gmail.com")
		s.Error(err)
		s.Nil(user)
	})
}

func (s *UserServiceTestSuite) TestDeleteUser() {
	s.Run("Success", func() {
		// Mock repo
		s.repo.On("Delete", uint(1)).Return(nil).Once()

		// Call service
		err := s.service.DeleteUser(1)
		s.NoError(err)
	})

	s.Run("Error", func() {
		// Mock repo
		s.repo.On("Delete", uint(999)).Return(errors.New("user not found")).Once()

		// Call service
		err := s.service.DeleteUser(999)
		s.Error(err)
	})
}

func (s *UserServiceTestSuite) TestGetUserByToken() {
	s.Run("Success", func() {
		// Mock repo
		expectedUser := &models.User{ID: 1, Email: "email@example.com", Password: "password123"}
		s.repo.On("FindByField", "token", "valid_token").Return(expectedUser, nil).Once()
		// Call service
		user, err := s.service.GetUserByToken("valid_token")
		s.NoError(err)
		s.Equal(expectedUser, user)
	})
	s.Run("Error", func() {
		// Mock repo
		s.repo.On("FindByField", "token", "invalid_token").Return(&models.User{}, errors.New("user not found")).Once()

		// Call service
		user, err := s.service.GetUserByToken("invalid_token")
		s.Error(err)
		s.Nil(user)
	})
}

func (s *UserServiceTestSuite) TestGetProfile() {
	s.Run("Success", func() {
		// Mock repo
		expectedUser := &models.User{ID: 1, Email: "email@example.com", Password: "password123"}
		s.repo.On("GetProfile", uint(1)).Return(expectedUser, nil).Once()
		// Call service
		user, err := s.service.GetProfile(1)
		s.NoError(err)
		s.Equal(expectedUser, user)
	})
	s.Run("Error", func() {
		// Mock repo
		s.repo.On("GetProfile", uint(999)).Return(&models.User{}, errors.New("profile not found")).Once()

		// Call service
		user, err := s.service.GetProfile(999)
		s.Error(err)
		s.Nil(user)
	})
}

func (s *UserServiceTestSuite) TestUpdateProfile() {
	s.Run("Success", func() {
		// Mock repo
		user := &models.User{ID: 1, Email: "", Password: "newpassword123"}
		s.repo.On("UpdateProfile", user).Return(nil).Once()
		// Call service
		err := s.service.UpdateProfile(user)
		s.NoError(err)
	})
	s.Run("Error", func() {
		// Mock repo
		user := &models.User{ID: 999, Email: "", Password: "newpassword123"}
		s.repo.On("UpdateProfile", user).Return(errors.New("update failed")).Once()

		// Call service
		err := s.service.UpdateProfile(user)
		s.Error(err)
	})
}

func (s *UserServiceTestSuite) TestUpdateUser() {
	s.Run("Success", func() {
		// Mock repo
		user := &models.User{ID: 1, Email: "updated@example.com", Password: "newpassword123"}
		s.repo.On("Update", user).Return(nil).Once()

		// Call service
		err := s.service.UpdateUser(user)
		s.NoError(err)
	})

	s.Run("Error", func() {
		// Mock repo
		user := &models.User{ID: 999, Email: "updated@example.com", Password: "newpassword123"}
		s.repo.On("Update", user).Return(errors.New("update failed")).Once()

		// Call service
		err := s.service.UpdateUser(user)
		s.Error(err)
	})
}

func TestUserServiceTestSuite(t *testing.T) {
	suite.Run(t, new(UserServiceTestSuite))
}
