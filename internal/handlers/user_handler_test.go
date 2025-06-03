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
	appError "github.com/vfa-khuongdv/golang-cms/pkg/errors"
	"github.com/vfa-khuongdv/golang-cms/tests/mocks"
)

func TestPaginateUser(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("Success", func(t *testing.T) {
		userService := new(mocks.MockUserService)
		redisService := new(mocks.MockRedisService)
		bcryptService := new(mocks.MockBcryptService)

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
		userService.On("PaginateUser", 1, 10).Return(expected, nil)

		handler := handlers.NewUserHandler(userService, redisService, bcryptService)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		c.Request, _ = http.NewRequest("GET", "/api/v1/users?page=1&limit=10", nil)

		handler.PaginationUser(c)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.JSONEq(t, `{"page":1,"limit":10,"totalItems":2,"totalPages":1,"data":[{"id":1,"email":"email1@example.com","name":"User One","gender":1,"createdAt":"2023-10-01T00:00:00Z","updatedAt":"2023-10-01T00:00:00Z","deletedAt":null},{"id":2,"email":"email2@example.com","name":"User Two","gender":2,"createdAt":"2023-10-02T00:00:00Z","updatedAt":"2023-10-02T00:00:00Z","deletedAt":null}]}`, w.Body.String())
	})

	t.Run("Error", func(t *testing.T) {
		userService := new(mocks.MockUserService)
		redisService := new(mocks.MockRedisService)
		bcryptService := new(mocks.MockBcryptService)
		userService.On("PaginateUser", 1, 10).Return(&utils.Pagination{}, appError.New(appError.ErrDBQuery, "Database query error"))

		handler := handlers.NewUserHandler(userService, redisService, bcryptService)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		c.Request, _ = http.NewRequest("GET", "/api/v1/users?page=1&limit=10", nil)

		handler.PaginationUser(c)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.JSONEq(t, fmt.Sprintf(`{"code":%d,"message":"Database query error"}`, appError.ErrDBQuery), w.Body.String())
	})
}

func TestCreateUser(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Initialize the validator
	utils.InitValidator()

	t.Run("Success", func(t *testing.T) {
		userService := new(mocks.MockUserService)
		redisService := new(mocks.MockRedisService)
		bcryptService := new(mocks.MockBcryptService)

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

		handler := handlers.NewUserHandler(userService, redisService, bcryptService)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		c.Request, _ = http.NewRequest("POST", "/api/v1/users", bytes.NewBuffer(body))
		c.Request.Header.Set("Content-Type", "application/json")

		handler.CreateUser(c)

		assert.Equal(t, http.StatusCreated, w.Code)
		assert.JSONEq(t, `{"message":"Create user successfully"}`, w.Body.String())
	})

	t.Run("Validation Error", func(t *testing.T) {
		userService := new(mocks.MockUserService)
		redisService := new(mocks.MockRedisService)
		bcryptService := new(mocks.MockBcryptService)

		tests := []struct {
			name         string
			reqBody      string
			expectedCode float64
			expectedMsg  string
		}{
			{
				name:         "MissingEmail",
				reqBody:      `{}`,
				expectedCode: float64(4001),
				expectedMsg:  "email is required",
			},
			{
				name:         "EmptyEmail",
				reqBody:      `{"email":""}`,
				expectedCode: float64(4001),
				expectedMsg:  "email is required",
			},
			{
				name:         "InvalidEmailFormat",
				reqBody:      `{"email":"not-an-email"}`,
				expectedCode: float64(4001),
				expectedMsg:  "email must be a valid email address",
			},
			{
				name:         "MissingPassword",
				reqBody:      `{"email":"email@example.com"}`,
				expectedCode: float64(4001),
				expectedMsg:  "password is required",
			},
			{
				name:         "EmptyPassword",
				reqBody:      `{"email":"email@example.com","password":""}`,
				expectedCode: float64(4001),
				expectedMsg:  "password is required",
			},
			{
				name:         "ShortPassword",
				reqBody:      `{"email":"email@example.com","password":"123"}`,
				expectedCode: float64(4001),
				expectedMsg:  "password must be at least 6 characters long or numeric",
			},
			{
				name:         "LongPassword",
				reqBody:      `{"email":"email@example.com","password":"` + strings.Repeat("a", 256) + `"}`,
				expectedCode: float64(4001),
				expectedMsg:  "password must be at most 255 characters long or numeric",
			},
			{
				name:         "MissingName",
				reqBody:      `{"email":"email@example.com","password":"password"}`,
				expectedCode: float64(4001),
				expectedMsg:  "name is required",
			},
			{
				name:         "NameNotBlank",
				reqBody:      `{"email":"email@example.com","password":"password","name":"  "}`,
				expectedCode: float64(4001),
				expectedMsg:  "name must not be blank",
			},
			{
				name:         "EmptyName",
				reqBody:      `{"email":"email@example.com","password":"password","name":""}`,
				expectedCode: float64(4001),
				expectedMsg:  "name is required",
			},
			{
				name:         "LongName",
				reqBody:      `{"email":"email@example.com","password":"password","name": "` + strings.Repeat("a", 46) + `","birthday":"2000-01-01","address":"address","gender":1}`,
				expectedCode: float64(4001),
				expectedMsg:  "name must be at most 45 characters long or numeric",
			},
			{
				name:         "MissingBirthday",
				reqBody:      `{"email":"email@example.com","password":"password","name": "Bob"}`,
				expectedCode: float64(4001),
				expectedMsg:  "birthday is required",
			},
			{
				name:         "InvalidBirthdayFormat",
				reqBody:      `{"email":"email@example.com","password":"password","name": "Bob","birthday":"invalid-date"}`,
				expectedCode: float64(4001),
				expectedMsg:  "birthday must be a valid date (YYYY-MM-DD) and not in the future",
			},
			{
				name:         "FutureBirthday",
				reqBody:      `{"email":"email@example.com","password":"password","name": "Bob","birthday":"3000-01-01"}`,
				expectedCode: float64(4001),
				expectedMsg:  "birthday must be a valid date (YYYY-MM-DD) and not in the future",
			},
			{
				name:         "MissingAddress",
				reqBody:      `{"email":"email@example.com","password":"password","name": "Bob","birthday":"2000-01-01"}`,
				expectedCode: float64(4001),
				expectedMsg:  "address is required",
			},
			{
				name:         "AddressNotBlank",
				reqBody:      `{"email":"email@example.com","password":"password","name": "Bob","birthday":"2000-01-01","address":"  "}`,
				expectedCode: float64(4001),
				expectedMsg:  "address must not be blank",
			},
			{
				name:         "EmptyAddress",
				reqBody:      `{"email":"email@example.com","password":"password","name": "Bob","birthday":"2000-01-01","address":""}`,
				expectedCode: float64(4001),
				expectedMsg:  "address must be at least 1 characters long or numeric",
			},
			{
				name:         "LongAddress",
				reqBody:      `{"email":"email@example.com","password":"password","name": "Bob","birthday":"2000-01-01","address":"` + strings.Repeat("a", 256) + `"}`,
				expectedCode: float64(4001),
				expectedMsg:  "address must be at most 255 characters long or numeric",
			},
			{
				name:         "MissingGender",
				reqBody:      `{"email":"email@example.com","password":"password","name": "Bob","birthday":"2000-01-01","address":"address"}`,
				expectedCode: float64(4001),
				expectedMsg:  "gender is required",
			},
			{
				name:         "InvalidGender",
				reqBody:      `{"email":"email@example.com","password":"password","name": "Bob","birthday":"2000-01-01","address":"address", "gender": 4}`,
				expectedCode: float64(4001),
				expectedMsg:  "gender must be one of [1 2 3]",
			},
			{
				name:         "StringGender",
				reqBody:      `{"email":"email@example.com","password":"password","name": "Bob","birthday":"2000-01-01","address":"address", "gender": "not_numeric"}`,
				expectedCode: float64(4001),
				expectedMsg:  "json: cannot unmarshal string into Go struct field .gender of type int16",
			},
			{
				name:         "MissingRoleIDs",
				reqBody:      `{"email":"email@example.com","password":"password","name": "Bob","birthday":"2000-01-01","address":"address","gender":1}`,
				expectedCode: float64(4001),
				expectedMsg:  "roleids is required",
			},
			{
				name:         "EmptyRoleIDs",
				reqBody:      `{"email":"email@example.com","password":"password","name": "Bob","birthday":"2000-01-01","address":"address","gender":1,"role_ids":""}`,
				expectedCode: float64(4001),
				expectedMsg:  "json: cannot unmarshal string into Go struct field .role_ids of type []uint",
			},
			{
				name:         "InvalidRoleIDsIdNotNumeric",
				reqBody:      `{"email":"email@example.com","password":"password","name": "Bob","birthday":"2000-01-01","address":"address","gender":1,"role_ids":["not_numeric"]}`,
				expectedCode: float64(4001),
				expectedMsg:  "json: cannot unmarshal string into Go struct field .role_ids of type uint",
			},
		}

		for _, tc := range tests {
			t.Run(tc.name, func(t *testing.T) {
				w := httptest.NewRecorder()
				c, _ := gin.CreateTestContext(w)
				c.Request, _ = http.NewRequest("POST", "/api/v1/user", bytes.NewBufferString(tc.reqBody))
				c.Request.Header.Set("Content-Type", "application/json")

				handler := handlers.NewUserHandler(userService, redisService, bcryptService)

				handler.CreateUser(c)

				assert.Equal(t, http.StatusBadRequest, w.Code)

				expectedBody := map[string]any{
					"code":    tc.expectedCode,
					"message": tc.expectedMsg,
				}

				var actualBody map[string]any
				err := json.Unmarshal(w.Body.Bytes(), &actualBody)
				require.NoError(t, err)

				assert.Equal(t, expectedBody["code"], actualBody["code"])
				assert.Equal(t, expectedBody["message"], actualBody["message"])
			})
		}
	})

	t.Run("Create user Error", func(t *testing.T) {
		userService := new(mocks.MockUserService)
		redisService := new(mocks.MockRedisService)
		bcryptService := new(mocks.MockBcryptService)

		// Mock the HashPassword method
		bcryptService.On("HashPassword", "password").Return("$2a$10$examplehash", nil)
		// Mock the CreateUser method to return an error
		userService.On("CreateUser", mock.AnythingOfType("*models.User"), mock.AnythingOfType("[]uint")).Return(appError.New(appError.ErrDBInsert, "Database insert error"))

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

		handler := handlers.NewUserHandler(userService, redisService, bcryptService)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest("POST", "/api/v1/user", bytes.NewBuffer(body))
		c.Request.Header.Set("Content-Type", "application/json")

		// Call the CreateUser handler
		handler.CreateUser(c)

		// Assert the response
		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.JSONEq(t, fmt.Sprintf(`{"code":%d,"message":"Database insert error"}`, appError.ErrDBInsert), w.Body.String())
	})

	t.Run("Error Bcrypt Hash Password", func(t *testing.T) {
		userService := new(mocks.MockUserService)
		redisService := new(mocks.MockRedisService)
		bcryptService := new(mocks.MockBcryptService)

		// Mock the HashPassword method to return an error
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
		handler := handlers.NewUserHandler(userService, redisService, bcryptService)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest("POST", "/api/v1/user", bytes.NewBuffer(body))
		c.Request.Header.Set("Content-Type", "application/json")
		c.Request.Header.Set("Authorization", "Bearer test-token")
		c.Set("UserID", uint(1))
		// Call the CreateUser handler
		handler.CreateUser(c)
		// Assert the response
		assert.Equal(t, http.StatusBadRequest, w.Code)

		expectedBody := map[string]any{
			"code":    float64(appError.ErrPasswordHashFailed),
			"message": "Hash password failed",
		}

		var actualBody map[string]any
		err := json.Unmarshal(w.Body.Bytes(), &actualBody)
		require.NoError(t, err)
		assert.Equal(t, expectedBody["code"], actualBody["code"])
		assert.Equal(t, expectedBody["message"], actualBody["message"])

		userService.AssertNotCalled(t, "CreateUser", mock.AnythingOfType("*models.User"), mock.AnythingOfType("[]uint"))

	})
}

