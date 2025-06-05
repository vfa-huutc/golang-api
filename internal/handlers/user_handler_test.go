package handlers_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/vfa-khuongdv/golang-cms/internal/constants"
	"github.com/vfa-khuongdv/golang-cms/internal/handlers"
	"github.com/vfa-khuongdv/golang-cms/internal/models"
	"github.com/vfa-khuongdv/golang-cms/internal/utils"
	"github.com/vfa-khuongdv/golang-cms/pkg/apperror"
	"github.com/vfa-khuongdv/golang-cms/tests/mocks"
)

func TestPaginateUser(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("PaginateUser Success", func(t *testing.T) {
		userService := new(mocks.MockUserService)
		redisService := new(mocks.MockRedisService)
		bcryptService := new(mocks.MockBcryptService)
		handler := handlers.NewUserHandler(userService, redisService, bcryptService)

		expected := &utils.Pagination{
			Page:       1,
			Limit:      10,
			TotalItems: 2,
			TotalPages: 1,
			Data: []models.User{
				{ID: 1, Email: "email1@example.com", Name: "User One", Gender: 1, CreatedAt: time.Date(2023, 10, 1, 0, 0, 0, 0, time.UTC), UpdatedAt: time.Date(2023, 10, 1, 0, 0, 0, 0, time.UTC)},
				{ID: 2, Email: "email2@example.com", Name: "User Two", Gender: 2, CreatedAt: time.Date(2023, 10, 2, 0, 0, 0, 0, time.UTC), UpdatedAt: time.Date(2023, 10, 2, 0, 0, 0, 0, time.UTC)},
			},
		}
		// Mock the service method
		userService.On("PaginateUser", 1, 10).Return(expected, nil)

		// Create a test context
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest("GET", "/api/v1/users?page=1&limit=10", nil)

		// Call the handler
		handler.PaginationUser(c)

		// Assert the response
		assert.Equal(t, http.StatusOK, w.Code)
		assert.JSONEq(t, `{"page":1,"limit":10,"totalItems":2,"totalPages":1,"data":[{"id":1,"email":"email1@example.com","name":"User One","gender":1,"createdAt":"2023-10-01T00:00:00Z","updatedAt":"2023-10-01T00:00:00Z","deletedAt":null},{"id":2,"email":"email2@example.com","name":"User Two","gender":2,"createdAt":"2023-10-02T00:00:00Z","updatedAt":"2023-10-02T00:00:00Z","deletedAt":null}]}`, w.Body.String())

		// Assert mocks
		userService.AssertExpectations(t)
		redisService.AssertExpectations(t)
		bcryptService.AssertExpectations(t)
	})

	t.Run("PaginateUser - Error", func(t *testing.T) {
		userService := new(mocks.MockUserService)
		redisService := new(mocks.MockRedisService)
		bcryptService := new(mocks.MockBcryptService)
		handler := handlers.NewUserHandler(userService, redisService, bcryptService)

		// Mock the service method
		userService.On("PaginateUser", 1, 10).Return(&utils.Pagination{}, apperror.NewDBQueryError("Database query error"))

		// Create a test context
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest("GET", "/api/v1/users?page=1&limit=10", nil)

		// Call the handler
		handler.PaginationUser(c)

		// Assert the response

		var expectedBody = map[string]any{
			"code":    float64(apperror.ErrDBQuery),
			"message": "Database query error",
		}
		var actualBody map[string]any
		json.Unmarshal(w.Body.Bytes(), &actualBody)
		assert.Equal(t, expectedBody["code"], actualBody["code"])
		assert.Equal(t, expectedBody["message"], actualBody["message"])
		assert.Equal(t, http.StatusInternalServerError, w.Code)

		// Assert mocks
		userService.AssertExpectations(t)
		redisService.AssertExpectations(t)
		bcryptService.AssertExpectations(t)
	})
}

