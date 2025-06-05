package handlers_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/vfa-khuongdv/golang-cms/internal/handlers"
	"github.com/vfa-khuongdv/golang-cms/internal/services"
	"github.com/vfa-khuongdv/golang-cms/internal/utils"
	"github.com/vfa-khuongdv/golang-cms/pkg/apperror"
	"github.com/vfa-khuongdv/golang-cms/tests/mocks"
)

func TestLogin(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("Login - Success", func(t *testing.T) {
		mockService := new(mocks.MockAuthService)
		handler := handlers.NewAuthHandler(mockService)

		// Mock the service method
		mockService.On("Login", "email@gmail.com", "testpassword", mock.Anything).Return(
			&services.LoginResponse{
				AccessToken: services.JwtResult{
					Token:     "testtoken",
					ExpiresAt: 0,
				},
				RefreshToken: services.JwtResult{
					Token:     "testrefreshtoken",
					ExpiresAt: 0,
				},
			}, nil,
		)

		requestBody := map[string]string{
			"email":    "email@gmail.com",
			"password": "testpassword",
		}

		reqBody, _ := json.Marshal(requestBody)

		// Create a test context
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest("POST", "/api/v1/login", bytes.NewBuffer(reqBody))
		c.Request.Header.Set("Content-Type", "application/json")

		// Call the handler
		handler.Login(c)

		// Assert the response
		assert.Equal(t, http.StatusOK, w.Code)
		assert.JSONEq(t, `
		{
			"accessToken": {"token":"testtoken","expiresAt":0},
			"refreshToken": {"token":"testrefreshtoken","expiresAt":0}
		}
		`, w.Body.String())
		// Assert that the mock service method was called
		mockService.AssertExpectations(t)
	})

	t.Run("Login - Create Error", func(t *testing.T) {
		mockService := new(mocks.MockAuthService)
		handler := handlers.NewAuthHandler(mockService)

		// Mock the service method
		mockService.On("Login", "email@gmail.com", "testpassword", mock.Anything).Return(nil, apperror.NewUnauthorizedError("Invalid email or password"))

		requestBody := map[string]string{
			"email":    "email@gmail.com",
			"password": "testpassword",
		}
		reqBody, _ := json.Marshal(requestBody)

		// Create a test context
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest("POST", "/api/v1/login", bytes.NewBuffer(reqBody))
		c.Request.Header.Set("Content-Type", "application/json")

		// Call the handler
		handler.Login(c)

		// Assert the response
		expectedBody := map[string]any{
			"code":    float64(apperror.ErrUnauthorized),
			"message": "Invalid email or password",
		}
		var actualBody map[string]any
		json.Unmarshal(w.Body.Bytes(), &actualBody)
		assert.Equal(t, expectedBody["code"], actualBody["code"])
		assert.Equal(t, expectedBody["message"], actualBody["message"])

		// Assert mocks
		mockService.AssertExpectations(t)
	})

	t.Run("Login - Validation Errors", func(t *testing.T) {

		// Create a mock service and handler
		mockService := new(mocks.MockAuthService)
		handler := handlers.NewAuthHandler(mockService)

		tests := []struct {
			name           string
			reqBody        string
			expectedCode   float64
			expectedMsg    string
			expectedFields []apperror.FieldError
		}{
			{
				name:         "MissingEmailAndPassword",
				reqBody:      `{}`,
				expectedCode: float64(4001),
				expectedMsg:  "Validation failed",
				expectedFields: []apperror.FieldError{
					{Field: "email", Message: "email is required"},
					{Field: "password", Message: "password is required"},
				},
			},
			{
				name:         "MissingEmail",
				reqBody:      `{"password":"validPassword123"}`,
				expectedCode: float64(4001),
				expectedMsg:  "Validation failed",
				expectedFields: []apperror.FieldError{
					{Field: "email", Message: "email is required"},
				},
			},
			{
				name:         "InvalidEmailFormat",
				reqBody:      `{"email":"not-an-email","password":"validPassword123"}`,
				expectedCode: float64(4001),
				expectedMsg:  "Validation failed",
				expectedFields: []apperror.FieldError{
					{Field: "email", Message: "email must be a valid email address"},
				},
			},
			{
				name:         "EmptyEmail",
				reqBody:      `{"email":"","password":"validPassword123"}`,
				expectedCode: float64(4001),
				expectedMsg:  "Validation failed",
				expectedFields: []apperror.FieldError{
					{Field: "email", Message: "email is required"},
				},
			},
			{
				name:         "PasswordTooShort",
				reqBody:      `{"email":"user@example.com","password":"123"}`,
				expectedCode: float64(4001),
				expectedMsg:  "Validation failed",
				expectedFields: []apperror.FieldError{
					{Field: "password", Message: "password must be at least 6 characters long or numeric"},
				},
			},
			{
				name:         "PasswordTooLong",
				reqBody:      `{"email":"user@example.com","password":"` + strings.Repeat("a", 256) + `"}`,
				expectedCode: float64(4001),
				expectedMsg:  "Validation failed",
				expectedFields: []apperror.FieldError{
					{Field: "password", Message: "password must be at most 255 characters long or numeric"},
				},
			},
			{
				name:         "EmptyPassword",
				reqBody:      `{"email":"user@example.com","password":""}`,
				expectedCode: float64(4001),
				expectedMsg:  "Validation failed",
				expectedFields: []apperror.FieldError{
					{Field: "password", Message: "password is required"},
				},
			},
			{
				name:         "MissingPassword",
				reqBody:      `{"email":"user@example.com"}`,
				expectedCode: float64(4001),
				expectedMsg:  "Validation failed",
				expectedFields: []apperror.FieldError{
					{Field: "password", Message: "password is required"},
				},
			},
		}

		for _, tc := range tests {
			t.Run(tc.name, func(t *testing.T) {
				// Create a test context
				w := httptest.NewRecorder()
				c, _ := gin.CreateTestContext(w)
				c.Request, _ = http.NewRequest("POST", "/api/v1/login", bytes.NewBufferString(tc.reqBody))

				// Call the handler method
				handler.Login(c)

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
				assert.Equal(t, expectedBody["fields"], utils.MapJsonToFieldErrors(actualBody["fields"]))

				// Assert mocks
				mockService.AssertExpectations(t)

			})
		}
	})
}