func TestUpdateProfile(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Initialize the validator
	utils.InitValidator()

	t.Run("Success", func(t *testing.T) {
		userService := new(mocks.MockUserService)
		redisService := new(mocks.MockRedisService)
		bcryptService := new(mocks.MockBcryptService)

		user := &models.User{
			ID:    1,
			Email: "email@example.com",
			Name:  "User",
		}

		// Assuming the cache key is constructed as "profile:<user_id>"
		profileKey := constants.PROFILE + string(rune(user.ID))

		// Mock the GetUser method
		userService.On("GetUser", uint(1)).Return(user, nil)
		// Mock the UpdateUser method
		userService.On("UpdateUser", user).Return(nil)
		redisService.On("Delete", profileKey).Return(nil)

		requestBody := map[string]any{
			"name":     "Updated User",
			"birthday": "2000-01-01",
			"address":  "456 New Street",
			"gender":   2,
		}

		body, _ := json.Marshal(requestBody)

		handler := handlers.NewUserHandler(userService, redisService, bcryptService)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest("PUT", "/api/v1/profile", bytes.NewBuffer(body))
		c.Request.Header.Set("Content-Type", "application/json")
		c.Request.Header.Set("Authorization", "Bearer test-token")
		c.Set("UserID", uint(1))

		// Call the UpdateProfile handler
		handler.UpdateProfile(c)

		// Assert the response
		assert.Equal(t, http.StatusOK, w.Code)
		assert.JSONEq(t, `{"message":"Update profile successfully"}`, w.Body.String())

	})

	t.Run("Validation Error", func(t *testing.T) {
		tests := []struct {
			name         string
			reqBody      string
			expectedCode float64
			expectedMsg  string
		}{
			{
				name:         "EmptyName",
				reqBody:      `{"name":""}`,
				expectedCode: float64(4001),
				expectedMsg:  "name must be at least 1 characters long or numeric",
			},
			{
				name:         "NameNotBlank",
				reqBody:      `{"name":"  "}`,
				expectedCode: float64(4001),
				expectedMsg:  "name must not be blank",
			},
			{
				name:         "LongName",
				reqBody:      `{"name": "` + strings.Repeat("a", 46) + `"}`,
				expectedCode: float64(4001),
				expectedMsg:  "name must be at most 45 characters long or numeric",
			},
			{
				name:         "InvalidBirthdayFormat",
				reqBody:      `{"name": "User", "birthday": "invalid-date"}`,
				expectedCode: float64(4001),
				expectedMsg:  "birthday must be a valid date (YYYY-MM-DD) and not in the future",
			},
			{
				name:         "FutureBirthday",
				reqBody:      `{"name": "User", "birthday": "3000-01-01"}`,
				expectedCode: float64(4001),
				expectedMsg:  "birthday must be a valid date (YYYY-MM-DD) and not in the future",
			},
			{
				name:         "EmptyAddress",
				reqBody:      `{"name": "User", "birthday": "2000-01-01", "address": ""}`,
				expectedCode: float64(4001),
				expectedMsg:  "address must be at least 1 characters long or numeric",
			},
			{
				name:         "LongAddress",
				reqBody:      `{"name": "User", "birthday": "2000-01-01", "address": "` + strings.Repeat("a", 256) + `"}`,
				expectedCode: float64(4001),
				expectedMsg:  "address must be at most 255 characters long or numeric",
			},
			{
				name:         "AddressNotBlank",
				reqBody:      `{"name": "User", "birthday": "2000-01-01", "address": "  "}`,
				expectedCode: float64(4001),
				expectedMsg:  "address must not be blank",
			},
			{
				name:         "InvalidGender 0",
				reqBody:      `{"name": "User", "birthday": "2000-01-01", "address": "123 Street", "gender": 0}`,
				expectedCode: float64(4001),
				expectedMsg:  "gender must be one of [1 2 3]",
			},
			{
				name:         "InvalidGender 4",
				reqBody:      `{"name": "User", "birthday": "2000-01-01", "address": "123 Street", "gender": 4}`,
				expectedCode: float64(4001),
				expectedMsg:  "gender must be one of [1 2 3]",
			},
			{
				name:         "StringGender",
				reqBody:      `{"name": "User", "birthday": "2000-01-01", "address": "123 Street", "gender": "male"}`,
				expectedCode: float64(4001),
				expectedMsg:  "json: cannot unmarshal string into Go struct field .gender of type int16",
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				userService := new(mocks.MockUserService)
				redisService := new(mocks.MockRedisService)
				bcryptService := new(mocks.MockBcryptService)

				handler := handlers.NewUserHandler(userService, redisService, bcryptService)
				w := httptest.NewRecorder()
				c, _ := gin.CreateTestContext(w)
				c.Request, _ = http.NewRequest("PUT", "/api/v1/profile", bytes.NewBufferString(tt.reqBody))
				c.Request.Header.Set("Content-Type", "application/json")
				c.Request.Header.Set("Authorization", "Bearer test-token")
				c.Set("UserID", uint(1))

				// Call the UpdateProfile handler
				handler.UpdateProfile(c)

				// Assert the response
				assert.Equal(t, http.StatusBadRequest, w.Code)

				expectedBody := map[string]any{
					"code":    tt.expectedCode,
					"message": tt.expectedMsg,
				}

				var actualBody map[string]any
				err := json.Unmarshal(w.Body.Bytes(), &actualBody)
				require.NoError(t, err)

				assert.Equal(t, expectedBody["code"], actualBody["code"])
				assert.Equal(t, expectedBody["message"], actualBody["message"])
			})
		}
	})

	t.Run("Error User Not found from ctx", func(t *testing.T) {
		userService := new(mocks.MockUserService)
		redisService := new(mocks.MockRedisService)
		bcryptService := new(mocks.MockBcryptService)

		handler := handlers.NewUserHandler(userService, redisService, bcryptService)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest("PUT", "/api/v1/profile", nil)
		c.Request.Header.Set("Content-Type", "application/json")
		c.Request.Header.Set("Authorization", "Bearer test-token")
		c.Set("UserID", uint(0)) // Invalid User ID

		// Call the UpdateProfile handler
		handler.UpdateProfile(c)

		// Assert the response
		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.JSONEq(t, fmt.Sprintf(`{"code":%d,"message":"Invalid UserID"}`, appError.ErrParseError), w.Body.String())
	})

	t.Run("Error User Not Found", func(t *testing.T) {
		userService := new(mocks.MockUserService)
		redisService := new(mocks.MockRedisService)
		bcryptService := new(mocks.MockBcryptService)

		requestBody := map[string]any{
			"name":     "Updated User",
			"birthday": "2000-01-01",
			"address":  "456 New Street",
			"gender":   2,
		}

		body, _ := json.Marshal(requestBody)

		// Mock the GetUser method to return an error
		userService.On("GetUser", uint(1)).Return(&models.User{}, appError.New(appError.ErrDBQuery, "Query error"))

		handler := handlers.NewUserHandler(userService, redisService, bcryptService)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest("PUT", "/api/v1/profile", bytes.NewBuffer(body))
		c.Request.Header.Set("Content-Type", "application/json")
		c.Request.Header.Set("Authorization", "Bearer test-token")
		c.Set("UserID", uint(1))

		// Call the UpdateProfile handler
		handler.UpdateProfile(c)

		// Assert the response
		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.JSONEq(t, fmt.Sprintf(`{"code":%d,"message":"Query error"}`, appError.ErrDBQuery), w.Body.String())
	})

	t.Run("Error Update User", func(t *testing.T) {
		userService := new(mocks.MockUserService)
		redisService := new(mocks.MockRedisService)
		bcryptService := new(mocks.MockBcryptService)

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

		// Mock the GetUser method
		userService.On("GetUser", uint(1)).Return(user, nil)
		// Mock the UpdateUser method to return an error
		userService.On("UpdateUser", user).Return(appError.New(appError.ErrDBUpdate, "Update error"))
		// Mock the Redis Delete method
		profileKey := constants.PROFILE + string(rune(user.ID))
		redisService.On("Delete", profileKey).Return(nil)
		handler := handlers.NewUserHandler(userService, redisService, bcryptService)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest("PUT", "/api/v1/profile", bytes.NewBuffer(body))
		c.Request.Header.Set("Content-Type", "application/json")
		c.Request.Header.Set("Authorization", "Bearer test-token")
		c.Set("UserID", uint(1))

		// Call the UpdateProfile handler
		handler.UpdateProfile(c)

		// Assert the response
		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.JSONEq(t, fmt.Sprintf(`{"code":%d,"message":"Update error"}`, appError.ErrDBUpdate), w.Body.String())
	})

	t.Run("Error Delete Cache", func(t *testing.T) {
		userService := new(mocks.MockUserService)
		redisService := new(mocks.MockRedisService)
		bcryptService := new(mocks.MockBcryptService)

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
		// Mock the GetUser method
		userService.On("GetUser", uint(1)).Return(user, nil)
		// Mock the UpdateUser method
		userService.On("UpdateUser", user).Return(nil)
		// Mock the Redis Delete method to return an error
		profileKey := constants.PROFILE + string(rune(user.ID))
		redisService.On("Delete", profileKey).Return(errors.New("Redis delete error"))
		handler := handlers.NewUserHandler(userService, redisService, bcryptService)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest("PUT", "/api/v1/profile", bytes.NewBuffer(body))
		c.Request.Header.Set("Content-Type", "application/json")
		c.Request.Header.Set("Authorization", "Bearer test-token")
		c.Set("UserID", uint(1))

		// Call the UpdateProfile handler
		handler.UpdateProfile(c)

		// Assert the response
		assert.Equal(t, http.StatusOK, w.Code)
		assert.JSONEq(t, fmt.Sprintf(`{"message":"%s"}`, "Update profile successfully"), w.Body.String())
	})
}