func TestCreateUser(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Initialize the validator
	utils.InitValidator()

	t.Run("CreateUser - Success", func(t *testing.T) {
		userService := new(mocks.MockUserService)
		redisService := new(mocks.MockRedisService)
		bcryptService := new(mocks.MockBcryptService)
		handler := handlers.NewUserHandler(userService, redisService, bcryptService)

		// Mock the CreateUser method
		userService.On("CreateUser", mock.AnythingOfType("*models.User"), mock.AnythingOfType("[]uint")).Return(nil)
		bcryptService.On("HashPassword", "password").Return("$2a$10$examplehash", nil)

		requestBody := map[string]any{
			"email":    "email@example.com",
			"password": "password",
			"name":     "User",
			"birthday": "2000-01-01",
			"address":  "123 Street",
			"gender":   1,
			"role_ids": []uint{1, 2},
		}
		body, _ := json.Marshal(requestBody)

		// Create a test context
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest("POST", "/api/v1/users", bytes.NewBuffer(body))

		// Call the CreateUser handler
		handler.CreateUser(c)

		// Assert the response
		assert.Equal(t, http.StatusCreated, w.Code)
		assert.JSONEq(t, `{"message":"Create user successfully"}`, w.Body.String())

		// Assert mocks
		userService.AssertExpectations(t)
		redisService.AssertExpectations(t)
		bcryptService.AssertExpectations(t)
	})

	t.Run("CreateUser - Validation Error", func(t *testing.T) {
		userService := new(mocks.MockUserService)
		redisService := new(mocks.MockRedisService)
		bcryptService := new(mocks.MockBcryptService)

		tests := []struct {
			name           string
			reqBody        string
			expectedCode   float64
			expectedMsg    string
			expectedFields []apperror.FieldError
		}{
			{
				name:         "MissingEmail",
				reqBody:      `{}`,
				expectedCode: float64(4001),
				expectedMsg:  "Validation failed",
				expectedFields: []apperror.FieldError{
					{Field: "email", Message: "email is required"},
					{Field: "password", Message: "password is required"},
					{Field: "name", Message: "name is required"},
					{Field: "birthday", Message: "birthday is required"},
					{Field: "address", Message: "address is required"},
					{Field: "gender", Message: "gender is required"},
					{Field: "role_ids", Message: "role_ids is required"},
				},
			},
			{
				name:         "EmptyEmail",
				reqBody:      `{"email":""}`,
				expectedCode: float64(4001),
				expectedMsg:  "Validation failed",
				expectedFields: []apperror.FieldError{
					{Field: "email", Message: "email is required"},
					{Field: "password", Message: "password is required"},
					{Field: "name", Message: "name is required"},
					{Field: "birthday", Message: "birthday is required"},
					{Field: "address", Message: "address is required"},
					{Field: "gender", Message: "gender is required"},
					{Field: "role_ids", Message: "role_ids is required"},
				},
			},
			{
				name:         "InvalidEmailFormat",
				reqBody:      `{"email":"not-an-email"}`,
				expectedCode: float64(4001),
				expectedMsg:  "Validation failed",
				expectedFields: []apperror.FieldError{
					{Field: "email", Message: "email must be a valid email address"},
					{Field: "password", Message: "password is required"},
					{Field: "name", Message: "name is required"},
					{Field: "birthday", Message: "birthday is required"},
					{Field: "address", Message: "address is required"},
					{Field: "gender", Message: "gender is required"},
					{Field: "role_ids", Message: "role_ids is required"},
				},
			},
			{
				name:         "MissingPassword",
				reqBody:      `{"email":"email@example.com"}`,
				expectedCode: float64(4001),
				expectedMsg:  "Validation failed",
				expectedFields: []apperror.FieldError{
					{Field: "password", Message: "password is required"},
					{Field: "name", Message: "name is required"},
					{Field: "birthday", Message: "birthday is required"},
					{Field: "address", Message: "address is required"},
					{Field: "gender", Message: "gender is required"},
					{Field: "role_ids", Message: "role_ids is required"},
				},
			},
			{
				name:         "EmptyPassword",
				reqBody:      `{"email":"email@example.com","password":""}`,
				expectedCode: float64(4001),
				expectedMsg:  "Validation failed",
				expectedFields: []apperror.FieldError{
					{Field: "password", Message: "password is required"},
					{Field: "name", Message: "name is required"},
					{Field: "birthday", Message: "birthday is required"},
					{Field: "address", Message: "address is required"},
					{Field: "gender", Message: "gender is required"},
					{Field: "role_ids", Message: "role_ids is required"},
				},
			},
			{
				name:         "ShortPassword",
				reqBody:      `{"email":"email@example.com","password":"123"}`,
				expectedCode: float64(4001),
				expectedMsg:  "Validation failed",
				expectedFields: []apperror.FieldError{
					{Field: "password", Message: "password must be at least 6 characters long or numeric"},
					{Field: "name", Message: "name is required"},
					{Field: "birthday", Message: "birthday is required"},
					{Field: "address", Message: "address is required"},
					{Field: "gender", Message: "gender is required"},
					{Field: "role_ids", Message: "role_ids is required"},
				},
			},
			{
				name:         "LongPassword",
				reqBody:      `{"email":"email@example.com","password":"` + strings.Repeat("a", 256) + `"}`,
				expectedCode: float64(4001),
				expectedMsg:  "Validation failed",
				expectedFields: []apperror.FieldError{
					{Field: "password", Message: "password must be at most 255 characters long or numeric"},
					{Field: "name", Message: "name is required"},
					{Field: "birthday", Message: "birthday is required"},
					{Field: "address", Message: "address is required"},
					{Field: "gender", Message: "gender is required"},
					{Field: "role_ids", Message: "role_ids is required"},
				},
			},
			{
				name:         "MissingName",
				reqBody:      `{"email":"email@example.com","password":"password"}`,
				expectedCode: float64(4001),
				expectedMsg:  "Validation failed",
				expectedFields: []apperror.FieldError{
					{Field: "name", Message: "name is required"},
					{Field: "birthday", Message: "birthday is required"},
					{Field: "address", Message: "address is required"},
					{Field: "gender", Message: "gender is required"},
					{Field: "role_ids", Message: "role_ids is required"},
				},
			},
			{
				name:         "NameNotBlank",
				reqBody:      `{"email":"email@example.com","password":"password","name":"  "}`,
				expectedCode: float64(4001),
				expectedMsg:  "Validation failed",
				expectedFields: []apperror.FieldError{
					{Field: "name", Message: "name must not be blank"},
					{Field: "birthday", Message: "birthday is required"},
					{Field: "address", Message: "address is required"},
					{Field: "gender", Message: "gender is required"},
					{Field: "role_ids", Message: "role_ids is required"},
				},
			},
			{
				name:         "EmptyName",
				reqBody:      `{"email":"email@example.com","password":"password","name":""}`,
				expectedCode: float64(4001),
				expectedMsg:  "Validation failed",
				expectedFields: []apperror.FieldError{
					{Field: "name", Message: "name is required"},
					{Field: "birthday", Message: "birthday is required"},
					{Field: "address", Message: "address is required"},
					{Field: "gender", Message: "gender is required"},
					{Field: "role_ids", Message: "role_ids is required"},
				},
			},
			{
				name:         "LongName",
				reqBody:      `{"email":"email@example.com","password":"password","name": "` + strings.Repeat("a", 46) + `","birthday":"2000-01-01","address":"address","gender":1}`,
				expectedCode: float64(4001),
				expectedMsg:  "Validation failed",
				expectedFields: []apperror.FieldError{
					{Field: "name", Message: "name must be at most 45 characters long or numeric"},
					{Field: "role_ids", Message: "role_ids is required"},
				},
			},
			{
				name:         "MissingBirthday",
				reqBody:      `{"email":"email@example.com","password":"password","name": "Bob"}`,
				expectedCode: float64(4001),
				expectedMsg:  "Validation failed",
				expectedFields: []apperror.FieldError{
					{Field: "birthday", Message: "birthday is required"},
					{Field: "address", Message: "address is required"},
					{Field: "gender", Message: "gender is required"},
					{Field: "role_ids", Message: "role_ids is required"},
				},
			},
			{
				name:         "InvalidBirthdayFormat",
				reqBody:      `{"email":"email@example.com","password":"password","name": "Bob","birthday":"invalid-date"}`,
				expectedCode: float64(4001),
				expectedMsg:  "Validation failed",
				expectedFields: []apperror.FieldError{
					{Field: "birthday", Message: "birthday must be a valid date (YYYY-MM-DD) and not in the future"},
					{Field: "address", Message: "address is required"},
					{Field: "gender", Message: "gender is required"},
					{Field: "role_ids", Message: "role_ids is required"},
				},
			},
			{
				name:         "FutureBirthday",
				reqBody:      `{"email":"email@example.com","password":"password","name": "Bob","birthday":"3000-01-01"}`,
				expectedCode: float64(4001),
				expectedMsg:  "Validation failed",
				expectedFields: []apperror.FieldError{
					{Field: "birthday", Message: "birthday must be a valid date (YYYY-MM-DD) and not in the future"},
					{Field: "address", Message: "address is required"},
					{Field: "gender", Message: "gender is required"},
					{Field: "role_ids", Message: "role_ids is required"},
				},
			},
			{
				name:         "MissingAddress",
				reqBody:      `{"email":"email@example.com","password":"password","name": "Bob","birthday":"2000-01-01"}`,
				expectedCode: float64(4001),
				expectedMsg:  "Validation failed",
				expectedFields: []apperror.FieldError{
					{Field: "address", Message: "address is required"},
					{Field: "gender", Message: "gender is required"},
					{Field: "role_ids", Message: "role_ids is required"},
				},
			},
			{
				name:         "AddressNotBlank",
				reqBody:      `{"email":"email@example.com","password":"password","name": "Bob","birthday":"2000-01-01","address":"  "}`,
				expectedCode: float64(4001),
				expectedMsg:  "Validation failed",
				expectedFields: []apperror.FieldError{
					{Field: "address", Message: "address must not be blank"},
					{Field: "gender", Message: "gender is required"},
					{Field: "role_ids", Message: "role_ids is required"},
				},
			},
			{
				name:         "EmptyAddress",
				reqBody:      `{"email":"email@example.com","password":"password","name": "Bob","birthday":"2000-01-01","address":""}`,
				expectedCode: float64(4001),
				expectedMsg:  "Validation failed",
				expectedFields: []apperror.FieldError{
					{Field: "address", Message: "address must be at least 1 characters long or numeric"},
					{Field: "gender", Message: "gender is required"},
					{Field: "role_ids", Message: "role_ids is required"},
				},
			},
			{
				name:         "LongAddress",
				reqBody:      `{"email":"email@example.com","password":"password","name": "Bob","birthday":"2000-01-01","address":"` + strings.Repeat("a", 256) + `"}`,
				expectedCode: float64(4001),
				expectedMsg:  "Validation failed",
				expectedFields: []apperror.FieldError{
					{Field: "address", Message: "address must be at most 255 characters long or numeric"},
					{Field: "gender", Message: "gender is required"},
					{Field: "role_ids", Message: "role_ids is required"},
				},
			},
			{
				name:         "MissingGender",
				reqBody:      `{"email":"email@example.com","password":"password","name": "Bob","birthday":"2000-01-01","address":"address"}`,
				expectedCode: float64(4001),
				expectedMsg:  "Validation failed",
				expectedFields: []apperror.FieldError{
					{Field: "gender", Message: "gender is required"},
					{Field: "role_ids", Message: "role_ids is required"},
				},
			},
			{
				name:         "InvalidGender",
				reqBody:      `{"email":"email@example.com","password":"password","name": "Bob","birthday":"2000-01-01","address":"address", "gender": 4}`,
				expectedCode: float64(4001),
				expectedMsg:  "Validation failed",
				expectedFields: []apperror.FieldError{
					{Field: "gender", Message: "gender must be one of [1 2 3]"},
					{Field: "role_ids", Message: "role_ids is required"},
				},
			},
			{
				name:           "StringGender",
				reqBody:        `{"email":"email@example.com","password":"password","name": "Bob","birthday":"2000-01-01","address":"address", "gender": "not_numeric"}`,
				expectedCode:   float64(4001),
				expectedMsg:    "json: cannot unmarshal string into Go struct field .gender of type int16",
				expectedFields: nil, // specific error case
			},
			{
				name:         "MissingRoleIDs",
				reqBody:      `{"email":"email@example.com","password":"password","name": "Bob","birthday":"2000-01-01","address":"address","gender":1}`,
				expectedCode: float64(4001),
				expectedMsg:  "Validation failed",
				expectedFields: []apperror.FieldError{
					{Field: "role_ids", Message: "role_ids is required"},
				},
			},
			{
				name:           "EmptyRoleIDs",
				reqBody:        `{"email":"email@example.com","password":"password","name": "Bob","birthday":"2000-01-01","address":"address","gender":1,"role_ids":""}`,
				expectedCode:   float64(4001),
				expectedMsg:    "json: cannot unmarshal string into Go struct field .role_ids of type []uint",
				expectedFields: nil, // specific error case
			},
			{
				name:           "InvalidRoleIDsIdNotNumeric",
				reqBody:        `{"email":"email@example.com","password":"password","name": "Bob","birthday":"2000-01-01","address":"address","gender":1,"role_ids":["not_numeric"]}`,
				expectedCode:   float64(4001),
				expectedMsg:    "json: cannot unmarshal string into Go struct field .role_ids of type uint",
				expectedFields: nil, // specific error case
			},
			{
				name:         "InvalidRoleIDIsEmptyArray",
				reqBody:      `{"email":"email@example.com","password":"password","name": "Bob","birthday":"2000-01-01","address":"address","gender":1,"role_ids":[]}`,
				expectedCode: float64(4001),
				expectedMsg:  "Validation failed",
				expectedFields: []apperror.FieldError{
					// TODO: fixed in the future assume that role_ids is array of uint, so it must be at least 1 characters long or numeric
					{Field: "role_ids", Message: "role_ids must be at least 1 characters long or numeric"},
				},
			},
		}

		for _, tc := range tests {
			t.Run(tc.name, func(t *testing.T) {
				handler := handlers.NewUserHandler(userService, redisService, bcryptService)

				// Create a test context
				w := httptest.NewRecorder()
				c, _ := gin.CreateTestContext(w)
				c.Request, _ = http.NewRequest("POST", "/api/v1/user", bytes.NewBufferString(tc.reqBody))

				// Call the CreateUser handler
				handler.CreateUser(c)

				// Assert the response
				expectedBody := map[string]any{
					"code":    tc.expectedCode,
					"message": tc.expectedMsg,
					"fields":  tc.expectedFields,
				}
				var actualBody map[string]any
				json.Unmarshal(w.Body.Bytes(), &actualBody)
				assert.Equal(t, http.StatusBadRequest, w.Code)
				assert.Equal(t, expectedBody["code"], actualBody["code"])
				assert.Equal(t, expectedBody["message"], actualBody["message"])
				assert.Equal(t, tc.expectedFields, utils.MapJsonToFieldErrors(actualBody["fields"]))

				// Assert mocks
				userService.AssertExpectations(t)
				redisService.AssertExpectations(t)
				bcryptService.AssertExpectations(t)
			})
		}
	})

	t.Run("Create user Error", func(t *testing.T) {
		userService := new(mocks.MockUserService)
		redisService := new(mocks.MockRedisService)
		bcryptService := new(mocks.MockBcryptService)
		handler := handlers.NewUserHandler(userService, redisService, bcryptService)

		// Mock the service methods
		bcryptService.On("HashPassword", "password").Return("$2a$10$examplehash", nil)
		userService.On("CreateUser", mock.AnythingOfType("*models.User"), mock.AnythingOfType("[]uint")).
			Return(apperror.NewDBInsertError("Database insert error"))

		requestBody := map[string]any{
			"email":    "email@example.com",
			"password": "password",
			"name":     "Bob",
			"birthday": "2000-01-01",
			"address":  "123 Street",
			"gender":   1,
			"role_ids": []uint{1, 2},
		}

		body, _ := json.Marshal(requestBody)

		// Create a test context
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest("POST", "/api/v1/user", bytes.NewBuffer(body))

		// Call the handler
		handler.CreateUser(c)

		// Assert the response
		var expectedBody = map[string]any{
			"code":    float64(apperror.ErrDBInsert),
			"message": "Database insert error",
		}
		var actualBody map[string]any
		json.Unmarshal(w.Body.Bytes(), &actualBody)
		assert.Equal(t, expectedBody["code"], actualBody["code"])
		assert.Equal(t, expectedBody["message"], actualBody["message"])
		assert.Equal(t, http.StatusInternalServerError, w.Code)

		// Assert mocks
		userService.AssertExpectations(t)
		redisService.AssertExpectations(t)
		bcryptService.AssertExpectations(t)

	})

	t.Run("Error Bcrypt Hash Password", func(t *testing.T) {
		userService := new(mocks.MockUserService)
		redisService := new(mocks.MockRedisService)
		bcryptService := new(mocks.MockBcryptService)
		handler := handlers.NewUserHandler(userService, redisService, bcryptService)

		// Mock the service methods
		bcryptService.On("HashPassword", "password").Return("", errors.New("bcrypt error"))
		requestBody := map[string]any{
			"email":    "example@gmail.com",
			"password": "password",
			"name":     "User",
			"birthday": "2000-01-01",
			"address":  "123 Street",
			"gender":   1,
			"role_ids": []uint{1, 2},
		}
		body, _ := json.Marshal(requestBody)

		// Create a test context
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest("POST", "/api/v1/user", bytes.NewBuffer(body))
		c.Set("UserID", uint(1))

		// Call the handler
		handler.CreateUser(c)

		// Assert the response
		expectedBody := map[string]any{
			"code":    float64(apperror.ErrPasswordHashFailed),
			"message": "Failed to hash password",
		}
		var actualBody map[string]any
		json.Unmarshal(w.Body.Bytes(), &actualBody)
		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.Equal(t, expectedBody["code"], actualBody["code"])
		assert.Equal(t, expectedBody["message"], actualBody["message"])

		// Assert mocks
		userService.AssertExpectations(t)
		redisService.AssertExpectations(t)
		bcryptService.AssertExpectations(t)

	})
}

