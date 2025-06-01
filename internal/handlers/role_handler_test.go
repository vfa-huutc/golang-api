package handlers_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/vfa-khuongdv/golang-cms/internal/handlers"
	"github.com/vfa-khuongdv/golang-cms/internal/models"
	"github.com/vfa-khuongdv/golang-cms/pkg/errors"
	"github.com/vfa-khuongdv/golang-cms/tests/mocks"
)

func TestCreateRole(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Test cases
	testCases := []struct {
		name           string
		requestBody    map[string]interface{}
		mockSetup      func(*mocks.MockRoleService)
		expectedStatus int
		expectedBody   map[string]interface{}
	}{
		{
			name: "Success - Role Created",
			requestBody: map[string]interface{}{
				"name":         "admin",
				"display_name": "Administrator",
			},
			mockSetup: func(mrs *mocks.MockRoleService) {
				mrs.On("Create", mock.AnythingOfType("*models.Role")).Return(nil)
			},
			expectedStatus: http.StatusCreated,
			expectedBody: map[string]interface{}{
				"message": "Create new role successfully",
			},
		},
		{
			name: "Fail - Invalid Request Body",
			requestBody: map[string]interface{}{
				"name":         "ad", // too short, less than min:3
				"display_name": "Administrator",
			},
			mockSetup: func(mrs *mocks.MockRoleService) {
				// No mock setup needed as validation should fail
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody: map[string]interface{}{
				"code":    float64(errors.ErrInvalidData),
				"message": mock.Anything, // We don't test the exact message
			},
		},
		{
			name: "Fail - Service Error",
			requestBody: map[string]interface{}{
				"name":         "admin",
				"display_name": "Administrator",
			},
			mockSetup: func(mrs *mocks.MockRoleService) {
				mrs.On("Create", mock.AnythingOfType("*models.Role")).Return(
					errors.New(errors.ErrDBInsert, "Database error"),
				)
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody: map[string]interface{}{
				"code":    float64(errors.ErrDBInsert),
				"message": "Database error",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Setup mock
			mockService := new(mocks.MockRoleService)
			if tc.mockSetup != nil {
				tc.mockSetup(mockService)
			}

			// Create handler with mock service
			handler := handlers.NewRoleHandler(mockService)

			// Create a response recorder
			w := httptest.NewRecorder()

			// Create a request
			jsonValue, _ := json.Marshal(tc.requestBody)
			req, _ := http.NewRequest("POST", "/api/v1/roles", bytes.NewBuffer(jsonValue))
			req.Header.Set("Content-Type", "application/json")

			// Create a Gin context
			c, _ := gin.CreateTestContext(w)
			c.Request = req

			// Call the handler function
			handler.CreateRole(c)

			// Assert status code
			assert.Equal(t, tc.expectedStatus, w.Code)

			// Parse response body
			var response map[string]interface{}
			err := json.Unmarshal(w.Body.Bytes(), &response)
			assert.NoError(t, err)

			// Check if response matches expected structure and values
			if tc.expectedStatus == http.StatusCreated {
				// For success case
				assert.Equal(t, tc.expectedBody["message"], response["message"])
			} else {
				// For error case
				assert.Equal(t, tc.expectedBody["code"], response["code"])
				// If we care about the exact message
				if tc.expectedBody["message"] != mock.Anything {
					assert.Equal(t, tc.expectedBody["message"], response["message"])
				}
			}

			// Verify that all expected mock interactions occurred
			mockService.AssertExpectations(t)
		})
	}
}

func TestUpdateRole(t *testing.T) {
	gin.SetMode(gin.TestMode)

	testCases := []struct {
		name           string
		paramID        string
		requestBody    map[string]interface{}
		mockSetup      func(*mocks.MockRoleService)
		expectedStatus int
	}{
		{
			name:    "Success - Role Updated",
			paramID: "1",
			requestBody: map[string]interface{}{
				"display_name": "New Name",
			},
			mockSetup: func(m *mocks.MockRoleService) {
				m.On("GetByID", int64(1)).Return(&models.Role{ID: 1, Name: "admin"}, nil)
				m.On("Update", mock.AnythingOfType("*models.Role")).Return(nil)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:    "Fail - Invalid ID",
			paramID: "abc",
			requestBody: map[string]interface{}{
				"display_name": "New Name",
			},
			mockSetup:      nil,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:    "Fail - Validation Error",
			paramID: "1",
			requestBody: map[string]interface{}{
				"display_name": "x", // too short
			},
			mockSetup:      nil,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:    "Fail - GetByID Error",
			paramID: "1",
			requestBody: map[string]interface{}{
				"display_name": "Valid Name",
			},
			mockSetup: func(m *mocks.MockRoleService) {
				m.On("GetByID", int64(1)).Return(nil, errors.New(errors.ErrNotFound, "not found"))
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:    "Fail - Update Error",
			paramID: "1",
			requestBody: map[string]interface{}{
				"display_name": "Valid Name",
			},
			mockSetup: func(m *mocks.MockRoleService) {
				m.On("GetByID", int64(1)).Return(&models.Role{ID: 1}, nil)
				m.On("Update", mock.AnythingOfType("*models.Role")).Return(errors.New(errors.ErrDBUpdate, "update error"))
			},
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockService := new(mocks.MockRoleService)
			if tc.mockSetup != nil {
				tc.mockSetup(mockService)
			}

			handler := handlers.NewRoleHandler(mockService)
			w := httptest.NewRecorder()

			jsonValue, _ := json.Marshal(tc.requestBody)
			req, _ := http.NewRequest("PUT", "/api/v1/roles/"+tc.paramID, bytes.NewBuffer(jsonValue))
			req.Header.Set("Content-Type", "application/json")

			c, _ := gin.CreateTestContext(w)
			c.Params = gin.Params{{Key: "id", Value: tc.paramID}}
			c.Request = req

			handler.UpdateRole(c)
			assert.Equal(t, tc.expectedStatus, w.Code)
			mockService.AssertExpectations(t)
		})
	}
}

func TestGetRole(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockService := new(mocks.MockRoleService)
	handler := handlers.NewRoleHandler(mockService)

	t.Run("Success - Get Role", func(t *testing.T) {
		mockService.On("GetByID", int64(1)).Return(&models.Role{ID: 1, Name: "admin"}, nil)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/api/v1/roles/1", nil)

		c, _ := gin.CreateTestContext(w)
		c.Params = gin.Params{{Key: "id", Value: "1"}}
		c.Request = req

		handler.GetRole(c)

		assert.Equal(t, http.StatusOK, w.Code)
		mockService.AssertExpectations(t)
	})

	t.Run("Fail - Invalid ID", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/api/v1/roles/abc", nil)

		c, _ := gin.CreateTestContext(w)
		c.Params = gin.Params{{Key: "id", Value: "abc"}}
		c.Request = req

		handler.GetRole(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("Fail - Service Error", func(t *testing.T) {
		mockService.On("GetByID", int64(2)).Return(nil, errors.New(errors.ErrNotFound, "not found"))

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/api/v1/roles/2", nil)

		c, _ := gin.CreateTestContext(w)
		c.Params = gin.Params{{Key: "id", Value: "2"}}
		c.Request = req

		handler.GetRole(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		mockService.AssertExpectations(t)
	})
}

func TestDeleteRole(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockService := new(mocks.MockRoleService)
	handler := handlers.NewRoleHandler(mockService)

	t.Run("Success - Delete Role", func(t *testing.T) {
		mockService.On("Delete", int64(1)).Return(nil)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("DELETE", "/api/v1/roles/1", nil)

		c, _ := gin.CreateTestContext(w)
		c.Params = gin.Params{{Key: "id", Value: "1"}}
		c.Request = req

		handler.DeleteRole(c)

		assert.Equal(t, http.StatusOK, w.Code)
		mockService.AssertExpectations(t)
	})

	t.Run("Fail - Invalid ID", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("DELETE", "/api/v1/roles/abc", nil)

		c, _ := gin.CreateTestContext(w)
		c.Params = gin.Params{{Key: "id", Value: "abc"}}
		c.Request = req

		handler.DeleteRole(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("Fail - Service Error", func(t *testing.T) {
		mockService.On("Delete", int64(2)).Return(errors.New(errors.ErrNotFound, "not found"))

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("DELETE", "/api/v1/roles/2", nil)

		c, _ := gin.CreateTestContext(w)
		c.Params = gin.Params{{Key: "id", Value: "2"}}
		c.Request = req

		handler.DeleteRole(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		mockService.AssertExpectations(t)
	})
}