func TestGetProfile(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("Success get profile from database", func(t *testing.T) {
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
		// Mock the GetUser method
		userService.On("GetProfile", uint(1)).Return(user, nil)
		// Mock the Redis Get method
		redisService.On("Get", profileKey).Return("", nil)

		// Parse the user into a JSON string
		profileData, _ := json.Marshal(user)
		// Set the TTL for the cache
		ttl := 60 * time.Minute
		// Mock the Redis Set method to cache the profile
		redisService.On("Set", profileKey, profileData, ttl).Return(nil)

		handler := handlers.NewUserHandler(userService, redisService, bcryptService)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest("GET", "/api/v1/profile", nil)
		c.Request.Header.Set("Authorization", "Bearer test-token")
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

		assert.Equal(t, expectedBody, actualBody)
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
		c.Request.Header.Set("Authorization", "Bearer test-token")
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

		assert.Equal(t, expectedBody, actualBody)
	})

	t.Run("Error Invalid User ID", func(t *testing.T) {
		userService := new(mocks.MockUserService)
		redisService := new(mocks.MockRedisService)
		bcryptService := new(mocks.MockBcryptService)

		handler := handlers.NewUserHandler(userService, redisService, bcryptService)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest("GET", "/api/v1/profile", nil)
		c.Request.Header.Set("Authorization", "Bearer test-token")
		c.Set("UserID", uint(0)) // Invalid User ID

		// Call the GetProfile handler
		handler.GetProfile(c)

		// Assert the response
		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.JSONEq(t, fmt.Sprintf(`{"code":%d,"message":"Invalid UserID"}`, appError.ErrParseError), w.Body.String())
	})

	t.Run("Error User Not Found", func(t *testing.T) {
		userService := new(mocks.MockUserService)
		redisService := new(mocks.MockRedisService)
		bcryptService := new(mocks.MockBcryptService)

		profileKey := constants.PROFILE + string(rune(1))
		// Mock the GetUser method to return an error
		userService.On("GetProfile", uint(1)).Return(&models.User{}, appError.New(appError.ErrNotFound, "User not found"))
		// Mock the Redis Get method to return an empty string
		redisService.On("Get", profileKey).Return("", nil)

		handler := handlers.NewUserHandler(userService, redisService, bcryptService)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest("GET", "/api/v1/profile", nil)
		c.Request.Header.Set("Authorization", "Bearer test-token")
		c.Set("UserID", uint(1))

		// Call the GetProfile handler
		handler.GetProfile(c)

		// Assert the response
		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.JSONEq(t, fmt.Sprintf(`{"code":%d,"message":"User not found"}`, appError.ErrNotFound), w.Body.String())
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
		c.Request.Header.Set("Authorization", "Bearer test-token")
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

	t.Run("Success Get Profile from Redis but Parse Error", func(t *testing.T) {
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
		c.Request.Header.Set("Authorization", "Bearer test-token")
		c.Set("UserID", uint(1))

		// Call the GetProfile handler
		handler.GetProfile(c)

		// Assert the response
		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.JSONEq(t, fmt.Sprintf(`{"code":%d,"message":"Failed to parse user data from cache"}`, appError.ErrParseError), w.Body.String())
	})
}

func TestGetUser(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("Success", func(t *testing.T) {
		userService := new(mocks.MockUserService)
		redisService := new(mocks.MockRedisService)
		bcryptService := new(mocks.MockBcryptService)

		router := gin.Default()

		user := &models.User{
			ID:        1,
			Email:     "email@example.com",
			Name:      "User",
			Gender:    1,
			CreatedAt: time.Date(2023, 10, 1, 0, 0, 0, 0, time.UTC),
			UpdatedAt: time.Date(2023, 10, 1, 0, 0, 0, 0, time.UTC),
		}

		handler := handlers.NewUserHandler(userService, redisService, bcryptService)
		router.GET("/users/:id", handler.GetUser)

		userService.On("GetUser", uint(1)).Return(user, nil)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/users/1", nil)
		req.Header.Set("Authorization", "Bearer test-token")

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var actualBody map[string]any
		err := json.Unmarshal(w.Body.Bytes(), &actualBody)
		require.NoError(t, err)

		expectedBody := map[string]any{
			"id":        float64(1),
			"email":     "email@example.com",
			"name":      "User",
			"gender":    float64(1),
			"createdAt": "2023-10-01T00:00:00Z",
			"updatedAt": "2023-10-01T00:00:00Z",
			"deletedAt": nil,
		}

		assert.Equal(t, expectedBody, actualBody)
	})

	t.Run("User Not Found", func(t *testing.T) {
		userService := new(mocks.MockUserService)
		redisService := new(mocks.MockRedisService)
		bcryptService := new(mocks.MockBcryptService)

		router := gin.Default()

		handler := handlers.NewUserHandler(userService, redisService, bcryptService)
		router.GET("/users/:id", handler.GetUser)

		userService.On("GetUser", uint(1)).Return(&models.User{}, appError.New(appError.ErrNotFound, "User not found"))

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/users/1", nil)
		req.Header.Set("Authorization", "Bearer test-token")

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
		assert.JSONEq(t, fmt.Sprintf(`{"code":%d,"message":"User not found"}`, appError.ErrNotFound), w.Body.String())
	})

	t.Run("Invalid User ID", func(t *testing.T) {
		userService := new(mocks.MockUserService)
		redisService := new(mocks.MockRedisService)
		bcryptService := new(mocks.MockBcryptService)

		router := gin.Default()

		handler := handlers.NewUserHandler(userService, redisService, bcryptService)
		router.GET("/users/:id", handler.GetUser)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/users/invalid-id", nil)
		req.Header.Set("Authorization", "Bearer test-token")

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.JSONEq(t, fmt.Sprintf(`{"code":%d,"message":"strconv.Atoi: parsing \"invalid-id\": invalid syntax"}`, appError.ErrParseError), w.Body.String())
	})
}

func TestChangePassword(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Initialize the validator
	utils.InitValidator()

	t.Run("Success", func(t *testing.T) {
		userService := new(mocks.MockUserService)
		redisService := new(mocks.MockRedisService)
		bcryptService := new(mocks.MockBcryptService)

		user := &models.User{
			ID:        1,
			Email:     "email@example.com",
			Name:      "User",
			Password:  "$2a$10$I/L5VegpCyOlJPoa1.KrmeCdezSBIandsEL5S2dd4Ap0YIWk0Iuka", // bcrypt hash of "12345678"
			Gender:    1,
			CreatedAt: time.Date(2023, 10, 1, 0, 0, 0, 0, time.UTC),
			UpdatedAt: time.Date(2023, 10, 1, 0, 0, 0, 0, time.UTC),
		}

		handler := handlers.NewUserHandler(userService, redisService, bcryptService)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		requestBody := map[string]any{
			"old_password":     "12345678",
			"new_password":     "newpassword",
			"confirm_password": "newpassword",
		}
		body, _ := json.Marshal(requestBody)
		c.Request, _ = http.NewRequest("PUT", "/api/v1/change-password", bytes.NewBuffer(body))
		c.Request.Header.Set("Content-Type", "application/json")
		c.Request.Header.Set("Authorization", "Bearer test-token")

		c.Set("UserID", uint(1))

		// mock method GetUser to return a user
		userService.On("GetUser", uint(1)).Return(user, nil)
		// mock method UpdateUser to return nil
		userService.On("UpdateUser", user).Return(nil)
		// mock method BcryptCompare to return nil (passwords match)
		bcryptService.On("CheckPasswordHash", "12345678", user.Password).Return(true)
		// mock method HashPassword to return a new hashed password
		bcryptService.On("HashPassword", "newpassword").Return("$2a$10$hashedNewPassword", nil)

		// call the ChangePassword handler
		handler.ChangePassword(c)
		// Assert the response
		assert.Equal(t, http.StatusOK, w.Code)
		assert.JSONEq(t, `{"message":"Change password successfully"}`, w.Body.String())
	})

	t.Run("Validation Error", func(t *testing.T) {
		tests := []struct {
			name         string
			reqBody      string
			expectedCode float64
			expectedMsg  string
		}{
			{
				name:         "EmptyOldPassword",
				reqBody:      `{"old_password":""}`,
				expectedCode: float64(4001),
				expectedMsg:  "oldpassword is required",
			},
			{
				name:         "ShortOldPassword",
				reqBody:      `{"old_password":"short"}`,
				expectedCode: float64(4001),
				expectedMsg:  "oldpassword must be at least 6 characters long or numeric",
			},
			{
				name:         "LongOldPassword",
				reqBody:      `{"old_password":"` + strings.Repeat("a", 256) + `"}`,
				expectedCode: float64(4001),
				expectedMsg:  "oldpassword must be at most 255 characters long or numeric",
			},
			{
				name:         "EmptyNewPassword",
				reqBody:      `{"old_password":"12345678","new_password":""}`,
				expectedCode: float64(4001),
				expectedMsg:  "newpassword is required",
			},
			{
				name:         "ShortNewPassword",
				reqBody:      `{"old_password":"12345678","new_password":"short"}`,
				expectedCode: float64(4001),
				expectedMsg:  "newpassword must be at least 6 characters long or numeric",
			},
			{
				name:         "LongNewPassword",
				reqBody:      `{"old_password":"12345678","new_password":"` + strings.Repeat("a", 256) + `"}`,
				expectedCode: float64(4001),
				expectedMsg:  "newpassword must be at most 255 characters long or numeric",
			},
			{
				name:         "EmptyConfirmPassword",
				reqBody:      `{"old_password":"12345678","new_password":"newpassword","confirm_password":""}`,
				expectedCode: float64(4001),
				expectedMsg:  "confirmpassword is required",
			},
			{
				name:         "ShortConfirmPassword",
				reqBody:      `{"old_password":"12345678","new_password":"newpassword","confirm_password":"short"}`,
				expectedCode: float64(4001),
				expectedMsg:  "confirmpassword must be at least 6 characters long or numeric",
			},
			{
				name:         "LongConfirmPassword",
				reqBody:      `{"old_password":"12345678","new_password":"newpassword","confirm_password":"` + strings.Repeat("a", 256) + `"}`,
				expectedCode: float64(4001),
				expectedMsg:  "confirmpassword must be at most 255 characters long or numeric",
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				userService := new(mocks.MockUserService)
				redisService := new(mocks.MockRedisService)
				bcryptService := new(mocks.MockBcryptService)

				handler := handlers.NewUserHandler(userService, redisService, bcryptService)
				w := httptest.NewRecorder()
				c, _ := gin.CreateTestContext(w)

				c.Request, _ = http.NewRequest("PUT", "/api/v1/change-password", bytes.NewBufferString(tt.reqBody))
				c.Request.Header.Set("Content-Type", "application/json")
				c.Request.Header.Set("Authorization", "Bearer test-token")
				c.Set("UserID", uint(1))

				// Call the ChangePassword handler
				handler.ChangePassword(c)

				// Assert the response
				assert.Equal(t, http.StatusBadRequest, w.Code)

				expectedBody := map[string]any{
					"code":    tt.expectedCode,
					"message": tt.expectedMsg,
				}

				var actualBody map[string]any
				err := json.Unmarshal(w.Body.Bytes(), &actualBody)
				require.NoError(t, err)

				assert.Equal(t, expectedBody["code"], actualBody["code"])
				assert.Equal(t, expectedBody["message"], actualBody["message"])
			})
		}
	})

	t.Run("Error User Not Found", func(t *testing.T) {
		userService := new(mocks.MockUserService)
		redisService := new(mocks.MockRedisService)
		bcryptService := new(mocks.MockBcryptService)

		requestBody := map[string]any{
			"old_password":     "12345678",
			"new_password":     "newpassword",
			"confirm_password": "newpassword",
		}
		body, _ := json.Marshal(requestBody)

		handler := handlers.NewUserHandler(userService, redisService, bcryptService)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest("PUT", "/api/v1/change-password", bytes.NewBuffer(body))
		c.Request.Header.Set("Content-Type", "application/json")
		c.Request.Header.Set("Authorization", "Bearer test-token")
		c.Set("UserID", uint(1))

		// Mock the GetUser method to return an error
		userService.On("GetUser", uint(1)).Return(&models.User{}, appError.New(appError.ErrNotFound, "User not found"))

		// Call the ChangePassword handler
		handler.ChangePassword(c)

		// Assert the response
		assert.Equal(t, http.StatusBadRequest, w.Code)

		expectedBody := map[string]any{
			"code":    float64(appError.ErrNotFound),
			"message": "User not found",
		}

		var actualBody map[string]any
		err := json.Unmarshal(w.Body.Bytes(), &actualBody)
		require.NoError(t, err)

		assert.Equal(t, expectedBody["code"], actualBody["code"])
		assert.Equal(t, expectedBody["message"], actualBody["message"])
	})

	t.Run("Error Old Password Mismatch", func(t *testing.T) {
		userService := new(mocks.MockUserService)
		redisService := new(mocks.MockRedisService)
		bcryptService := new(mocks.MockBcryptService)

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

		handler := handlers.NewUserHandler(userService, redisService, bcryptService)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest("PUT", "/api/v1/change-password", bytes.NewBuffer(body))
		c.Request.Header.Set("Content-Type", "application/json")
		c.Request.Header.Set("Authorization", "Bearer test-token")
		c.Set("UserID", uint(1))

		// Mock the GetUser method to return the user
		userService.On("GetUser", uint(1)).Return(user, nil)
		// Mock the Bcrypt.CheckPasswordHash method to return an error (passwords do not match)
		bcryptService.On("CheckPasswordHash", "wrongpassword", user.Password).Return(false)

		// Call the ChangePassword handler
		handler.ChangePassword(c)

		// Assert the response
		assert.Equal(t, http.StatusBadRequest, w.Code)

		expectedBody := map[string]any{
			"code":    float64(appError.ErrInvalidPassword),
			"message": "Old password is incorrect",
		}

		var actualBody map[string]any
		err := json.Unmarshal(w.Body.Bytes(), &actualBody)
		require.NoError(t, err)

		assert.Equal(t, expectedBody["code"], actualBody["code"])
		assert.Equal(t, expectedBody["message"], actualBody["message"])
	})

	t.Run("Error new_password and confirm_password mismatch", func(t *testing.T) {
		userService := new(mocks.MockUserService)
		redisService := new(mocks.MockRedisService)
		bcryptService := new(mocks.MockBcryptService)

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
		handler := handlers.NewUserHandler(userService, redisService, bcryptService)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest("PUT", "/api/v1/change-password", bytes.NewBuffer(body))
		c.Request.Header.Set("Content-Type", "application/json")
		c.Request.Header.Set("Authorization", "Bearer test-token")
		c.Set("UserID", uint(1))
		// Mock the GetUser method to return the user
		userService.On("GetUser", uint(1)).Return(user, nil)
		// Mock the Bcrypt.CheckPasswordHash method to return true (passwords match)
		bcryptService.On("CheckPasswordHash", requestBody["old_password"], user.Password).Return(true)

		// Call the ChangePassword handler
		handler.ChangePassword(c)
		// Assert the response
		assert.Equal(t, http.StatusBadRequest, w.Code)

		expectedBody := map[string]any{
			"code":    float64(appError.ErrPasswordMismatch),
			"message": "New password and confirm password do not match",
		}
		var actualBody map[string]any
		err := json.Unmarshal(w.Body.Bytes(), &actualBody)
		require.NoError(t, err)
		assert.Equal(t, expectedBody["code"], actualBody["code"])
		assert.Equal(t, expectedBody["message"], actualBody["message"])

	})

	t.Run("Error Update User", func(t *testing.T) {
		userService := new(mocks.MockUserService)
		redisService := new(mocks.MockRedisService)
		bcryptService := new(mocks.MockBcryptService)

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
		handler := handlers.NewUserHandler(userService, redisService, bcryptService)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest("PUT", "/api/v1/change-password", bytes.NewBuffer(body))
		c.Request.Header.Set("Content-Type", "application/json")
		c.Request.Header.Set("Authorization", "Bearer test-token")
		c.Set("UserID", uint(1))

		// Mock the GetUser method to return the user
		userService.On("GetUser", uint(1)).Return(user, nil)
		// Mock the UpdateUser method to return an error
		userService.On("UpdateUser", user).Return(appError.New(appError.ErrDBUpdate, "Update error"))
		// mock the Bcrypt.CheckPasswordHash method to return true (passwords match)
		bcryptService.On("CheckPasswordHash", "12345678", user.Password).Return(true)
		// mock the Bcrypt.HashPassword method to return a new hashed password
		bcryptService.On("HashPassword", "newpassword").Return("$2a$10$hashedNewPassword", nil)

		// Call the ChangePassword handler
		handler.ChangePassword(c)

		// Assert the response
		assert.Equal(t, http.StatusInternalServerError, w.Code)

		expectedBody := map[string]any{
			"code":    float64(appError.ErrDBUpdate),
			"message": "Update error",
		}

		var actualBody map[string]any
		err := json.Unmarshal(w.Body.Bytes(), &actualBody)
		require.NoError(t, err)

		assert.Equal(t, expectedBody["code"], actualBody["code"])
		assert.Equal(t, expectedBody["message"], actualBody["message"])
	})

	t.Run("Error User Not found from ctx", func(t *testing.T) {
		userService := new(mocks.MockUserService)
		redisService := new(mocks.MockRedisService)
		bcryptService := new(mocks.MockBcryptService)

		handler := handlers.NewUserHandler(userService, redisService, bcryptService)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest("PUT", "/api/v1/change-password", nil)
		c.Request.Header.Set("Content-Type", "application/json")
		c.Request.Header.Set("Authorization", "Bearer test-token")
		c.Set("UserID", uint(0)) // Invalid User ID

		// Call the ChangePassword handler
		handler.ChangePassword(c)

		// Assert the response
		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.JSONEq(t, fmt.Sprintf(`{"code":%d,"message":"Invalid UserID"}`, appError.ErrParseError), w.Body.String())
	})

	t.Run("Old Password equal to New Password", func(t *testing.T) {
		userService := new(mocks.MockUserService)
		redisService := new(mocks.MockRedisService)
		bcryptService := new(mocks.MockBcryptService)

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
		handler := handlers.NewUserHandler(userService, redisService, bcryptService)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest("PUT", "/api/v1/change-password", bytes.NewBuffer(body))
		c.Request.Header.Set("Content-Type", "application/json")
		c.Request.Header.Set("Authorization", "Bearer test-token")
		c.Set("UserID", uint(1))
		// Mock the GetUser method to return the user
		userService.On("GetUser", uint(1)).Return(user, nil)
		// Mock the Bcrypt.CheckPasswordHash method to return true (passwords match)
		bcryptService.On("CheckPasswordHash", requestBody["old_password"], user.Password).Return(true)
		// Call the ChangePassword handler
		handler.ChangePassword(c)
		// Assert the response
		assert.Equal(t, http.StatusBadRequest, w.Code)
		expectedBody := map[string]any{
			"code":    float64(appError.ErrPasswordMismatch),
			"message": "New password must be different from old password",
		}
		var actualBody map[string]any
		err := json.Unmarshal(w.Body.Bytes(), &actualBody)
		require.NoError(t, err)
		assert.Equal(t, expectedBody["code"], actualBody["code"])
		assert.Equal(t, expectedBody["message"], actualBody["message"])
	})

	t.Run("Error Bcrypt Hashing", func(t *testing.T) {
		userService := new(mocks.MockUserService)
		redisService := new(mocks.MockRedisService)
		bcryptService := new(mocks.MockBcryptService)

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
		handler := handlers.NewUserHandler(userService, redisService, bcryptService)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest("PUT", "/api/v1/change-password", bytes.NewBuffer(body))
		c.Request.Header.Set("Content-Type", "application/json")
		c.Request.Header.Set("Authorization", "Bearer test-token")
		c.Set("UserID", uint(1))
		// Mock the GetUser method to return the user
		userService.On("GetUser", uint(1)).Return(user, nil)
		// Mock the Bcrypt.CheckPasswordHash method to return true
		bcryptService.On("CheckPasswordHash", requestBody["old_password"], user.Password).Return(true)
		// Mock the Bcrypt.HashPassword method to return an error
		bcryptService.On("HashPassword", "newpassword").Return("", appError.New(appError.ErrPasswordHashFailed, "Hash password failed"))
		// Call the ChangePassword handler
		handler.ChangePassword(c)
		// Assert the response
		assert.Equal(t, http.StatusInternalServerError, w.Code)
		expectedBody := map[string]any{
			"code":    float64(appError.ErrPasswordHashFailed),
			"message": "Hash password failed",
		}
		var actualBody map[string]any
		err := json.Unmarshal(w.Body.Bytes(), &actualBody)
		require.NoError(t, err)
		assert.Equal(t, expectedBody["code"], actualBody["code"])
		assert.Equal(t, expectedBody["message"], actualBody["message"])
	})

}

func TestUpdateUser(t *testing.T) {
	gin.SetMode(gin.TestMode)

	utils.InitValidator()

	t.Run("Success Update User", func(t *testing.T) {
		router := gin.Default()

		userService := new(mocks.MockUserService)
		redisService := new(mocks.MockRedisService)
		bcryptService := new(mocks.MockBcryptService)

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

		handler := handlers.NewUserHandler(userService, redisService, bcryptService)
		router.PATCH("/api/v1/users/:id", handler.UpdateUser)

		req, _ := http.NewRequest("PATCH", "/api/v1/users/1", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer test-token")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.JSONEq(t, `{"message":"Update user successfully"}`, w.Body.String())
	})
	t.Run("Validation Error", func(t *testing.T) {
		tests := []struct {
			name         string
			reqBody      string
			expectedCode float64
			expectedMsg  string
		}{
			{
				name:         "EmptyName",
				reqBody:      `{"name":""}`,
				expectedCode: 4001,
				expectedMsg:  "name must be at least 1 characters long or numeric",
			},
			{
				name:         "NameNotBlank",
				reqBody:      `{"name":"  "}`,
				expectedCode: 4001,
				expectedMsg:  "name must not be blank",
			},
			{
				name:         "LongName",
				reqBody:      `{"name": "` + strings.Repeat("a", 46) + `"}`,
				expectedCode: float64(4001),
				expectedMsg:  "name must be at most 45 characters long or numeric",
			},
			{
				name:         "InvalidBirthdayFormat",
				reqBody:      `{"name":"User","birthday":"invalid-date"}`,
				expectedCode: 4001,
				expectedMsg:  "birthday must be a valid date (YYYY-MM-DD) and not in the future",
			},
			{
				name:         "FutureBirthday",
				reqBody:      `{"name":"User","birthday":"3000-01-01"}`,
				expectedCode: 4001,
				expectedMsg:  "birthday must be a valid date (YYYY-MM-DD) and not in the future",
			},
			{
				name:         "EmptyAddress",
				reqBody:      `{"name":"User","birthday":"2000-01-01","address":""}`,
				expectedCode: 4001,
				expectedMsg:  "address must be at least 1 characters long or numeric",
			},
			{
				name:         "LongAddress",
				reqBody:      `{"name":"User","birthday":"2000-01-01","address":"` + strings.Repeat("a", 256) + `"}`,
				expectedCode: 4001,
				expectedMsg:  "address must be at most 255 characters long or numeric",
			},
			{
				name:         "AddressNotBlank",
				reqBody:      `{"name":"User","birthday":"2000-01-01","address":"  "}`,
				expectedCode: 4001,
				expectedMsg:  "address must not be blank",
			},
			{
				name:         "InvalidGender 4",
				reqBody:      `{"name":"User","birthday":"2000-01-01","address":"123 Street","gender":4}`,
				expectedCode: 4001,
				expectedMsg:  "gender must be one of [1 2 3]",
			},
			{
				name:         "InvalidGender 0",
				reqBody:      `{"name":"User","birthday":"2000-01-01","address":"123 Street","gender":0}`,
				expectedCode: 4001,
				expectedMsg:  "gender must be one of [1 2 3]",
			},
			{
				name:         "StringGender",
				reqBody:      `{"name":"User","birthday":"2000-01-01","address":"123 Street","gender":"male"}`,
				expectedCode: 4001,
				expectedMsg:  "json: cannot unmarshal string into Go struct field .gender of type int16",
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				// Setup router
				router := gin.Default()
				// Mock services
				userService := new(mocks.MockUserService)
				redisService := new(mocks.MockRedisService)
				bcryptService := new(mocks.MockBcryptService)

				handler := handlers.NewUserHandler(userService, redisService, bcryptService)
				router.PATCH("/api/v1/users/:id", handler.UpdateUser)

				w := httptest.NewRecorder()
				req, _ := http.NewRequest("PATCH", "/api/v1/users/1", bytes.NewBufferString(tt.reqBody))
				req.Header.Set("Content-Type", "application/json")
				req.Header.Set("Authorization", "Bearer test-token")

				router.ServeHTTP(w, req)

				assert.Equal(t, http.StatusBadRequest, w.Code)

				var actualBody map[string]any
				err := json.Unmarshal(w.Body.Bytes(), &actualBody)
				require.NoError(t, err)

				assert.Equal(t, tt.expectedCode, actualBody["code"])
				assert.Equal(t, tt.expectedMsg, actualBody["message"])

			})
		}
	})

	t.Run("Error Parse ID", func(t *testing.T) {
		router := gin.Default()

		userService := new(mocks.MockUserService)
		redisService := new(mocks.MockRedisService)
		bcryptService := new(mocks.MockBcryptService)

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
		userService.On("UpdateUser", user).Return(nil)

		handler := handlers.NewUserHandler(userService, redisService, bcryptService)
		router.PATCH("/api/v1/users/:id", func(c *gin.Context) {
			c.Set("UserID", uint(1))
			handler.UpdateUser(c)
		})

		req, _ := http.NewRequest("PATCH", "/api/v1/users/invalid-id", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer test-token")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		var actualBody map[string]any
		err := json.Unmarshal(w.Body.Bytes(), &actualBody)
		require.NoError(t, err)

		expectedBody := map[string]any{
			"code":    float64(appError.ErrParseError),
			"message": "strconv.Atoi: parsing \"invalid-id\": invalid syntax",
		}

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Equal(t, expectedBody, actualBody)
	})

	t.Run("Error User Not Found", func(t *testing.T) {
		router := gin.Default()

		userService := new(mocks.MockUserService)
		redisService := new(mocks.MockRedisService)
		bcryptService := new(mocks.MockBcryptService)

		var requestBody = map[string]any{
			"name":     "Updated User",
			"birthday": "2000-01-01",
			"address":  "456 New Street",
			"gender":   1,
		}
		body, _ := json.Marshal(requestBody)

		userService.On("GetUser", uint(1)).Return(&models.User{}, appError.New(appError.ErrNotFound, "User not found"))

		handler := handlers.NewUserHandler(userService, redisService, bcryptService)
		router.PATCH("/api/v1/users/:id", func(c *gin.Context) {
			c.Set("UserID", uint(1))
			handler.UpdateUser(c)
		})

		req, _ := http.NewRequest("PATCH", "/api/v1/users/1", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer test-token")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		var actualBody map[string]any
		err := json.Unmarshal(w.Body.Bytes(), &actualBody)
		require.NoError(t, err)

		expectedBody := map[string]any{
			"code":    float64(appError.ErrNotFound),
			"message": "User not found",
		}

		assert.Equal(t, http.StatusNotFound, w.Code)
		assert.Equal(t, expectedBody, actualBody)
	})

	t.Run("Error Update User", func(t *testing.T) {
		router := gin.Default()

		userService := new(mocks.MockUserService)
		redisService := new(mocks.MockRedisService)
		bcryptService := new(mocks.MockBcryptService)

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
		userService.On("UpdateUser", user).Return(appError.New(appError.ErrDBUpdate, "Update error"))

		handler := handlers.NewUserHandler(userService, redisService, bcryptService)
		router.PATCH("/api/v1/users/:id", func(c *gin.Context) {
			c.Set("UserID", uint(1))
			handler.UpdateUser(c)
		})
		req, _ := http.NewRequest("PATCH", "/api/v1/users/1", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer test-token")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		var actualBody map[string]any
		err := json.Unmarshal(w.Body.Bytes(), &actualBody)
		require.NoError(t, err)
		expectedBody := map[string]any{
			"code":    float64(appError.ErrDBUpdate),
			"message": "Update error",
		}

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Equal(t, expectedBody, actualBody)
	})

}

func TestDeleteUser(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("Success Delete User", func(t *testing.T) {
		router := gin.Default()
		userService := new(mocks.MockUserService)
		redisService := new(mocks.MockRedisService)
		bcryptService := new(mocks.MockBcryptService)

		user := &models.User{
			ID:    1,
			Email: "email@example.com",
			Name:  "User",
		}
		// Mock the GetUser method to return a user
		userService.On("GetUser", uint(1)).Return(user, nil)
		// Mock the DeleteUser method to return nil
		userService.On("DeleteUser", uint(1)).Return(nil)

		// Create a new UserHandler with parameters
		handler := handlers.NewUserHandler(userService, redisService, bcryptService)
		router.DELETE("/api/v1/users/:id", func(ctx *gin.Context) {
			// set the UserID in the context
			ctx.Set("UserID", uint(1))
			// Call the DeleteUser handler
			handler.DeleteUser(ctx)
		})

		req, _ := http.NewRequest("DELETE", "/api/v1/users/1", nil)
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer test-token")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		// Assert the response
		assert.Equal(t, http.StatusOK, w.Code)
		assert.JSONEq(t, `{"message":"Delete user successfully"}`, w.Body.String())
	})

	t.Run("Parse ID from URL error", func(t *testing.T) {
		router := gin.Default()
		userService := new(mocks.MockUserService)
		redisService := new(mocks.MockRedisService)
		bcryptService := new(mocks.MockBcryptService)

		handler := handlers.NewUserHandler(userService, redisService, bcryptService)

		router.DELETE("/api/v1/users/:id", func(ctx *gin.Context) {
			ctx.Set("UserID", uint(1))
			handler.DeleteUser(ctx)
		})

		req, _ := http.NewRequest("DELETE", "/api/v1/users/invalid-id", nil)
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer test-token")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		var actualBody map[string]any
		err := json.Unmarshal(w.Body.Bytes(), &actualBody)
		require.NoError(t, err)

		expectedBody := map[string]any{
			"code":    float64(appError.ErrParseError),
			"message": "Invalid UserID",
		}

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Equal(t, expectedBody, actualBody)
	})

	t.Run("User Not Found", func(t *testing.T) {
		router := gin.Default()
		userService := new(mocks.MockUserService)
		redisService := new(mocks.MockRedisService)
		bcryptService := new(mocks.MockBcryptService)

		// Mock the GetUser method to return an error
		userService.On("GetUser", uint(1)).Return(&models.User{}, appError.New(appError.ErrNotFound, "User not found"))

		handler := handlers.NewUserHandler(userService, redisService, bcryptService)
		router.DELETE("/api/v1/users/:id", func(ctx *gin.Context) {
			ctx.Set("UserID", uint(1))
			handler.DeleteUser(ctx)
		})

		req, _ := http.NewRequest("DELETE", "/api/v1/users/1", nil)
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer test-token")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		var actualBody map[string]any
		err := json.Unmarshal(w.Body.Bytes(), &actualBody)
		require.NoError(t, err)

		expectedBody := map[string]any{
			"code":    float64(appError.ErrNotFound),
			"message": "User not found",
		}

		assert.Equal(t, http.StatusNotFound, w.Code)
		assert.Equal(t, expectedBody, actualBody)
	})

	t.Run("Error Delete User", func(t *testing.T) {
		router := gin.Default()
		userService := new(mocks.MockUserService)
		redisService := new(mocks.MockRedisService)
		bcryptService := new(mocks.MockBcryptService)

		user := &models.User{
			ID:    1,
			Email: "email@example.com",
			Name:  "User",
		}
		// Mock the GetUser method to return a user
		userService.On("GetUser", uint(1)).Return(user, nil)
		// Mock the DeleteUser method to return an error
		userService.On("DeleteUser", uint(1)).Return(appError.New(appError.ErrDBDelete, "Delete error"))
		// Create a new UserHandler with parameters
		handler := handlers.NewUserHandler(userService, redisService, bcryptService)
		router.DELETE("/api/v1/users/:id", func(ctx *gin.Context) {
			// set the UserID in the context
			ctx.Set("UserID", uint(1))
			// Call the DeleteUser handler
			handler.DeleteUser(ctx)
		})
		req, _ := http.NewRequest("DELETE", "/api/v1/users/1", nil)
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer test-token")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		var actualBody map[string]any
		err := json.Unmarshal(w.Body.Bytes(), &actualBody)
		require.NoError(t, err)
		expectedBody := map[string]any{
			"code":    float64(appError.ErrDBDelete),
			"message": "Delete error",
		}
		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.Equal(t, expectedBody, actualBody)

	})
}

func TestResetPassword(t *testing.T) {
	gin.SetMode(gin.TestMode)

	utils.InitValidator()

	t.Run("Success Reset Password", func(t *testing.T) {
		userService := new(mocks.MockUserService)
		redisService := new(mocks.MockRedisService)
		bcryptService := new(mocks.MockBcryptService)

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
		handler := handlers.NewUserHandler(userService, redisService, bcryptService)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest("POST", "/api/v1/reset-password", bytes.NewBuffer(body))
		c.Request.Header.Set("Content-Type", "application/json")
		c.Request.Header.Set("Authorization", "Bearer test-token")
		c.Set("UserID", uint(1))

		// Mock the GetUserByToken method to return a user
		userService.On("GetUserByToken", "token").Return(user, nil)
		// Mock the Bcrypt.CheckPasswordHash method to return true (passwords match)
		bcryptService.On("CheckPasswordHash", "newpassword", user.Password).Return(true)
		// Mock the Bcrypt.HashPassword method to return a new hashed password
		bcryptService.On("HashPassword", "newpassword").Return("$2a$10$hashedNewPassword", nil)
		// Mock the UpdateUser method to return nil
		userService.On("UpdateUser", user).Return(nil)

		// Call the ResetPassword handler
		handler.ResetPassword(c)

		// Assert the response
		assert.Equal(t, http.StatusOK, w.Code)
		assert.JSONEq(t, `{"message":"Reset password successfully"}`, w.Body.String())

	})

	t.Run("Not found user by token", func(t *testing.T) {
		userService := new(mocks.MockUserService)
		redisService := new(mocks.MockRedisService)
		bcryptService := new(mocks.MockBcryptService)

		requestBody := map[string]any{
			"token":        "invalid-token",
			"password":     "newpassword",
			"new_password": "newpassword",
		}
		body, _ := json.Marshal(requestBody)

		handler := handlers.NewUserHandler(userService, redisService, bcryptService)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest("POST", "/api/v1/reset-password", bytes.NewBuffer(body))
		c.Request.Header.Set("Content-Type", "application/json")
		c.Request.Header.Set("Authorization", "Bearer test-token")
		c.Set("UserID", uint(1))

		// Mock the GetUserByToken method to return an error
		userService.On("GetUserByToken", "invalid-token").Return(&models.User{}, appError.New(appError.ErrNotFound, "User not found"))

		// Call the ResetPassword handler
		handler.ResetPassword(c)

		// Assert the response
		assert.Equal(t, http.StatusNotFound, w.Code)

		expectedBody := map[string]any{
			"code":    float64(appError.ErrNotFound),
			"message": "User not found",
		}

		var actualBody map[string]any
		err := json.Unmarshal(w.Body.Bytes(), &actualBody)
		require.NoError(t, err)

		assert.Equal(t, expectedBody["code"], actualBody["code"])
		assert.Equal(t, expectedBody["message"], actualBody["message"])
	})

	t.Run("Error Token expired", func(t *testing.T) {
		userService := new(mocks.MockUserService)
		redisService := new(mocks.MockRedisService)
		bcryptService := new(mocks.MockBcryptService)

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

		handler := handlers.NewUserHandler(userService, redisService, bcryptService)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest("POST", "/api/v1/reset-password", bytes.NewBuffer(body))
		c.Request.Header.Set("Content-Type", "application/json")
		c.Request.Header.Set("Authorization", "Bearer test-token")
		c.Set("UserID", uint(1))

		// Mock the GetUserByToken method to return a user
		userService.On("GetUserByToken", "invalid-token").Return(user, nil)

		// Call the ResetPassword handler
		handler.ResetPassword(c)

		// Assert the response
		assert.Equal(t, http.StatusBadRequest, w.Code)
		expectedBody := map[string]any{
			"code":    float64(appError.ErrTokenExpired),
			"message": "Token expired",
		}

		var actualBody map[string]any
		err := json.Unmarshal(w.Body.Bytes(), &actualBody)
		require.NoError(t, err)

		assert.Equal(t, expectedBody["code"], actualBody["code"])
		assert.Equal(t, expectedBody["message"], actualBody["message"])

	})

	t.Run("Error Passwords incorrect", func(t *testing.T) {
		userService := new(mocks.MockUserService)
		redisService := new(mocks.MockRedisService)
		bcryptService := new(mocks.MockBcryptService)

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

		handler := handlers.NewUserHandler(userService, redisService, bcryptService)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest("POST", "/api/v1/reset-password", bytes.NewBuffer(body))
		c.Request.Header.Set("Content-Type", "application/json")
		c.Request.Header.Set("Authorization", "Bearer test-token")
		c.Set("UserID", uint(1))

		// Mock the GetUserByToken method to return a user
		userService.On("GetUserByToken", "token").Return(user, nil)
		// Mock the Bcrypt.CheckPasswordHash method to return false (passwords do not match)
		bcryptService.On("CheckPasswordHash", "newpassword", user.Password).Return(false)

		// Call the ResetPassword handler
		handler.ResetPassword(c)

		// Assert the response
		assert.Equal(t, http.StatusBadRequest, w.Code)
		expectedBody := map[string]any{
			"code":    float64(appError.ErrInvalidPassword),
			"message": "Password is incorrect",
		}

		var actualBody map[string]any
		err := json.Unmarshal(w.Body.Bytes(), &actualBody)
		require.NoError(t, err)

		assert.Equal(t, expectedBody["code"], actualBody["code"])
		assert.Equal(t, expectedBody["message"], actualBody["message"])

	})

	t.Run("Error hashed password failed", func(t *testing.T) {
		userService := new(mocks.MockUserService)
		redisService := new(mocks.MockRedisService)
		bcryptService := new(mocks.MockBcryptService)

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

		handler := handlers.NewUserHandler(userService, redisService, bcryptService)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest("POST", "/api/v1/reset-password", bytes.NewBuffer(body))
		c.Request.Header.Set("Content-Type", "application/json")
		c.Request.Header.Set("Authorization", "Bearer test-token")
		c.Set("UserID", uint(1))

		// Mock the GetUserByToken method to return a user
		userService.On("GetUserByToken", "token").Return(user, nil)
		// Mock the Bcrypt.CheckPasswordHash method to return true
		bcryptService.On("CheckPasswordHash", "newpassword", user.Password).Return(true)
		// Mock the Bcrypt.HashPassword method to return an error
		bcryptService.On("HashPassword", "newpassword").Return("", appError.New(appError.ErrPasswordHashFailed, "Hash password failed"))

		// Call the ResetPassword handler
		handler.ResetPassword(c)

		// Assert the response
		assert.Equal(t, http.StatusInternalServerError, w.Code)
		expectedBody := map[string]any{
			"code":    float64(appError.ErrPasswordHashFailed),
			"message": "Failed to hash password",
		}

		var actualBody map[string]any
		err := json.Unmarshal(w.Body.Bytes(), &actualBody)
		require.NoError(t, err)

		assert.Equal(t, expectedBody["code"], actualBody["code"])
		assert.Equal(t, expectedBody["message"], actualBody["message"])

	})

	t.Run("Error failed to UpdateUser", func(t *testing.T) {
		userService := new(mocks.MockUserService)
		redisService := new(mocks.MockRedisService)
		bcryptService := new(mocks.MockBcryptService)

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

		handler := handlers.NewUserHandler(userService, redisService, bcryptService)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest("POST", "/api/v1/reset-password", bytes.NewBuffer(body))
		c.Request.Header.Set("Content-Type", "application/json")
		c.Request.Header.Set("Authorization", "Bearer test-token")
		c.Set("UserID", uint(1))

		// Mock the GetUserByToken method to return a user
		userService.On("GetUserByToken", "token").Return(user, nil)
		// Mock the Bcrypt.CheckPasswordHash method to return true
		bcryptService.On("CheckPasswordHash", "newpassword", user.Password).Return(true)
		// Mock the Bcrypt.HashPassword method to return an hashed password
		bcryptService.On("HashPassword", "newpassword").Return("$2a$10$hashedNewPassword", nil)
		// Mock the UpdateUser method to return an error
		userService.On("UpdateUser", user).Return(appError.New(appError.ErrDBUpdate, "Failed to update user"))

		// Call the ResetPassword handler
		handler.ResetPassword(c)

		// Assert the response
		assert.Equal(t, http.StatusBadRequest, w.Code)
		expectedBody := map[string]any{
			"code":    float64(appError.ErrDBUpdate),
			"message": "Failed to update user",
		}

		var actualBody map[string]any
		err := json.Unmarshal(w.Body.Bytes(), &actualBody)
		require.NoError(t, err)

		assert.Equal(t, expectedBody["code"], actualBody["code"])
		assert.Equal(t, expectedBody["message"], actualBody["message"])
	})

	t.Run("Validation Error", func(t *testing.T) {
		tests := []struct {
			name         string
			reqBody      string
			expectedCode float64
			expectedMsg  string
		}{
			{
				name:         "EmptyToken",
				reqBody:      `{"token":""}`,
				expectedCode: 4001,
				expectedMsg:  "token is required",
			},
			{
				name:         "EmptyPassword",
				reqBody:      `{"token":"valid-token","password":""}`,
				expectedCode: 4001,
				expectedMsg:  "password is required",
			},
			{
				name:         "PasswordTooShort",
				reqBody:      `{"token":"valid-token","password":"short"}`,
				expectedCode: 4001,
				expectedMsg:  "password must be at least 6 characters long or numeric",
			},
			{
				name:         "PasswordTooLong",
				reqBody:      `{"token":"valid-token","password":"` + strings.Repeat("a", 256) + `"}`,
				expectedCode: 4001,
				expectedMsg:  "password must be at most 255 characters long or numeric",
			},
			{
				name:         "EmptyNewPassword",
				reqBody:      `{"token":"valid-token","password":"newpassword","new_password":""}`,
				expectedCode: 4001,
				expectedMsg:  "newpassword is required",
			},
			{
				name:         "NewPasswordTooShort",
				reqBody:      `{"token":"valid-token","password":"newpassword","new_password":"short"}`,
				expectedCode: 4001,
				expectedMsg:  "newpassword must be at least 6 characters long or numeric",
			},
			{
				name:         "NewPasswordTooLong",
				reqBody:      `{"token":"valid-token","password":"newpassword","new_password":"` + strings.Repeat("a", 256) + `"}`,
				expectedCode: 4001,
				expectedMsg:  "newpassword must be at most 255 characters long or numeric",
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				// Setup router
				// Mock services
				userService := new(mocks.MockUserService)
				redisService := new(mocks.MockRedisService)
				bcryptService := new(mocks.MockBcryptService)

				// Create a new UserHandler with parameters
				handler := handlers.NewUserHandler(userService, redisService, bcryptService)

				w := httptest.NewRecorder()
				c, _ := gin.CreateTestContext(w)
				c.Request, _ = http.NewRequest("POST", "/api/v1/reset-password", bytes.NewBufferString(tt.reqBody))
				c.Request.Header.Set("Content-Type", "application/json")
				c.Request.Header.Set("Authorization", "Bearer test-token")
				c.Set("UserID", uint(1))

				// Call the ResetPassword handler
				handler.ResetPassword(c)

				assert.Equal(t, http.StatusBadRequest, w.Code)

				var actualBody map[string]any
				err := json.Unmarshal(w.Body.Bytes(), &actualBody)
				require.NoError(t, err)

				assert.Equal(t, tt.expectedCode, actualBody["code"])
				assert.Equal(t, tt.expectedMsg, actualBody["message"])

			})
		}
	})

}