func TestUpdateProfile(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Initialize the validator
	utils.InitValidator()

	t.Run("UpdateProfile - Success", func(t *testing.T) {
		userService := new(mocks.MockUserService)
		redisService := new(mocks.MockRedisService)
		bcryptService := new(mocks.MockBcryptService)
		handler := handlers.NewUserHandler(userService, redisService, bcryptService)

		user := &models.User{
			ID:    1,
			Email: "email@example.com",
			Name:  "User",
		}

		// Assuming the cache key is constructed as "profile:<user_id>"
		profileKey := constants.PROFILE + string(rune(user.ID))

		// Mock the service methods
		userService.On("GetUser", uint(1)).Return(user, nil)
		userService.On("UpdateUser", user).Return(nil)
		redisService.On("Delete", profileKey).Return(nil)

		requestBody := map[string]any{
			"name":     "Updated User",
			"birthday": "2000-01-01",
			"address":  "456 New Street",
			"gender":   2,
		}

		body, _ := json.Marshal(requestBody)

		// Create a test context
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest("PUT", "/api/v1/profile", bytes.NewBuffer(body))
		c.Set("UserID", uint(1))

		// Call the handler
		handler.UpdateProfile(c)

		// Assert the response
		assert.Equal(t, http.StatusOK, w.Code)
		assert.JSONEq(t, `{"message":"Update profile successfully"}`, w.Body.String())

		// Assert mocks
		userService.AssertExpectations(t)
		redisService.AssertExpectations(t)
		bcryptService.AssertExpectations(t)
	})

	t.Run("UpdateProfile - Validation Error", func(t *testing.T) {
		tests := []struct {
			name           string
			reqBody        string
			expectedCode   float64
			expectedMsg    string
			expectedFields []apperror.FieldError
		}{
			{
				name:         "EmptyName",
				reqBody:      `{"name":""}`,
				expectedCode: float64(4001),
				expectedMsg:  "Validation failed",
				expectedFields: []apperror.FieldError{
					{Field: "name", Message: "name must be at least 1 characters long or numeric"},
				},
			},
			{
				name:         "NameNotBlank",
				reqBody:      `{"name":"  "}`,
				expectedCode: float64(4001),
				expectedMsg:  "Validation failed",
				expectedFields: []apperror.FieldError{
					{Field: "name", Message: "name must not be blank"},
				},
			},
			{
				name:         "LongName",
				reqBody:      `{"name": "` + strings.Repeat("a", 46) + `"}`,
				expectedCode: float64(4001),
				expectedMsg:  "Validation failed",
				expectedFields: []apperror.FieldError{
					{Field: "name", Message: "name must be at most 45 characters long or numeric"},
				},
			},
			{
				name:         "InvalidBirthdayFormat",
				reqBody:      `{"name": "User", "birthday": "invalid-date"}`,
				expectedCode: float64(4001),
				expectedMsg:  "Validation failed",
				expectedFields: []apperror.FieldError{
					{Field: "birthday", Message: "birthday must be a valid date (YYYY-MM-DD) and not in the future"},
				},
			},
			{
				name:         "FutureBirthday",
				reqBody:      `{"name": "User", "birthday": "3000-01-01"}`,
				expectedCode: float64(4001),
				expectedMsg:  "Validation failed",
				expectedFields: []apperror.FieldError{
					{Field: "birthday", Message: "birthday must be a valid date (YYYY-MM-DD) and not in the future"},
				},
			},
			{
				name:         "EmptyAddress",
				reqBody:      `{"name": "User", "birthday": "2000-01-01", "address": ""}`,
				expectedCode: float64(4001),
				expectedMsg:  "Validation failed",
				expectedFields: []apperror.FieldError{
					{Field: "address", Message: "address must be at least 1 characters long or numeric"},
				},
			},
			{
				name:         "LongAddress",
				reqBody:      `{"name": "User", "birthday": "2000-01-01", "address": "` + strings.Repeat("a", 256) + `"}`,
				expectedCode: float64(4001),
				expectedMsg:  "Validation failed",
				expectedFields: []apperror.FieldError{
					{Field: "address", Message: "address must be at most 255 characters long or numeric"},
				},
			},
			{
				name:         "AddressNotBlank",
				reqBody:      `{"name": "User", "birthday": "2000-01-01", "address": "  "}`,
				expectedCode: float64(4001),
				expectedMsg:  "Validation failed",
				expectedFields: []apperror.FieldError{
					{Field: "address", Message: "address must not be blank"},
				},
			},
			{
				name:         "InvalidGender 0",
				reqBody:      `{"name": "User", "birthday": "2000-01-01", "address": "123 Street", "gender": 0}`,
				expectedCode: float64(4001),
				expectedMsg:  "Validation failed",
				expectedFields: []apperror.FieldError{
					{Field: "gender", Message: "gender must be one of [1 2 3]"},
				},
			},
			{
				name:         "InvalidGender 4",
				reqBody:      `{"name": "User", "birthday": "2000-01-01", "address": "123 Street", "gender": 4}`,
				expectedCode: float64(4001),
				expectedMsg:  "Validation failed",
				expectedFields: []apperror.FieldError{
					{Field: "gender", Message: "gender must be one of [1 2 3]"},
				},
			},
			{
				name:           "StringGender",
				reqBody:        `{"name": "User", "birthday": "2000-01-01", "address": "123 Street", "gender": "male"}`,
				expectedCode:   float64(4001),
				expectedMsg:    "json: cannot unmarshal string into Go struct field .gender of type int16",
				expectedFields: nil, // specific error case
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				userService := new(mocks.MockUserService)
				redisService := new(mocks.MockRedisService)
				bcryptService := new(mocks.MockBcryptService)
				handler := handlers.NewUserHandler(userService, redisService, bcryptService)

				// Create a test context
				w := httptest.NewRecorder()
				c, _ := gin.CreateTestContext(w)
				c.Request, _ = http.NewRequest("PUT", "/api/v1/profile", bytes.NewBufferString(tt.reqBody))
				c.Set("UserID", uint(1))

				// Call the handler
				handler.UpdateProfile(c)

				// Assert the response
				expectedBody := map[string]any{
					"code":    tt.expectedCode,
					"message": tt.expectedMsg,
					"fields":  tt.expectedFields,
				}

				var actualBody map[string]any
				json.Unmarshal(w.Body.Bytes(), &actualBody)

				assert.Equal(t, http.StatusBadRequest, w.Code)
				assert.Equal(t, expectedBody["code"], actualBody["code"])
				assert.Equal(t, expectedBody["message"], actualBody["message"])
				assert.Equal(t, tt.expectedFields, utils.MapJsonToFieldErrors(actualBody["fields"]))

				// Assert mocks
				userService.AssertExpectations(t)
				redisService.AssertExpectations(t)
				bcryptService.AssertExpectations(t)
			})
		}
	})

	t.Run("UpdateProfile - Invalid UserID ctx", func(t *testing.T) {
		userService := new(mocks.MockUserService)
		redisService := new(mocks.MockRedisService)
		bcryptService := new(mocks.MockBcryptService)
		handler := handlers.NewUserHandler(userService, redisService, bcryptService)

		// Create a test context
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest("PUT", "/api/v1/profile", nil)
		c.Set("UserID", uint(0)) // Invalid User ID

		// Call the handler
		handler.UpdateProfile(c)

		// Assert the response
		var expectedBody = map[string]any{
			"code":    float64(apperror.ErrParseError),
			"message": "Invalid UserID",
		}
		var actualBody map[string]any
		json.Unmarshal(w.Body.Bytes(), &actualBody)
		assert.Equal(t, expectedBody["code"], actualBody["code"])
		assert.Equal(t, expectedBody["message"], actualBody["message"])
		assert.Equal(t, http.StatusBadRequest, w.Code)

		// Assert mocks
		userService.AssertExpectations(t)
		redisService.AssertExpectations(t)
		bcryptService.AssertExpectations(t)
	})

	t.Run("UpdateProfile - User Not Found", func(t *testing.T) {
		userService := new(mocks.MockUserService)
		redisService := new(mocks.MockRedisService)
		bcryptService := new(mocks.MockBcryptService)
		handler := handlers.NewUserHandler(userService, redisService, bcryptService)

		requestBody := map[string]any{
			"name":     "Updated User",
			"birthday": "2000-01-01",
			"address":  "456 New Street",
			"gender":   2,
		}

		body, _ := json.Marshal(requestBody)

		// Mock the service method
		userService.On("GetUser", uint(1)).Return(&models.User{}, apperror.NewNotFoundError("User not found"))

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest("PUT", "/api/v1/profile", bytes.NewBuffer(body))
		c.Set("UserID", uint(1))

		// Call the handler
		handler.UpdateProfile(c)

		// Assert the response
		var expectedBody = map[string]any{
			"code":    float64(apperror.ErrNotFound),
			"message": "User not found",
		}
		var actualBody map[string]any
		json.Unmarshal(w.Body.Bytes(), &actualBody)
		assert.Equal(t, expectedBody["code"], actualBody["code"])
		assert.Equal(t, expectedBody["message"], actualBody["message"])
		assert.Equal(t, http.StatusNotFound, w.Code)

		// Assert mocks
		userService.AssertExpectations(t)
		redisService.AssertExpectations(t)
		bcryptService.AssertExpectations(t)
	})

	t.Run("Error Update User", func(t *testing.T) {
		userService := new(mocks.MockUserService)
		redisService := new(mocks.MockRedisService)
		bcryptService := new(mocks.MockBcryptService)
		handler := handlers.NewUserHandler(userService, redisService, bcryptService)

		user := &models.User{
			ID:    1,
			Email: "email@example.com",
			Name:  "User",
		}

		requestBody := map[string]any{
			"name":     "Updated User",
			"birthday": "2000-01-01",
			"address":  "456 New Street",
			"gender":   2,
		}

		body, _ := json.Marshal(requestBody)

		// Mock the service method
		userService.On("GetUser", uint(1)).Return(user, nil)
		userService.On("UpdateUser", user).Return(apperror.NewDBUpdateError("Update error"))

		// Create a test context
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest("PUT", "/api/v1/profile", bytes.NewBuffer(body))
		c.Set("UserID", uint(1))

		// Call the handler
		handler.UpdateProfile(c)

		// Assert the response
		expectedBody := map[string]any{
			"code":    float64(apperror.ErrDBUpdate),
			"message": "Update error",
		}
		var actualBody map[string]any
		json.Unmarshal(w.Body.Bytes(), &actualBody)
		assert.Equal(t, expectedBody["code"], actualBody["code"])
		assert.Equal(t, expectedBody["message"], actualBody["message"])
		assert.Equal(t, http.StatusInternalServerError, w.Code)

		// Assert mocks
		userService.AssertExpectations(t)
		redisService.AssertExpectations(t)
		bcryptService.AssertExpectations(t)
	})

	t.Run("Error Delete Cache", func(t *testing.T) {
		userService := new(mocks.MockUserService)
		redisService := new(mocks.MockRedisService)
		bcryptService := new(mocks.MockBcryptService)
		handler := handlers.NewUserHandler(userService, redisService, bcryptService)

		user := &models.User{
			ID:    1,
			Email: "email@example.com",
			Name:  "User",
		}
		requestBody := map[string]any{
			"name":     "Updated User",
			"birthday": "2000-01-01",
			"address":  "456 New Street",
			"gender":   2,
		}

		body, _ := json.Marshal(requestBody)
		// Mock the service method
		userService.On("GetUser", uint(1)).Return(user, nil)
		userService.On("UpdateUser", user).Return(nil)
		profileKey := constants.PROFILE + string(rune(user.ID))
		redisService.On("Delete", profileKey).Return(errors.New("Redis delete error"))

		// Create a test context
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest("PUT", "/api/v1/profile", bytes.NewBuffer(body))
		c.Set("UserID", uint(1))

		// Call the handler
		handler.UpdateProfile(c)

		// Assert the response
		expectedBody := map[string]any{
			"message": "Update profile successfully",
		}
		var actualBody map[string]any
		json.Unmarshal(w.Body.Bytes(), &actualBody)
		assert.Equal(t, expectedBody["message"], actualBody["message"])
		assert.Equal(t, http.StatusOK, w.Code)

		// Assert mocks
		userService.AssertExpectations(t)
		redisService.AssertExpectations(t)
		bcryptService.AssertExpectations(t)
	})
}