func TestRefreshToken(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("RefreshToken - Success", func(t *testing.T) {
		mockService := new(mocks.MockAuthService)
		handler := handlers.NewAuthHandler(mockService)

		// Mock the service method
		mockService.On("RefreshToken", "testrefreshtoken", mock.Anything).Return(
			&services.LoginResponse{
				AccessToken: services.JwtResult{
					Token:     "newtesttoken",
					ExpiresAt: 0,
				},
				RefreshToken: services.JwtResult{
					Token:     "newtestrefreshtoken",
					ExpiresAt: 0,
				},
			}, nil,
		)
		requestBody := map[string]string{
			"refresh_token": "testrefreshtoken",
		}
		reqBody, _ := json.Marshal(requestBody)

		// Create a test context
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest("POST", "/api/v1/refresh-token", bytes.NewBuffer(reqBody))
		c.Request.Header.Set("Content-Type", "application/json")

		// Call the handler
		handler.RefreshToken(c)

		// Assert the response
		assert.Equal(t, http.StatusOK, w.Code)
		assert.JSONEq(t, `
		{
			"accessToken": {"token":"newtesttoken","expiresAt":0},
			"refreshToken": {"token":"newtestrefreshtoken","expiresAt":0}
		}
		`, w.Body.String())

		// Assert mocks
		mockService.AssertExpectations(t)
	})

	t.Run("RefreshToken - Error Invalid Token", func(t *testing.T) {
		mockService := new(mocks.MockAuthService)
		handler := handlers.NewAuthHandler(mockService)

		// Mock the service method
		mockService.On("RefreshToken", "invalidtoken", mock.Anything).Return(nil, apperror.NewUnauthorizedError("Invalid refresh token"))
		reqBody := map[string]string{
			"refresh_token": "invalidtoken",
		}
		reqBodyBytes, _ := json.Marshal(reqBody)

		// Create a test context
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest("POST", "/api/v1/refresh-token", bytes.NewBuffer(reqBodyBytes))
		c.Request.Header.Set("Content-Type", "application/json")

		// Call the handler
		handler.RefreshToken(c)

		// Assert the response
		assert.Equal(t, http.StatusUnauthorized, w.Code)
		expectedBody := map[string]any{
			"code":    float64(apperror.ErrUnauthorized),
			"message": "Invalid refresh token",
		}

		var actualBody map[string]any
		json.Unmarshal(w.Body.Bytes(), &actualBody)
		assert.Equal(t, expectedBody["code"], actualBody["code"])
		assert.Equal(t, expectedBody["message"], actualBody["message"])
		assert.Equal(t, expectedBody["code"], actualBody["code"])
		assert.Equal(t, expectedBody["message"], actualBody["message"])

		// Assert mocks
		mockService.AssertExpectations(t)
	})

	t.Run("RefreshToken - Validation Errors", func(t *testing.T) {
		mockService := new(mocks.MockAuthService)
		handler := handlers.NewAuthHandler(mockService)

		tests := []struct {
			name           string
			reqBody        string
			expectedCode   float64
			expectedMsg    string
			expectedFields []apperror.FieldError
		}{
			{
				name:         "MissingRefreshToken",
				reqBody:      `{}`,
				expectedCode: float64(4001),
				expectedMsg:  "Validation failed",
				expectedFields: []apperror.FieldError{
					{Field: "refresh_token", Message: "refresh_token is required"},
				},
			},
			{
				name:         "EmptyRefreshToken",
				reqBody:      `{"refresh_token":""}`,
				expectedCode: float64(4001),
				expectedMsg:  "Validation failed",
				expectedFields: []apperror.FieldError{
					{Field: "refresh_token", Message: "refresh_token is required"},
				},
			},
		}

		for _, tc := range tests {
			t.Run(tc.name, func(t *testing.T) {
				// Create a test context
				w := httptest.NewRecorder()
				c, _ := gin.CreateTestContext(w)
				c.Request, _ = http.NewRequest("POST", "/api/v1/refresh-token", bytes.NewBufferString(tc.reqBody))

				// Call the handler method
				handler.RefreshToken(c)

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
				assert.Equal(t, expectedBody["fields"], utils.MapJsonToFieldErrors(actualBody["fields"]))

				// Assert mocks
				mockService.AssertExpectations(t)
			})
		}
	})

}
