package tests_internal_handlers

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
	"github.com/vfa-khuongdv/golang-cms/internal/services"
	"github.com/vfa-khuongdv/golang-cms/pkg/errors"
)

// MockRoleService is a mock implementation of the role service
type MockRoleService struct {
	mock.Mock
	services.IRoleService
}

// Create mocks the Create method of the role service interface.
// It records the call with the provided role model and returns the error
// configured in the test setup.
func (m *MockRoleService) Create(role *models.Role) error {
	args := m.Called(role)
	return args.Error(0)
}

// GetByID mocks the retrieval of a role by its ID.
// It records the call with the provided ID and returns the configured role model and error.
func (m *MockRoleService) GetByID(id int64) (*models.Role, error) {
	args := m.Called(id)
	var role *models.Role
	if r := args.Get(0); r != nil {
		role = r.(*models.Role)
	}
	return role, args.Error(1)
}

// Update mocks the Update method of the role service interface.
// It records the call with the provided role model and returns the error
// configured in the test setup.
func (m *MockRoleService) Update(role *models.Role) error {
	args := m.Called(role)
	return args.Error(0)
}

// Delete mocks the Delete method of the role service interface.
// It records the call with the provided role ID and returns the error
// configured in the test setup.
func (m *MockRoleService) Delete(id int64) error {
	args := m.Called(id)
	return args.Error(0)
}

// AssignPermissions mocks the AssignPermissions method of the role service interface.
// It records the call with the provided role ID and permission IDs array, then returns
// the error configured in the test setup.
func (m *MockRoleService) AssignPermissions(roleID uint, permissionIDs []uint) error {
	args := m.Called(roleID, permissionIDs)
	return args.Error(0)
}

// GetRolePermissions mocks the GetRolePermissions method of the role service interface.
// It records the call with the provided role ID and returns the configured list of
// permission models and error from the test setup.
func (m *MockRoleService) GetRolePermissions(roleID uint) ([]models.Permission, error) {
	args := m.Called(roleID)
	return args.Get(0).([]models.Permission), args.Error(1)
}

func TestCreateRole(t *testing.T) {
	// Set Gin to Test Mode
	gin.SetMode(gin.TestMode)

	// Test cases
	testCases := []struct {
		name           string
		requestBody    map[string]interface{}
		mockSetup      func(*MockRoleService)
		expectedStatus int
		expectedBody   map[string]interface{}
	}{
		{
			name: "Success - Role Created",
			requestBody: map[string]interface{}{
				"name":         "admin",
				"display_name": "Administrator",
			},
			mockSetup: func(mrs *MockRoleService) {
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
			mockSetup: func(mrs *MockRoleService) {
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
			mockSetup: func(mrs *MockRoleService) {
				mrs.On("Create", mock.AnythingOfType("*models.Role")).Return(
					errors.New(errors.ErrDatabaseInsert, "Database error"),
				)
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody: map[string]interface{}{
				"code":    float64(errors.ErrDatabaseInsert),
				"message": "Database error",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Setup mock
			mockService := new(MockRoleService)
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