func TestGetProfile(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("Success get profile from database", func(t *testing.T) {
		userService := new(mocks.MockUserService)
		redisService := new(mocks.MockRedisService)
		bcryptService := new(mocks.MockBcryptService)
		handler := handlers.NewUserHandler(userService, redisService, bcryptService)

		user := &models.User{
			ID:        1,
			Email:     "email@example.com",
			Name:      "User",
			Gender:    1,
			CreatedAt: time.Date(2023, 10, 1, 0, 0, 0, 0, time.UTC),
			UpdatedAt: time.Date(2023, 10, 1, 0, 0, 0, 0, time.UTC),
		}
		profileKey := constants.PROFILE + string(rune(user.ID))
		// Mock the service method
		userService.On("GetProfile", uint(1)).Return(user, nil)
		redisService.On("Get", profileKey).Return("", nil)
		// Parse the user into a JSON string
		profileData, _ := json.Marshal(user)
		// Set the TTL for the cache
		ttl := 60 * time.Minute
		// Mock the Redis Set method to cache the profile
		redisService.On("Set", profileKey, profileData, ttl).Return(nil)

		// Create a test context
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest("GET", "/api/v1/profile", nil)
		c.Set("UserID", uint(1))

		// Call the handler
		handler.GetProfile(c)

		// Assert the response
		expectedBody := map[string]any{
			"id":        float64(1),
			"email":     "email@example.com",
			"name":      "User",
			"gender":    float64(1),
			"createdAt": "2023-10-01T00:00:00Z",
			"updatedAt": "2023-10-01T00:00:00Z",
			"deletedAt": nil,
		}
		var actualBody map[string]any
		json.Unmarshal(w.Body.Bytes(), &actualBody)
		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, expectedBody, actualBody)

		// Assert mocks
		userService.AssertExpectations(t)
		redisService.AssertExpectations(t)
		bcryptService.AssertExpectations(t)
	})

	t.Run("Success get profile from redis cache", func(t *testing.T) {
		userService := new(mocks.MockUserService)
		redisService := new(mocks.MockRedisService)
		bcryptService := new(mocks.MockBcryptService)

		user := &models.User{
			ID:        1,
			Email:     "email@example.com",
			Name:      "User",
			Gender:    1,
			CreatedAt: time.Date(2023, 10, 1, 0, 0, 0, 0, time.UTC),
			UpdatedAt: time.Date(2023, 10, 1, 0, 0, 0, 0, time.UTC),
		}
		profileKey := constants.PROFILE + string(rune(user.ID))
		// Mock the Redis Get method to return a cached profile
		cachedProfile := fmt.Sprintf(`{"id":%d,"email":"%s","name":"%s","gender":%d,"createdAt":"%s","updatedAt":"%s","deletedAt":null}`,
			user.ID, user.Email, user.Name, user.Gender, user.CreatedAt.Format(time.RFC3339), user.UpdatedAt.Format(time.RFC3339))

		redisService.On("Get", profileKey).Return(cachedProfile, nil)

		handler := handlers.NewUserHandler(userService, redisService, bcryptService)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest("GET", "/api/v1/profile", nil)
		c.Set("UserID", uint(1))

		// Call the handler
		handler.GetProfile(c)

		// Assert the response
		expectedBody := map[string]any{
			"id":        float64(1),
			"email":     "email@example.com",
			"name":      "User",
			"gender":    float64(1),
			"createdAt": "2023-10-01T00:00:00Z",
			"updatedAt": "2023-10-01T00:00:00Z",
			"deletedAt": nil,
		}

		var actualBody map[string]any
		json.Unmarshal(w.Body.Bytes(), &actualBody)
		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, expectedBody, actualBody)

		// Assert mocks
		userService.AssertExpectations(t)
		redisService.AssertExpectations(t)
		bcryptService.AssertExpectations(t)
	})

	t.Run("Error Invalid User ID", func(t *testing.T) {
		userService := new(mocks.MockUserService)
		redisService := new(mocks.MockRedisService)
		bcryptService := new(mocks.MockBcryptService)

		handler := handlers.NewUserHandler(userService, redisService, bcryptService)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest("GET", "/api/v1/profile", nil)
		c.Set("UserID", uint(0)) // Invalid User ID

		// Call the GetProfile handler
		handler.GetProfile(c)

		var expectedBody = map[string]any{
			"code":    float64(apperror.ErrParseError),
			"message": "Invalid UserID",
		}
		var actualBody map[string]any
		json.Unmarshal(w.Body.Bytes(), &actualBody)
		assert.Equal(t, expectedBody["code"], actualBody["code"])
		assert.Equal(t, expectedBody["message"], actualBody["message"])
		assert.Equal(t, http.StatusBadRequest, w.Code)

		// Assert mocks
		userService.AssertExpectations(t)
		redisService.AssertExpectations(t)
		bcryptService.AssertExpectations(t)
	})

	t.Run("Error User Not Found", func(t *testing.T) {
		userService := new(mocks.MockUserService)
		redisService := new(mocks.MockRedisService)
		bcryptService := new(mocks.MockBcryptService)

		profileKey := constants.PROFILE + string(rune(1))
		// Mock the GetUser method to return an error
		userService.On("GetProfile", uint(1)).Return(&models.User{}, apperror.NewNotFoundError("User not found"))
		// Mock the Redis Get method to return an empty string
		redisService.On("Get", profileKey).Return("", nil)

		handler := handlers.NewUserHandler(userService, redisService, bcryptService)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest("GET", "/api/v1/profile", nil)
		c.Set("UserID", uint(1))

		// Call the GetProfile handler
		handler.GetProfile(c)

		// Assert the response
		expectedBody := map[string]any{
			"code":    float64(apperror.ErrNotFound),
			"message": "User not found",
		}
		var actualBody map[string]any
		json.Unmarshal(w.Body.Bytes(), &actualBody)
		assert.Equal(t, expectedBody["code"], actualBody["code"])
		assert.Equal(t, expectedBody["message"], actualBody["message"])
		assert.Equal(t, http.StatusNotFound, w.Code)

		// Assert mocks
		userService.AssertExpectations(t)
		redisService.AssertExpectations(t)
		bcryptService.AssertExpectations(t)

	})

	t.Run("Success Get Profile but Error Cache", func(t *testing.T) {
		userService := new(mocks.MockUserService)
		redisService := new(mocks.MockRedisService)
		bcryptService := new(mocks.MockBcryptService)

		profileKey := constants.PROFILE + string(rune(1))
		// Mock the GetUser method to return a user
		user := &models.User{
			ID:        1,
			Email:     "email@example.com",
			Name:      "User",
			Gender:    1,
			CreatedAt: time.Date(2023, 10, 1, 0, 0, 0, 0, time.UTC),
			UpdatedAt: time.Date(2023, 10, 1, 0, 0, 0, 0, time.UTC),
		}
		userService.On("GetProfile", uint(1)).Return(user, nil)
		// Mock the Redis Get method to return an error
		redisService.On("Get", profileKey).Return("", errors.New("Cache get error"))
		// Mock the Redis Set method to cache the profile
		profileData, _ := json.Marshal(user)
		ttl := 60 * time.Minute
		redisService.On("Set", profileKey, profileData, ttl).Return(errors.New("Cache set error"))
		handler := handlers.NewUserHandler(userService, redisService, bcryptService)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest("GET", "/api/v1/profile", nil)
		c.Set("UserID", uint(1))
		// Call the GetProfile handler
		handler.GetProfile(c)
		// Assert the response
		assert.Equal(t, http.StatusOK, w.Code)
		expectedBody := map[string]any{
			"id":        float64(1),
			"email":     "email@example.com",
			"name":      "User",
			"gender":    float64(1),
			"createdAt": "2023-10-01T00:00:00Z",
			"updatedAt": "2023-10-01T00:00:00Z",
			"deletedAt": nil,
		}

		var actualBody map[string]any
		err := json.Unmarshal(w.Body.Bytes(), &actualBody)
		require.NoError(t, err)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, expectedBody, actualBody)

	})

	t.Run("GetProfile - Could not parse user data from cache", func(t *testing.T) {
		userService := new(mocks.MockUserService)
		redisService := new(mocks.MockRedisService)
		bcryptService := new(mocks.MockBcryptService)

		profileKey := constants.PROFILE + string(rune(1))
		// Mock the Redis Get method to return an invalid JSON
		redisService.On("Get", profileKey).Return("invalid-json", nil)

		handler := handlers.NewUserHandler(userService, redisService, bcryptService)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest("GET", "/api/v1/profile", nil)
		c.Set("UserID", uint(1))

		// Call the GetProfile handler
		handler.GetProfile(c)

		// Assert the response
		actualBody := map[string]any{
			"code":    float64(apperror.ErrParseError),
			"message": "Invalid user data in cache",
		}
		var responseBody map[string]any
		json.Unmarshal(w.Body.Bytes(), &responseBody)
		assert.Equal(t, actualBody["code"], responseBody["code"])
		assert.Equal(t, actualBody["message"], responseBody["message"])
		assert.Equal(t, http.StatusBadRequest, w.Code)

		// Assert mocks
		userService.AssertExpectations(t)
		redisService.AssertExpectations(t)
		bcryptService.AssertExpectations(t)
	})
}

func TestGetUser(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("GetUser - Success", func(t *testing.T) {
		userService := new(mocks.MockUserService)
		redisService := new(mocks.MockRedisService)
		bcryptService := new(mocks.MockBcryptService)
		handler := handlers.NewUserHandler(userService, redisService, bcryptService)

		user := &models.User{
			ID:        1,
			Email:     "email@example.com",
			Name:      "User",
			Gender:    1,
			CreatedAt: time.Date(2023, 10, 1, 0, 0, 0, 0, time.UTC),
			UpdatedAt: time.Date(2023, 10, 1, 0, 0, 0, 0, time.UTC),
		}

		// Mock the service method
		userService.On("GetUser", uint(1)).Return(user, nil)

		// Create http request with valid UserID
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest("GET", "/users/:id", nil)
		c.Params = gin.Params{gin.Param{Key: "id", Value: "1"}}

		// Call the GetUser handler
		handler.GetUser(c)

		// Assert the response
		var actualBody map[string]any
		json.Unmarshal(w.Body.Bytes(), &actualBody)
		expectedBody := map[string]any{
			"id":        float64(1),
			"email":     "email@example.com",
			"name":      "User",
			"gender":    float64(1),
			"createdAt": "2023-10-01T00:00:00Z",
			"updatedAt": "2023-10-01T00:00:00Z",
			"deletedAt": nil,
		}
		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, expectedBody, actualBody)

		// Assert mocks
		userService.AssertExpectations(t)
		redisService.AssertExpectations(t)
		bcryptService.AssertExpectations(t)
	})

	t.Run("GetUser - Not found the user", func(t *testing.T) {
		userService := new(mocks.MockUserService)
		redisService := new(mocks.MockRedisService)
		bcryptService := new(mocks.MockBcryptService)
		handler := handlers.NewUserHandler(userService, redisService, bcryptService)

		// Mock the service method
		userService.On("GetUser", uint(1)).Return(&models.User{}, apperror.NewNotFoundError("User not found"))

		// Create http request with valid UserID
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest("GET", "/users/:id", nil)
		c.Params = gin.Params{gin.Param{Key: "id", Value: "1"}}

		// Call the GetUser handler
		handler.GetUser(c)

		// Assert the response
		actualBody := map[string]any{
			"code":    float64(apperror.ErrNotFound),
			"message": "User not found",
		}
		var responseBody map[string]any
		json.Unmarshal(w.Body.Bytes(), &responseBody)
		assert.Equal(t, http.StatusNotFound, w.Code)
		assert.Equal(t, actualBody["code"], responseBody["code"])
		assert.Equal(t, actualBody["message"], responseBody["message"])

		// Assert mocks
		userService.AssertExpectations(t)
		redisService.AssertExpectations(t)
		bcryptService.AssertExpectations(t)
	})

	t.Run("GetUser - Invalid UserID", func(t *testing.T) {
		userService := new(mocks.MockUserService)
		redisService := new(mocks.MockRedisService)
		bcryptService := new(mocks.MockBcryptService)
		handler := handlers.NewUserHandler(userService, redisService, bcryptService)

		// Create http request with invalid UserID
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest("GET", "/users/:id", nil)
		c.Params = gin.Params{gin.Param{Key: "id", Value: "invalid-id"}}

		// Call the GetUser handler
		handler.GetUser(c)

		// Assert the response
		expectedBody := map[string]any{
			"code":    float64(apperror.ErrParseError),
			"message": "Invalid UserID",
		}
		var actualBody map[string]any
		json.Unmarshal(w.Body.Bytes(), &actualBody)
		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Equal(t, expectedBody["code"], actualBody["code"])
		assert.Equal(t, expectedBody["message"], actualBody["message"])

		// Assert mocks
		userService.AssertExpectations(t)
		redisService.AssertExpectations(t)
		bcryptService.AssertExpectations(t)
	})
}

