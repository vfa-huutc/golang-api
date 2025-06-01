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
	"github.com/stretchr/testify/require"
	"github.com/vfa-khuongdv/golang-cms/internal/configs"
	"github.com/vfa-khuongdv/golang-cms/internal/handlers"
	"github.com/vfa-khuongdv/golang-cms/internal/services"
	appErrors "github.com/vfa-khuongdv/golang-cms/pkg/errors"
	"github.com/vfa-khuongdv/golang-cms/tests/mocks"
)

func TestLogin(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("Success", func(t *testing.T) {
		mockService := new(mocks.MockAuthService)

		// Create a new handler with the mock service
		handler := handlers.NewAuthHandler(mockService)

		mockService.On("Login", "email@gmail.com", "testpassword", mock.Anything).Return(
			&services.LoginResponse{
				AccessToken: configs.JwtResult{
					Token:     "testtoken",
					ExpiresAt: 0,
				},
				RefreshToken: configs.JwtResult{
					Token:     "testrefreshtoken",
					ExpiresAt: 0,
				},
			}, nil,
		)
		// Create a request with JSON body
		reqBody := `{"email":"email@gmail.com","password":"testpassword"}`
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest("POST", "/api/v1/login", bytes.NewBufferString(reqBody))
		c.Request.Header.Set("Content-Type", "application/json")

		handler.Login(c)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.JSONEq(t, `
		{
			"accessToken": {"token":"testtoken","expiresAt":0},
			"refreshToken": {"token":"testrefreshtoken","expiresAt":0}
		}
		`, w.Body.String())
		mockService.AssertExpectations(t)
	})

	t.Run("Error", func(t *testing.T) {
		mockService := new(mocks.MockAuthService)

		// Create a new handler with the mock service
		handler := handlers.NewAuthHandler(mockService)

		// Simulate an error in the service
		mockService.On("Login", "email@gmail.com", "testpassword", mock.Anything).Return(nil, appErrors.New(appErrors.ErrInvalidPassword, "Invalid email or password"))

		// Create a request with JSON body
		reqBody := `{"email":"email@gmail.com","password":"testpassword"}`
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest("POST", "/api/v1/login", bytes.NewBufferString(reqBody))
		c.Request.Header.Set("Content-Type", "application/json")

		// Simulate an error in the service
		handler.Login(c)
		assert.Equal(t, http.StatusUnauthorized, w.Code)
		expectedBody := map[string]any{
			"code":    float64(appErrors.ErrInvalidPassword),
			"message": "Invalid email or password",
		}

		var actualBody map[string]any
		err := json.Unmarshal(w.Body.Bytes(), &actualBody)
		require.NoError(t, err)
		assert.Equal(t, expectedBody["code"], actualBody["code"])
		assert.Equal(t, expectedBody["message"], actualBody["message"])
		mockService.AssertExpectations(t)
	})

	t.Run("Invalid input cases", func(t *testing.T) {

		mockService := new(mocks.MockAuthService)
		handler := handlers.NewAuthHandler(mockService)

		tests := []struct {
			name         string
			reqBody      string
			expectedCode float64
			expectedMsg  string
		}{
			{
				name:         "MissingEmailAndPassword",
				reqBody:      `{}`,
				expectedCode: float64(4001),
				expectedMsg: "Key: 'Email' Error:Field validation for 'Email' failed on the 'required' tag\n" +
					"Key: 'Password' Error:Field validation for 'Password' failed on the 'required' tag",
			},
			{
				name:         "MissingEmail",
				reqBody:      `{"password":"validPassword123"}`,
				expectedCode: float64(4001),
				expectedMsg:  "Key: 'Email' Error:Field validation for 'Email' failed on the 'required' tag",
			},
			{
				name:         "InvalidEmailFormat",
				reqBody:      `{"email":"not-an-email","password":"validPassword123"}`,
				expectedCode: float64(4001),
				expectedMsg:  "Key: 'Email' Error:Field validation for 'Email' failed on the 'email' tag",
			},
			{
				name:         "EmptyEmail",
				reqBody:      `{"email":"","password":"validPassword123"}`,
				expectedCode: float64(4001),
				expectedMsg:  "Key: 'Email' Error:Field validation for 'Email' failed on the 'required' tag",
			},
			{
				name:         "PasswordTooShort",
				reqBody:      `{"email":"user@example.com","password":"123"}`,
				expectedCode: float64(4001),
				expectedMsg:  "Key: 'Password' Error:Field validation for 'Password' failed on the 'min' tag",
			},
			{
				name:         "PasswordTooLong",
				reqBody:      `{"email":"user@example.com","password":"` + strings.Repeat("a", 256) + `"}`,
				expectedCode: float64(4001),
				expectedMsg:  "Key: 'Password' Error:Field validation for 'Password' failed on the 'max' tag",
			},
			{
				name:         "EmptyPassword",
				reqBody:      `{"email":"user@example.com","password":""}`,
				expectedCode: float64(4001),
				expectedMsg:  "Key: 'Password' Error:Field validation for 'Password' failed on the 'required' tag",
			},
			{
				name:         "MissingPassword",
				reqBody:      `{"email":"user@example.com"}`,
				expectedCode: float64(4001),
				expectedMsg:  "Key: 'Password' Error:Field validation for 'Password' failed on the 'required' tag",
			},
		}

		for _, tc := range tests {
			t.Run(tc.name, func(t *testing.T) {
				w := httptest.NewRecorder()
				c, _ := gin.CreateTestContext(w)
				c.Request, _ = http.NewRequest("POST", "/api/v1/login", bytes.NewBufferString(tc.reqBody))
				c.Request.Header.Set("Content-Type", "application/json")

				handler.Login(c)

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
}

func TestRefreshToken(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("Success", func(t *testing.T) {
		mockService := new(mocks.MockAuthService)

		// Create a new handler with the mock service
		handler := handlers.NewAuthHandler(mockService)

		mockService.On("RefreshToken", "testrefreshtoken", mock.Anything).Return(
			&services.LoginResponse{
				AccessToken: configs.JwtResult{
					Token:     "newtesttoken",
					ExpiresAt: 0,
				},
				RefreshToken: configs.JwtResult{
					Token:     "newtestrefreshtoken",
					ExpiresAt: 0,
				},
			}, nil,
		)

		reqBody := `{"refresh_token":"testrefreshtoken"}`
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest("POST", "/api/v1/refresh-token", bytes.NewBufferString(reqBody))
		c.Request.Header.Set("Content-Type", "application/json")

		handler.RefreshToken(c)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.JSONEq(t, `
		{
			"accessToken": {"token":"newtesttoken","expiresAt":0},
			"refreshToken": {"token":"newtestrefreshtoken","expiresAt":0}
		}
		`, w.Body.String())

		mockService.AssertExpectations(t)
	})

	t.Run("Error", func(t *testing.T) {
		mockService := new(mocks.MockAuthService)

		// Create a new handler with the mock service
		handler := handlers.NewAuthHandler(mockService)

		// Simulate an error in the service
		mockService.On("RefreshToken", "invalidtoken", mock.Anything).Return(nil, appErrors.New(appErrors.ErrUnauthorized, "Invalid refresh token"))

		reqBody := `{"refresh_token":"invalidtoken"}`
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest("POST", "/api/v1/refresh-token", bytes.NewBufferString(reqBody))
		c.Request.Header.Set("Content-Type", "application/json")

		handler.RefreshToken(c)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
		expectedBody := map[string]any{
			"code":    float64(appErrors.ErrUnauthorized),
			"message": "Invalid refresh token",
		}

		var actualBody map[string]any
		err := json.Unmarshal(w.Body.Bytes(), &actualBody)
		require.NoError(t, err)
		assert.Equal(t, expectedBody["code"], actualBody["code"])
		assert.Equal(t, expectedBody["message"], actualBody["message"])

		mockService.AssertExpectations(t)
	})

	t.Run("Invalid input cases", func(t *testing.T) {

		mockService := new(mocks.MockAuthService)
		handler := handlers.NewAuthHandler(mockService)

		tests := []struct {
			name         string
			reqBody      string
			expectedCode float64
			expectedMsg  string
		}{
			{
				name:         "MissingRefreshToken",
				reqBody:      `{}`,
				expectedCode: float64(4001),
				expectedMsg:  "Key: 'RefreshToken' Error:Field validation for 'RefreshToken' failed on the 'required' tag",
			},
			{
				name:         "EmptyRefreshToken",
				reqBody:      `{"refresh_token":""}`,
				expectedCode: float64(4001),
				expectedMsg:  "Key: 'RefreshToken' Error:Field validation for 'RefreshToken' failed on the 'required' tag",
			},
		}

		for _, tc := range tests {
			t.Run(tc.name, func(t *testing.T) {
				w := httptest.NewRecorder()
				c, _ := gin.CreateTestContext(w)
				c.Request, _ = http.NewRequest("POST", "/api/v1/refresh-token", bytes.NewBufferString(tc.reqBody))
				c.Request.Header.Set("Content-Type", "application/json")

				handler.RefreshToken(c)

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

}