func TestChangePassword(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Initialize the validator
	utils.InitValidator()

	t.Run("ChangePassword - Success", func(t *testing.T) {
		userService := new(mocks.MockUserService)
		redisService := new(mocks.MockRedisService)
		bcryptService := new(mocks.MockBcryptService)
		handler := handlers.NewUserHandler(userService, redisService, bcryptService)

		user := &models.User{
			ID:        1,
			Email:     "email@example.com",
			Name:      "User",
			Password:  "$2a$10$I/L5VegpCyOlJPoa1.KrmeCdezSBIandsEL5S2dd4Ap0YIWk0Iuka", // bcrypt hash of "12345678"
			Gender:    1,
			CreatedAt: time.Date(2023, 10, 1, 0, 0, 0, 0, time.UTC),
			UpdatedAt: time.Date(2023, 10, 1, 0, 0, 0, 0, time.UTC),
		}
		requestBody := map[string]any{
			"old_password":     "12345678",
			"new_password":     "newpassword",
			"confirm_password": "newpassword",
		}
		body, _ := json.Marshal(requestBody)

		// Mock the services methods
		userService.On("GetUser", uint(1)).Return(user, nil)
		userService.On("UpdateUser", user).Return(nil)
		bcryptService.On("CheckPasswordHash", "12345678", user.Password).Return(true)
		bcryptService.On("HashPassword", "newpassword").Return("$2a$10$hashedNewPassword", nil)

		// Create http request and context
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest("PUT", "/api/v1/change-password", bytes.NewBuffer(body))
		c.Set("UserID", uint(1))

		// Call the ChangePassword handler
		handler.ChangePassword(c)

		// Assert the response
		assert.Equal(t, http.StatusOK, w.Code)
		assert.JSONEq(t, `{"message":"Change password successfully"}`, w.Body.String())

		// Assert mocks
		userService.AssertExpectations(t)
		redisService.AssertExpectations(t)
		bcryptService.AssertExpectations(t)
	})

	t.Run("ChangePassword - Validation Error", func(t *testing.T) {
		tests := []struct {
			name           string
			reqBody        string
			expectedCode   float64
			expectedMsg    string
			expectedFields []apperror.FieldError
		}{
			{
				name:         "EmptyOldPassword",
				reqBody:      `{"old_password":""}`,
				expectedCode: float64(4001),
				expectedMsg:  "Validation failed",
				expectedFields: []apperror.FieldError{
					{Field: "old_password", Message: "old_password is required"},
					{Field: "new_password", Message: "new_password is required"},
					{Field: "confirm_password", Message: "confirm_password is required"},
				},
			},
			{
				name:         "ShortOldPassword",
				reqBody:      `{"old_password":"short"}`,
				expectedCode: float64(4001),
				expectedMsg:  "Validation failed",
				expectedFields: []apperror.FieldError{
					{Field: "old_password", Message: "old_password must be at least 6 characters long or numeric"},
					{Field: "new_password", Message: "new_password is required"},
					{Field: "confirm_password", Message: "confirm_password is required"},
				},
			},
			{
				name:         "LongOldPassword",
				reqBody:      `{"old_password":"` + strings.Repeat("a", 256) + `"}`,
				expectedCode: float64(4001),
				expectedMsg:  "Validation failed",
				expectedFields: []apperror.FieldError{
					{Field: "old_password", Message: "old_password must be at most 255 characters long or numeric"},
					{Field: "new_password", Message: "new_password is required"},
					{Field: "confirm_password", Message: "confirm_password is required"},
				},
			},
			{
				name:         "EmptyNewPassword",
				reqBody:      `{"old_password":"12345678","new_password":""}`,
				expectedCode: float64(4001),
				expectedMsg:  "Validation failed",
				expectedFields: []apperror.FieldError{
					{Field: "new_password", Message: "new_password is required"},
					{Field: "confirm_password", Message: "confirm_password is required"},
				},
			},
			{
				name:         "ShortNewPassword",
				reqBody:      `{"old_password":"12345678","new_password":"short"}`,
				expectedCode: float64(4001),
				expectedMsg:  "Validation failed",
				expectedFields: []apperror.FieldError{
					{Field: "new_password", Message: "new_password must be at least 6 characters long or numeric"},
					{Field: "confirm_password", Message: "confirm_password is required"},
				},
			},
			{
				name:         "LongNewPassword",
				reqBody:      `{"old_password":"12345678","new_password":"` + strings.Repeat("a", 256) + `"}`,
				expectedCode: float64(4001),
				expectedMsg:  "Validation failed",
				expectedFields: []apperror.FieldError{
					{Field: "new_password", Message: "new_password must be at most 255 characters long or numeric"},
					{Field: "confirm_password", Message: "confirm_password is required"},
				},
			},
			{
				name:         "EmptyConfirmPassword",
				reqBody:      `{"old_password":"12345678","new_password":"newpassword","confirm_password":""}`,
				expectedCode: float64(4001),
				expectedMsg:  "Validation failed",
				expectedFields: []apperror.FieldError{
					{Field: "confirm_password", Message: "confirm_password is required"},
				},
			},
			{
				name:         "ShortConfirmPassword",
				reqBody:      `{"old_password":"12345678","new_password":"newpassword","confirm_password":"short"}`,
				expectedCode: float64(4001),
				expectedMsg:  "Validation failed",
				expectedFields: []apperror.FieldError{
					{Field: "confirm_password", Message: "confirm_password must be at least 6 characters long or numeric"},
				},
			},
			{
				name:         "LongConfirmPassword",
				reqBody:      `{"old_password":"12345678","new_password":"newpassword","confirm_password":"` + strings.Repeat("a", 256) + `"}`,
				expectedCode: float64(4001),
				expectedMsg:  "Validation failed",
				expectedFields: []apperror.FieldError{
					{Field: "confirm_password", Message: "confirm_password must be at most 255 characters long or numeric"},
				},
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				userService := new(mocks.MockUserService)
				redisService := new(mocks.MockRedisService)
				bcryptService := new(mocks.MockBcryptService)
				handler := handlers.NewUserHandler(userService, redisService, bcryptService)

				// Create http request and context
				w := httptest.NewRecorder()
				c, _ := gin.CreateTestContext(w)
				c.Request, _ = http.NewRequest("PUT", "/api/v1/change-password", bytes.NewBufferString(tt.reqBody))
				c.Set("UserID", uint(1))

				// Call the ChangePassword handler
				handler.ChangePassword(c)

				// Assert the response
				expectedBody := map[string]any{
					"code":    tt.expectedCode,
					"message": tt.expectedMsg,
					"fields":  tt.expectedFields,
				}
				var actualBody map[string]any
				json.Unmarshal(w.Body.Bytes(), &actualBody)
				assert.Equal(t, http.StatusBadRequest, w.Code)
				assert.Equal(t, expectedBody["code"], actualBody["code"])
				assert.Equal(t, expectedBody["message"], actualBody["message"])
				assert.Equal(t, expectedBody["fields"], utils.MapJsonToFieldErrors(actualBody["fields"]))

				// Assert mock expectations
				userService.AssertExpectations(t)
				bcryptService.AssertExpectations(t)
				redisService.AssertExpectations(t)
			})
		}
	})

	t.Run("ChangePassword - NotFound User", func(t *testing.T) {
		userService := new(mocks.MockUserService)
		redisService := new(mocks.MockRedisService)
		bcryptService := new(mocks.MockBcryptService)
		handler := handlers.NewUserHandler(userService, redisService, bcryptService)

		requestBody := map[string]any{
			"old_password":     "12345678",
			"new_password":     "newpassword",
			"confirm_password": "newpassword",
		}
		body, _ := json.Marshal(requestBody)

		// Mock the GetUser method to return an error
		userService.On("GetUser", uint(1)).Return(&models.User{}, apperror.NewNotFoundError("User not found"))

		// Create http request and context
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest("PUT", "/api/v1/change-password", bytes.NewBuffer(body))
		c.Set("UserID", uint(1))

		// Call the ChangePassword handler
		handler.ChangePassword(c)

		// Assert the response
		expectedBody := map[string]any{
			"code":    float64(apperror.ErrNotFound),
			"message": "User not found",
		}
		var actualBody map[string]any
		json.Unmarshal(w.Body.Bytes(), &actualBody)
		assert.Equal(t, http.StatusNotFound, w.Code)
		assert.Equal(t, expectedBody["code"], actualBody["code"])
		assert.Equal(t, expectedBody["message"], actualBody["message"])

		// Assert mock expectations
		userService.AssertExpectations(t)
		bcryptService.AssertExpectations(t)
		redisService.AssertExpectations(t)
	})

	t.Run("ChangePassword - Old Password Mismatch", func(t *testing.T) {
		userService := new(mocks.MockUserService)
		redisService := new(mocks.MockRedisService)
		bcryptService := new(mocks.MockBcryptService)
		handler := handlers.NewUserHandler(userService, redisService, bcryptService)

		user := &models.User{
			ID:       1,
			Email:    "",
			Name:     "",
			Password: "$2a$10$I/L5VegpCyOlJPoa1.KrmeCdezSBIandsEL5S2dd4Ap0YIWk0Iuka", // bcrypt hash of "12345678"
		}
		requestBody := map[string]any{
			"old_password":     "wrongpassword",
			"new_password":     "newpassword",
			"confirm_password": "newpassword",
		}
		body, _ := json.Marshal(requestBody)

		// Mock the service methods
		userService.On("GetUser", uint(1)).Return(user, nil)
		bcryptService.On("CheckPasswordHash", "wrongpassword", user.Password).Return(false)

		// Create a new UserHandler instance
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest("PUT", "/api/v1/change-password", bytes.NewBuffer(body))
		c.Set("UserID", uint(1))

		// Call the ChangePassword handler
		handler.ChangePassword(c)

		// Assert the response
		expectedBody := map[string]any{
			"code":    float64(apperror.ErrInvalidPassword),
			"message": "Old password is incorrect",
		}

		var actualBody map[string]any
		json.Unmarshal(w.Body.Bytes(), &actualBody)
		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Equal(t, expectedBody["code"], actualBody["code"])
		assert.Equal(t, expectedBody["message"], actualBody["message"])

		// Assert mock expectations
		userService.AssertExpectations(t)
		bcryptService.AssertExpectations(t)
		redisService.AssertExpectations(t)
	})

	t.Run("ChangePassword - New Password and Confirm Password Mismatch", func(t *testing.T) {
		userService := new(mocks.MockUserService)
		redisService := new(mocks.MockRedisService)
		bcryptService := new(mocks.MockBcryptService)
		handler := handlers.NewUserHandler(userService, redisService, bcryptService)

		user := &models.User{
			ID:       1,
			Email:    "email@example.com",
			Name:     "John Doe",
			Password: "$2a$10$I/L5VegpCyOlJPoa1.KrmeCdezSBIandsEL5S2dd4Ap0YIWk0Iuka", // bcrypt hash of "12345678"
		}
		requestBody := map[string]any{
			"old_password":     "12345678",
			"new_password":     "123456789",
			"confirm_password": "differentpassword",
		}
		body, _ := json.Marshal(requestBody)

		// Mock the service methods
		userService.On("GetUser", uint(1)).Return(user, nil)
		bcryptService.On("CheckPasswordHash", requestBody["old_password"], user.Password).Return(true)

		// Create test context
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest("PUT", "/api/v1/change-password", bytes.NewBuffer(body))
		c.Set("UserID", uint(1))

		// Call the ChangePassword handler
		handler.ChangePassword(c)

		// Assert the response
		expectedBody := map[string]any{
			"code":    float64(apperror.ErrPasswordMismatch),
			"message": "New password and confirm password do not match",
		}
		var actualBody map[string]any
		json.Unmarshal(w.Body.Bytes(), &actualBody)
		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Equal(t, expectedBody["code"], actualBody["code"])
		assert.Equal(t, expectedBody["message"], actualBody["message"])

		// Assert mock expectations
		userService.AssertExpectations(t)
		bcryptService.AssertExpectations(t)
		redisService.AssertExpectations(t)
	})

	t.Run("ChangePassword - Failed To Update", func(t *testing.T) {
		userService := new(mocks.MockUserService)
		redisService := new(mocks.MockRedisService)
		bcryptService := new(mocks.MockBcryptService)
		handler := handlers.NewUserHandler(userService, redisService, bcryptService)

		user := &models.User{
			ID:       1,
			Email:    "email@example.com",
			Name:     "User",
			Password: "$2a$10$I/L5VegpCyOlJPoa1.KrmeCdezSBIandsEL5S2dd4Ap0YIWk0Iuka", // bcrypt hash of "12345678"
		}
		requestBody := map[string]any{
			"old_password":     "12345678",
			"new_password":     "newpassword",
			"confirm_password": "newpassword",
		}
		body, _ := json.Marshal(requestBody)

		// Mock the service methods
		userService.On("GetUser", uint(1)).Return(user, nil)
		userService.On("UpdateUser", user).Return(apperror.NewDBUpdateError("Update error"))
		bcryptService.On("CheckPasswordHash", "12345678", user.Password).Return(true)
		bcryptService.On("HashPassword", "newpassword").Return("$2a$10$hashedNewPassword", nil)

		// Create a test context
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest("PUT", "/api/v1/change-password", bytes.NewBuffer(body))
		c.Set("UserID", uint(1))

		// Call the handler
		handler.ChangePassword(c)

		// Assert the response
		expectedBody := map[string]any{
			"code":    float64(apperror.ErrDBUpdate),
			"message": "Update error",
		}
		var actualBody map[string]any
		json.Unmarshal(w.Body.Bytes(), &actualBody)
		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.Equal(t, expectedBody["code"], actualBody["code"])
		assert.Equal(t, expectedBody["message"], actualBody["message"])

		// Assert mock expectations
		userService.AssertExpectations(t)
		bcryptService.AssertExpectations(t)
		redisService.AssertExpectations(t)
	})

	t.Run("ChangePassword - User Not found from ctx", func(t *testing.T) {
		userService := new(mocks.MockUserService)
		redisService := new(mocks.MockRedisService)
		bcryptService := new(mocks.MockBcryptService)
		handler := handlers.NewUserHandler(userService, redisService, bcryptService)

		// Create a test context
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest("PUT", "/api/v1/change-password", nil)
		c.Set("UserID", uint(0)) // Invalid User ID

		// Call the handler
		handler.ChangePassword(c)

		// Assert the response
		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.JSONEq(t, fmt.Sprintf(`{"code":%d,"message":"Invalid UserID"}`, apperror.ErrParseError), w.Body.String())

		// Assert mocks
		userService.AssertExpectations(t)
		bcryptService.AssertExpectations(t)
		redisService.AssertExpectations(t)
	})

	t.Run("ChangePassword - Old Password equal to New Password", func(t *testing.T) {
		userService := new(mocks.MockUserService)
		redisService := new(mocks.MockRedisService)
		bcryptService := new(mocks.MockBcryptService)
		handler := handlers.NewUserHandler(userService, redisService, bcryptService)

		user := &models.User{
			ID:       1,
			Email:    "email@example.com",
			Name:     "User",
			Password: "$2a$10$I/L5VegpCyOlJPoa1.KrmeCdezSBIandsEL5S2dd4Ap0YIWk0Iuka", // bcrypt hash of "12345678"
		}
		requestBody := map[string]any{
			"old_password":     "12345678",
			"new_password":     "12345678",
			"confirm_password": "12345678",
		}
		body, _ := json.Marshal(requestBody)

		// Mock the service methods
		userService.On("GetUser", uint(1)).Return(user, nil)
		bcryptService.On("CheckPasswordHash", requestBody["old_password"], user.Password).Return(true)

		// Create a test context
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest("PUT", "/api/v1/change-password", bytes.NewBuffer(body))
		c.Set("UserID", uint(1))

		// Call the handler
		handler.ChangePassword(c)

		// Assert the response
		expectedBody := map[string]any{
			"code":    float64(apperror.ErrPasswordMismatch),
			"message": "New password must be different from old password",
		}
		var actualBody map[string]any
		json.Unmarshal(w.Body.Bytes(), &actualBody)
		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Equal(t, expectedBody["code"], actualBody["code"])
		assert.Equal(t, expectedBody["message"], actualBody["message"])

		// Assert mocks
		userService.AssertExpectations(t)
		bcryptService.AssertExpectations(t)
		redisService.AssertExpectations(t)
	})

	t.Run("ChangePassword - Hash Password Failed", func(t *testing.T) {
		userService := new(mocks.MockUserService)
		redisService := new(mocks.MockRedisService)
		bcryptService := new(mocks.MockBcryptService)
		handler := handlers.NewUserHandler(userService, redisService, bcryptService)

		user := &models.User{
			ID:       1,
			Email:    "email@example.com",
			Name:     "User",
			Password: "$2a$10$I/L5VegpCyOlJPoa1.KrmeCdezSBIandsEL5S2dd4Ap0YIWk0Iuka", // bcrypt hash of "12345678"
		}
		requestBody := map[string]any{
			"old_password":     "12345678",
			"new_password":     "newpassword",
			"confirm_password": "newpassword",
		}
		body, _ := json.Marshal(requestBody)

		// Mock the service methods
		userService.On("GetUser", uint(1)).Return(user, nil)
		bcryptService.On("CheckPasswordHash", requestBody["old_password"], user.Password).Return(true)
		bcryptService.On("HashPassword", "newpassword").Return("", apperror.NewInternalError("Hash password failed"))

		// Create a test context
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest("PUT", "/api/v1/change-password", bytes.NewBuffer(body))
		c.Set("UserID", uint(1))

		// Call the handler
		handler.ChangePassword(c)

		// Assert the response
		expectedBody := map[string]any{
			"code":    float64(apperror.ErrInternal),
			"message": "Hash password failed",
		}
		var actualBody map[string]any
		json.Unmarshal(w.Body.Bytes(), &actualBody)
		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.Equal(t, expectedBody["code"], actualBody["code"])
		assert.Equal(t, expectedBody["message"], actualBody["message"])

		// Assert mock expectations
		userService.AssertExpectations(t)
		bcryptService.AssertExpectations(t)
		redisService.AssertExpectations(t)
	})

}

func TestUpdateUser(t *testing.T) {
	gin.SetMode(gin.TestMode)

	utils.InitValidator()

	t.Run("UpdateUser - Success", func(t *testing.T) {
		// Mock the dependencies
		userService := new(mocks.MockUserService)
		redisService := new(mocks.MockRedisService)
		bcryptService := new(mocks.MockBcryptService)
		handler := handlers.NewUserHandler(userService, redisService, bcryptService)

		user := &models.User{
			ID:    1,
			Email: "email@example.com",
			Name:  "User",
		}

		requestBody := map[string]any{
			"name":     "Updated User",
			"birthday": "2000-01-01",
			"address":  "456 New Street",
			"gender":   1,
		}
		body, _ := json.Marshal(requestBody)

		// Mock methods
		userService.On("GetUser", uint(1)).Return(user, nil)
		userService.On("UpdateUser", user).Return(nil)

		// Create a test context
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest("PATCH", "/api/v1/users/id", bytes.NewBuffer(body))
		c.Params = gin.Params{gin.Param{Key: "id", Value: "1"}}

		// Call the  handler
		handler.UpdateUser(c)

		// Assert the response
		assert.Equal(t, http.StatusOK, w.Code)
		assert.JSONEq(t, `{"message":"Update user successfully"}`, w.Body.String())
	})
	t.Run("UpdateUser - Validation Error", func(t *testing.T) {
		tests := []struct {
			name           string
			reqBody        string
			expectedCode   float64
			expectedMsg    string
			expectedFields []apperror.FieldError
		}{
			{
				name:         "EmptyName",
				reqBody:      `{"name":""}`,
				expectedCode: 4001,
				expectedMsg:  "Validation failed",
				expectedFields: []apperror.FieldError{
					{Field: "name", Message: "name must be at least 1 characters long or numeric"},
				},
			},
			{
				name:         "NameNotBlank",
				reqBody:      `{"name":"  "}`,
				expectedCode: 4001,
				expectedMsg:  "Validation failed",
				expectedFields: []apperror.FieldError{
					{Field: "name", Message: "name must not be blank"},
				},
			},
			{
				name:         "LongName",
				reqBody:      `{"name": "` + strings.Repeat("a", 46) + `"}`,
				expectedCode: float64(4001),
				expectedMsg:  "Validation failed",
				expectedFields: []apperror.FieldError{
					{Field: "name", Message: "name must be at most 45 characters long or numeric"},
				},
			},
			{
				name:         "InvalidBirthdayFormat",
				reqBody:      `{"name":"User","birthday":"invalid-date"}`,
				expectedCode: 4001,
				expectedMsg:  "Validation failed",
				expectedFields: []apperror.FieldError{
					{Field: "birthday", Message: "birthday must be a valid date (YYYY-MM-DD) and not in the future"},
				},
			},
			{
				name:         "FutureBirthday",
				reqBody:      `{"name":"User","birthday":"3000-01-01"}`,
				expectedCode: 4001,
				expectedMsg:  "Validation failed",
				expectedFields: []apperror.FieldError{
					{Field: "birthday", Message: "birthday must be a valid date (YYYY-MM-DD) and not in the future"},
				},
			},
			{
				name:         "EmptyAddress",
				reqBody:      `{"name":"User","birthday":"2000-01-01","address":""}`,
				expectedCode: 4001,
				expectedMsg:  "Validation failed",
				expectedFields: []apperror.FieldError{
					{Field: "address", Message: "address must be at least 1 characters long or numeric"},
				},
			},
			{
				name:         "LongAddress",
				reqBody:      `{"name":"User","birthday":"2000-01-01","address":"` + strings.Repeat("a", 256) + `"}`,
				expectedCode: 4001,
				expectedMsg:  "Validation failed",
				expectedFields: []apperror.FieldError{
					{Field: "address", Message: "address must be at most 255 characters long or numeric"},
				},
			},
			{
				name:         "AddressNotBlank",
				reqBody:      `{"name":"User","birthday":"2000-01-01","address":"  "}`,
				expectedCode: 4001,
				expectedMsg:  "Validation failed",
				expectedFields: []apperror.FieldError{
					{Field: "address", Message: "address must not be blank"},
				},
			},
			{
				name:         "InvalidGender 4",
				reqBody:      `{"name":"User","birthday":"2000-01-01","address":"123 Street","gender":4}`,
				expectedCode: 4001,
				expectedMsg:  "Validation failed",
				expectedFields: []apperror.FieldError{
					{Field: "gender", Message: "gender must be one of [1 2 3]"},
				},
			},
			{
				name:         "InvalidGender 0",
				reqBody:      `{"name":"User","birthday":"2000-01-01","address":"123 Street","gender":0}`,
				expectedCode: 4001,
				expectedMsg:  "Validation failed",
				expectedFields: []apperror.FieldError{
					{Field: "gender", Message: "gender must be one of [1 2 3]"},
				},
			},
			{
				name:           "StringGender",
				reqBody:        `{"name":"User","birthday":"2000-01-01","address":"123 Street","gender":"male"}`,
				expectedCode:   4001,
				expectedMsg:    "json: cannot unmarshal string into Go struct field .gender of type int16",
				expectedFields: nil,
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				// Mock services
				userService := new(mocks.MockUserService)
				redisService := new(mocks.MockRedisService)
				bcryptService := new(mocks.MockBcryptService)
				handler := handlers.NewUserHandler(userService, redisService, bcryptService)

				// Create a test context
				w := httptest.NewRecorder()
				c, _ := gin.CreateTestContext(w)
				c.Request, _ = http.NewRequest("PATCH", "/api/v1/users/:id", bytes.NewBufferString(tt.reqBody))
				c.Params = gin.Params{gin.Param{Key: "id", Value: "1"}}

				// Call the handler
				handler.UpdateUser(c)

				// Assert the response
				var actualBody map[string]any
				json.Unmarshal(w.Body.Bytes(), &actualBody)

				assert.Equal(t, http.StatusBadRequest, w.Code)
				assert.Equal(t, tt.expectedCode, actualBody["code"])
				assert.Equal(t, tt.expectedMsg, actualBody["message"])

				// Assert mock expectations
				userService.AssertExpectations(t)
				redisService.AssertExpectations(t)
				bcryptService.AssertExpectations(t)

			})
		}
	})

	t.Run("UpdateUser - Error Parse ID", func(t *testing.T) {
		userService := new(mocks.MockUserService)
		redisService := new(mocks.MockRedisService)
		bcryptService := new(mocks.MockBcryptService)
		handler := handlers.NewUserHandler(userService, redisService, bcryptService)

		requestBody := map[string]any{
			"name":     "Updated User",
			"birthday": "2000-01-01",
			"address":  "456 New Street",
		}
		body, _ := json.Marshal(requestBody)

		// Create a test context
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest("PATCH", "/api/v1/users/:id", bytes.NewBuffer(body))
		c.Params = gin.Params{gin.Param{Key: "id", Value: "invalid-id"}}

		// Call the handler
		handler.UpdateUser(c)

		// Assert the response
		var actualBody map[string]any
		json.Unmarshal(w.Body.Bytes(), &actualBody)
		expectedBody := map[string]any{
			"code":    float64(apperror.ErrParseError),
			"message": "Invalid UserID",
		}
		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Equal(t, expectedBody, actualBody)

		// Assert mocks
		userService.AssertExpectations(t)
		redisService.AssertExpectations(t)
		bcryptService.AssertExpectations(t)
	})

	t.Run("UpdateUser - User Not Found", func(t *testing.T) {
		userService := new(mocks.MockUserService)
		redisService := new(mocks.MockRedisService)
		bcryptService := new(mocks.MockBcryptService)
		handler := handlers.NewUserHandler(userService, redisService, bcryptService)

		var requestBody = map[string]any{
			"name":     "Updated User",
			"birthday": "2000-01-01",
			"address":  "456 New Street",
			"gender":   1,
		}
		body, _ := json.Marshal(requestBody)
		userService.On("GetUser", uint(1)).Return(&models.User{}, apperror.NewNotFoundError("User not found"))

		// Create a test context
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest("PATCH", "/api/v1/users/1", bytes.NewBuffer(body))
		c.Params = gin.Params{gin.Param{Key: "id", Value: "1"}}

		// Call the handler
		handler.UpdateUser(c)

		// Assert the response
		var actualBody map[string]any
		json.Unmarshal(w.Body.Bytes(), &actualBody)
		expectedBody := map[string]any{
			"code":    float64(apperror.ErrNotFound),
			"message": "User not found",
		}
		assert.Equal(t, http.StatusNotFound, w.Code)
		assert.Equal(t, expectedBody, actualBody)

		// Assert mocks
		userService.AssertExpectations(t)
		redisService.AssertExpectations(t)
		bcryptService.AssertExpectations(t)
	})

	t.Run("UpdateUser - Update User Error", func(t *testing.T) {
		userService := new(mocks.MockUserService)
		redisService := new(mocks.MockRedisService)
		bcryptService := new(mocks.MockBcryptService)
		handler := handlers.NewUserHandler(userService, redisService, bcryptService)

		user := &models.User{
			ID:    1,
			Email: "email@example.com",
			Name:  "User",
		}
		requestBody := map[string]any{
			"name":     "Updated User",
			"birthday": "2000-01-01",
			"address":  "456 New Street",
		}
		body, _ := json.Marshal(requestBody)
		// Mock methods
		userService.On("GetUser", uint(1)).Return(user, nil)
		userService.On("UpdateUser", user).Return(apperror.NewDBUpdateError("Update error"))

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest("PATCH", "/api/v1/users/1", bytes.NewBuffer(body))
		c.Params = gin.Params{gin.Param{Key: "id", Value: "1"}}

		// Call the handler
		handler.UpdateUser(c)

		// Assert the response
		var actualBody map[string]any
		json.Unmarshal(w.Body.Bytes(), &actualBody)
		expectedBody := map[string]any{
			"code":    float64(apperror.ErrDBUpdate),
			"message": "Update error",
		}
		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.Equal(t, expectedBody, actualBody)

		// Assert mocks
		userService.AssertExpectations(t)
		redisService.AssertExpectations(t)
		bcryptService.AssertExpectations(t)
	})

}

func TestDeleteUser(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("DelelteUser - Success", func(t *testing.T) {
		userService := new(mocks.MockUserService)
		redisService := new(mocks.MockRedisService)
		bcryptService := new(mocks.MockBcryptService)
		handler := handlers.NewUserHandler(userService, redisService, bcryptService)

		user := &models.User{
			ID:    1,
			Email: "email@example.com",
			Name:  "User",
		}
		// Mock the service methods
		userService.On("GetUser", uint(1)).Return(user, nil)
		userService.On("DeleteUser", uint(1)).Return(nil)

		// Create a test context
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest("DELETE", "/api/v1/users/:id", nil)
		c.Params = gin.Params{gin.Param{Key: "id", Value: "1"}}

		// Call the handler
		handler.DeleteUser(c)

		// Assert the response
		assert.Equal(t, http.StatusOK, w.Code)
		assert.JSONEq(t, `{"message":"Delete user successfully"}`, w.Body.String())

		// Assert mocks
		userService.AssertExpectations(t)
		redisService.AssertExpectations(t)
		bcryptService.AssertExpectations(t)
	})

	t.Run("DeleteUser - Failed To Parse UserID", func(t *testing.T) {
		userService := new(mocks.MockUserService)
		redisService := new(mocks.MockRedisService)
		bcryptService := new(mocks.MockBcryptService)
		handler := handlers.NewUserHandler(userService, redisService, bcryptService)

		// Create a test context
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest("DELETE", "/api/v1/users/:id", nil)
		c.Params = gin.Params{gin.Param{Key: "id", Value: "invalid-id"}}

		// Call the handler
		handler.DeleteUser(c)

		// Assert the response
		var actualBody map[string]any
		json.Unmarshal(w.Body.Bytes(), &actualBody)
		expectedBody := map[string]any{
			"code":    float64(apperror.ErrParseError),
			"message": "Invalid UserID",
		}
		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Equal(t, expectedBody, actualBody)

		// Assert mocks
		userService.AssertExpectations(t)
		redisService.AssertExpectations(t)
		bcryptService.AssertExpectations(t)
	})

	t.Run("DeleteUser - User Not Found", func(t *testing.T) {
		userService := new(mocks.MockUserService)
		redisService := new(mocks.MockRedisService)
		bcryptService := new(mocks.MockBcryptService)
		handler := handlers.NewUserHandler(userService, redisService, bcryptService)

		// Mock the service method
		userService.On("GetUser", uint(1)).Return(&models.User{}, apperror.NewNotFoundError("User not found"))

		// Create a test context
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest("DELETE", "/api/v1/users/:id", nil)
		c.Params = gin.Params{gin.Param{Key: "id", Value: "1"}}

		// Call the  handler
		handler.DeleteUser(c)

		// Assert the response
		var actualBody map[string]any
		json.Unmarshal(w.Body.Bytes(), &actualBody)
		expectedBody := map[string]any{
			"code":    float64(apperror.ErrNotFound),
			"message": "User not found",
		}
		assert.Equal(t, http.StatusNotFound, w.Code)
		assert.Equal(t, expectedBody, actualBody)

		// Assert mocks
		userService.AssertExpectations(t)
		redisService.AssertExpectations(t)
		bcryptService.AssertExpectations(t)
	})

	t.Run("DeleteUser - Failed To Delete", func(t *testing.T) {
		userService := new(mocks.MockUserService)
		redisService := new(mocks.MockRedisService)
		bcryptService := new(mocks.MockBcryptService)
		handler := handlers.NewUserHandler(userService, redisService, bcryptService)

		user := &models.User{
			ID:    1,
			Email: "email@example.com",
			Name:  "User",
		}
		// Mock the service methods
		userService.On("GetUser", uint(1)).Return(user, nil)
		userService.On("DeleteUser", uint(1)).Return(apperror.NewDBDeleteError("Delete error"))

		// Create a test context
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest("DELETE", "/api/v1/users/:id", nil)
		c.Params = gin.Params{gin.Param{Key: "id", Value: "1"}}

		// Call the handler
		handler.DeleteUser(c)

		// Assert the response
		var actualBody map[string]any
		json.Unmarshal(w.Body.Bytes(), &actualBody)
		expectedBody := map[string]any{
			"code":    float64(apperror.ErrDBDelete),
			"message": "Delete error",
		}
		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.Equal(t, expectedBody, actualBody)

		// Assert mocks
		userService.AssertExpectations(t)
		redisService.AssertExpectations(t)
		bcryptService.AssertExpectations(t)
	})
}

func TestResetPassword(t *testing.T) {
	gin.SetMode(gin.TestMode)

	utils.InitValidator()

	t.Run("ResetPassword - Success", func(t *testing.T) {
		userService := new(mocks.MockUserService)
		redisService := new(mocks.MockRedisService)
		bcryptService := new(mocks.MockBcryptService)
		handler := handlers.NewUserHandler(userService, redisService, bcryptService)

		expiredAt := time.Now().Add(24 * time.Hour).Unix()
		user := &models.User{
			ID:        1,
			Email:     "email@example.com",
			Name:      "User",
			Password:  "$2a$10$I/L5VegpCyOlJPoa1.KrmeCdezSBIandsEL5S2dd4Ap0YIWk0Iuka", // bcrypt hash of "12345678"
			ExpiredAt: &expiredAt,
		}

		requestBody := map[string]any{
			"token":        "token",
			"password":     "newpassword",
			"new_password": "newpassword",
		}
		body, _ := json.Marshal(requestBody)

		// Mock the service methods
		userService.On("GetUserByToken", "token").Return(user, nil)
		bcryptService.On("CheckPasswordHash", "newpassword", user.Password).Return(true)
		bcryptService.On("HashPassword", "newpassword").Return("$2a$10$hashedNewPassword", nil)
		userService.On("UpdateUser", user).Return(nil)

		// Create a test context
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest("POST", "/api/v1/reset-password", bytes.NewBuffer(body))

		// Call the  handler
		handler.ResetPassword(c)

		// Assert the response
		assert.Equal(t, http.StatusOK, w.Code)
		assert.JSONEq(t, `{"message":"Reset password successfully"}`, w.Body.String())

		// Assert mocks
		userService.AssertExpectations(t)
		redisService.AssertExpectations(t)
		bcryptService.AssertExpectations(t)
	})

	t.Run("ResetPassword - Not found user by token", func(t *testing.T) {
		userService := new(mocks.MockUserService)
		redisService := new(mocks.MockRedisService)
		bcryptService := new(mocks.MockBcryptService)
		handler := handlers.NewUserHandler(userService, redisService, bcryptService)

		requestBody := map[string]any{
			"token":        "invalid-token",
			"password":     "newpassword",
			"new_password": "newpassword",
		}
		body, _ := json.Marshal(requestBody)

		// Mock the service method's behavior
		userService.On("GetUserByToken", "invalid-token").
			Return(&models.User{}, apperror.NewNotFoundError("User not found"))

		// Create a test context
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest("POST", "/api/v1/reset-password", bytes.NewBuffer(body))

		// Call the handler
		handler.ResetPassword(c)

		// Assert the response
		expectedBody := map[string]any{
			"code":    float64(apperror.ErrNotFound),
			"message": "User not found",
		}
		var actualBody map[string]any
		json.Unmarshal(w.Body.Bytes(), &actualBody)

		assert.Equal(t, http.StatusNotFound, w.Code)
		assert.Equal(t, expectedBody["code"], actualBody["code"])
		assert.Equal(t, expectedBody["message"], actualBody["message"])

		// Assert mocks
		userService.AssertExpectations(t)
		redisService.AssertExpectations(t)
		bcryptService.AssertExpectations(t)
	})

	t.Run("ResetPassword - Token Expired", func(t *testing.T) {
		userService := new(mocks.MockUserService)
		redisService := new(mocks.MockRedisService)
		bcryptService := new(mocks.MockBcryptService)
		handler := handlers.NewUserHandler(userService, redisService, bcryptService)

		expiredAt := time.Now().Add(-24 * time.Hour).Unix() // Token expired
		user := &models.User{
			ID:        1,
			Email:     "email@example.com",
			Name:      "User",
			Password:  "$2a$10$I/L5VegpCyOlJPoa1.KrmeCdezSBIandsEL5S2dd4Ap0YIWk0Iuka", // bcrypt hash of "12345678"
			ExpiredAt: &expiredAt,
		}

		requestBody := map[string]any{
			"token":        "invalid-token",
			"password":     "newpassword",
			"new_password": "newpassword",
		}
		body, _ := json.Marshal(requestBody)

		// Mock the service method
		userService.On("GetUserByToken", "invalid-token").Return(user, nil)

		// Create a test context
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest("POST", "/api/v1/reset-password", bytes.NewBuffer(body))

		// Call the handler
		handler.ResetPassword(c)

		// Assert the response
		expectedBody := map[string]any{
			"code":    float64(apperror.ErrTokenExpired),
			"message": "Token is expired",
		}
		var actualBody map[string]any
		json.Unmarshal(w.Body.Bytes(), &actualBody)
		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Equal(t, expectedBody["code"], actualBody["code"])
		assert.Equal(t, expectedBody["message"], actualBody["message"])

		// Assert mocks
		userService.AssertExpectations(t)
		redisService.AssertExpectations(t)
		bcryptService.AssertExpectations(t)

	})

	t.Run("ResetPassword - Passwords Incorrect", func(t *testing.T) {
		userService := new(mocks.MockUserService)
		redisService := new(mocks.MockRedisService)
		bcryptService := new(mocks.MockBcryptService)
		handler := handlers.NewUserHandler(userService, redisService, bcryptService)

		expiredAt := time.Now().Add(24 * time.Hour).Unix()
		user := &models.User{
			ID:        1,
			Email:     "email@example.com",
			Name:      "User",
			Password:  "$2a$10$I/L5VegpCyOlJPoa1.KrmeCdezSBIandsEL5S2dd4Ap0YIWk0Iuka", // bcrypt hash of "12345678"
			ExpiredAt: &expiredAt,
		}

		requestBody := map[string]any{
			"token":        "token",
			"password":     "newpassword",
			"new_password": "newpassword",
		}
		body, _ := json.Marshal(requestBody)

		// Mock the service methods
		userService.On("GetUserByToken", "token").Return(user, nil)
		bcryptService.On("CheckPasswordHash", "newpassword", user.Password).Return(false)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest("POST", "/api/v1/reset-password", bytes.NewBuffer(body))

		// Call the handler
		handler.ResetPassword(c)

		// Assert the response
		expectedBody := map[string]any{
			"code":    float64(apperror.ErrInvalidPassword),
			"message": "Old password is incorrect",
		}
		var actualBody map[string]any
		json.Unmarshal(w.Body.Bytes(), &actualBody)
		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Equal(t, expectedBody["code"], actualBody["code"])
		assert.Equal(t, expectedBody["message"], actualBody["message"])

		// Assert mocks
		userService.AssertExpectations(t)
		redisService.AssertExpectations(t)
		bcryptService.AssertExpectations(t)
	})

	t.Run("ResetPassword - Error Hashing Password Failed", func(t *testing.T) {
		userService := new(mocks.MockUserService)
		redisService := new(mocks.MockRedisService)
		bcryptService := new(mocks.MockBcryptService)
		handler := handlers.NewUserHandler(userService, redisService, bcryptService)

		expiredAt := time.Now().Add(24 * time.Hour).Unix()
		user := &models.User{
			ID:        1,
			Email:     "email@example.com",
			Name:      "User",
			Password:  "$2a$10$I/L5VegpCyOlJPoa1.KrmeCdezSBIandsEL5S2dd4Ap0YIWk0Iuka", // bcrypt hash of "12345678"
			ExpiredAt: &expiredAt,
		}

		requestBody := map[string]any{
			"token":        "token",
			"password":     "newpassword",
			"new_password": "newpassword",
		}
		body, _ := json.Marshal(requestBody)

		// Mock the service methods
		userService.On("GetUserByToken", "token").Return(user, nil)
		bcryptService.On("CheckPasswordHash", "newpassword", user.Password).Return(true)
		bcryptService.On("HashPassword", "newpassword").
			Return("", apperror.NewInternalError("Failed to hash password"))

		// Create a test context
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest("POST", "/api/v1/reset-password", bytes.NewBuffer(body))

		// Call the handler
		handler.ResetPassword(c)

		// Assert the response
		expectedBody := map[string]any{
			"code":    float64(apperror.ErrPasswordHashFailed),
			"message": "Failed to hash password",
		}
		var actualBody map[string]any
		json.Unmarshal(w.Body.Bytes(), &actualBody)
		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.Equal(t, expectedBody["code"], actualBody["code"])
		assert.Equal(t, expectedBody["message"], actualBody["message"])

		// Assert mocks
		userService.AssertExpectations(t)
		redisService.AssertExpectations(t)
		bcryptService.AssertExpectations(t)
	})

	t.Run("Error failed to UpdateUser", func(t *testing.T) {
		userService := new(mocks.MockUserService)
		redisService := new(mocks.MockRedisService)
		bcryptService := new(mocks.MockBcryptService)
		handler := handlers.NewUserHandler(userService, redisService, bcryptService)

		expiredAt := time.Now().Add(24 * time.Hour).Unix()
		user := &models.User{
			ID:        1,
			Email:     "email@example.com",
			Name:      "User",
			Password:  "$2a$10$I/L5VegpCyOlJPoa1.KrmeCdezSBIandsEL5S2dd4Ap0YIWk0Iuka", // bcrypt hash of "12345678"
			ExpiredAt: &expiredAt,
		}

		requestBody := map[string]any{
			"token":        "token",
			"password":     "newpassword",
			"new_password": "newpassword",
		}
		body, _ := json.Marshal(requestBody)

		// Mock the service methods
		userService.On("GetUserByToken", "token").Return(user, nil)
		bcryptService.On("CheckPasswordHash", "newpassword", user.Password).Return(true)
		bcryptService.On("HashPassword", "newpassword").Return("$2a$10$hashedNewPassword", nil)
		userService.On("UpdateUser", user).Return(apperror.NewDBUpdateError("Failed to update user"))

		// Create a test context
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest("POST", "/api/v1/reset-password", bytes.NewBuffer(body))
		c.Request.Header.Set("Content-Type", "application/json")
		c.Request.Header.Set("Authorization", "Bearer test-token")
		c.Set("UserID", uint(1))

		// Call the handler
		handler.ResetPassword(c)

		// Assert the response
		expectedBody := map[string]any{
			"code":    float64(apperror.ErrDBUpdate),
			"message": "Failed to update user",
		}
		var actualBody map[string]any
		json.Unmarshal(w.Body.Bytes(), &actualBody)
		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.Equal(t, expectedBody["code"], actualBody["code"])
		assert.Equal(t, expectedBody["message"], actualBody["message"])

		// Assert mocks
		userService.AssertExpectations(t)
		redisService.AssertExpectations(t)
		bcryptService.AssertExpectations(t)
	})

	t.Run("Validation Error", func(t *testing.T) {
		tests := []struct {
			name          string
			reqBody       string
			expectedCode  float64
			expectedMsg   string
			expectedField []apperror.FieldError
		}{
			{
				name:         "EmptyToken",
				reqBody:      `{"token":""}`,
				expectedCode: 4001,
				expectedMsg:  "Validation failed",
				expectedField: []apperror.FieldError{
					{
						Field:   "token",
						Message: "token is required",
					},
					{
						Field:   "password",
						Message: "password is required",
					},
					{
						Field:   "new_password",
						Message: "new_password is required",
					},
				},
			},
			{
				name:         "EmptyPassword",
				reqBody:      `{"token":"valid-token","password":""}`,
				expectedCode: 4001,
				expectedMsg:  "Validation failed",
				expectedField: []apperror.FieldError{
					{
						Field:   "password",
						Message: "password is required",
					},
					{
						Field:   "new_password",
						Message: "new_password is required",
					},
				},
			},
			{
				name:         "PasswordTooShort",
				reqBody:      `{"token":"valid-token","password":"short"}`,
				expectedCode: 4001,
				expectedMsg:  "Validation failed",
				expectedField: []apperror.FieldError{
					{
						Field:   "password",
						Message: "password must be at least 6 characters long or numeric",
					},
					{
						Field:   "new_password",
						Message: "new_password is required",
					},
				},
			},
			{
				name:         "PasswordTooLong",
				reqBody:      `{"token":"valid-token","password":"` + strings.Repeat("a", 256) + `"}`,
				expectedCode: 4001,
				expectedMsg:  "Validation failed",
				expectedField: []apperror.FieldError{
					{
						Field:   "password",
						Message: "password must be at most 255 characters long or numeric",
					},
					{
						Field:   "new_password",
						Message: "new_password is required",
					},
				},
			},
			{
				name:         "EmptyNewPassword",
				reqBody:      `{"token":"valid-token","password":"newpassword","new_password":""}`,
				expectedCode: 4001,
				expectedMsg:  "Validation failed",
				expectedField: []apperror.FieldError{
					{
						Field:   "new_password",
						Message: "new_password is required",
					},
				},
			},
			{
				name:         "NewPasswordTooShort",
				reqBody:      `{"token":"valid-token","password":"newpassword","new_password":"short"}`,
				expectedCode: 4001,
				expectedMsg:  "Validation failed",
				expectedField: []apperror.FieldError{
					{
						Field:   "new_password",
						Message: "new_password must be at least 6 characters long or numeric",
					},
				},
			},
			{
				name:         "NewPasswordTooLong",
				reqBody:      `{"token":"valid-token","password":"newpassword","new_password":"` + strings.Repeat("a", 256) + `"}`,
				expectedCode: 4001,
				expectedMsg:  "Validation failed",
				expectedField: []apperror.FieldError{
					{
						Field:   "new_password",
						Message: "new_password must be at most 255 characters long or numeric",
					},
				},
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				userService := new(mocks.MockUserService)
				redisService := new(mocks.MockRedisService)
				bcryptService := new(mocks.MockBcryptService)
				handler := handlers.NewUserHandler(userService, redisService, bcryptService)

				// Create a test context
				w := httptest.NewRecorder()
				c, _ := gin.CreateTestContext(w)
				c.Request, _ = http.NewRequest("POST", "/api/v1/reset-password", bytes.NewBufferString(tt.reqBody))
				c.Set("UserID", uint(1))

				// Call the handler
				handler.ResetPassword(c)

				// Assert the response
				var actualBody map[string]any
				json.Unmarshal(w.Body.Bytes(), &actualBody)
				assert.Equal(t, http.StatusBadRequest, w.Code)
				assert.Equal(t, tt.expectedCode, actualBody["code"])
				assert.Equal(t, tt.expectedMsg, actualBody["message"])
				assert.Equal(t, tt.expectedField, utils.MapJsonToFieldErrors(actualBody["fields"]))

				// Assert mocks
				userService.AssertExpectations(t)
				redisService.AssertExpectations(t)
				bcryptService.AssertExpectations(t)
			})
		}
	})

}
