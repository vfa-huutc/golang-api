package handlers_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/vfa-khuongdv/golang-cms/internal/handlers"
	"github.com/vfa-khuongdv/golang-cms/internal/models"
	"github.com/vfa-khuongdv/golang-cms/internal/utils"
	"github.com/vfa-khuongdv/golang-cms/pkg/apperror"
	"github.com/vfa-khuongdv/golang-cms/tests/mocks"
)

func TestCreateRole(t *testing.T) {
	gin.SetMode(gin.TestMode)

	var requestBody = map[string]any{
		"name":         "admin",
		"display_name": "Administrator",
	}
	role := models.Role{
		Name:        "admin",
		DisplayName: "Administrator",
	}

	t.Run("CreateRole - Success", func(t *testing.T) {
		mockService := new(mocks.MockRoleService)
		handler := handlers.NewRoleHandler(mockService)

		jsonValue, _ := json.Marshal(requestBody)
		// Mock service methods
		mockService.On("Create", &role).Return(nil)

		// Create a test context
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/api/v1/roles", bytes.NewBuffer(jsonValue))
		c, _ := gin.CreateTestContext(w)
		c.Request = req

		// Call the handler function
		handler.CreateRole(c)

		// Assert the response status code and body
		assert.Equal(t, http.StatusCreated, w.Code)
		var response map[string]any
		json.Unmarshal(w.Body.Bytes(), &response)
		assert.Equal(t, "Create new role successfully", response["message"])

		// Assert mocks
		mockService.AssertExpectations(t)
	})

	t.Run("CreateRole - Failed To Create", func(t *testing.T) {
		mockService := new(mocks.MockRoleService)
		handler := handlers.NewRoleHandler(mockService)

		// Mock service methods
		mockService.On("Create", &role).Return(apperror.NewDBInsertError("Database error"))

		// Create a test context
		w := httptest.NewRecorder()
		jsonValue, _ := json.Marshal(requestBody)
		req, _ := http.NewRequest("POST", "/api/v1/roles", bytes.NewBuffer(jsonValue))
		c, _ := gin.CreateTestContext(w)
		c.Request = req

		// Call the handler function
		handler.CreateRole(c)

		// Assert the response status code and body
		assert.Equal(t, http.StatusInternalServerError, w.Code)
		var response map[string]any
		json.Unmarshal(w.Body.Bytes(), &response)
		assert.Equal(t, float64(apperror.ErrDBInsert), response["code"])
		assert.Equal(t, "Database error", response["message"])

		// Assert mocks
		mockService.AssertExpectations(t)

	})

	t.Run("CreateRole Validation Error", func(t *testing.T) {
		tests := []struct {
			name           string
			reqBody        string
			expectedCode   float64
			expectedMsg    string
			expectedFields []apperror.FieldError
		}{
			{
				name:         "NameRequired",
				reqBody:      `{"display_name": "Administrator"}`,
				expectedCode: float64(apperror.ErrValidationFailed),
				expectedMsg:  "Validation failed",
				expectedFields: []apperror.FieldError{
					{Field: "name", Message: "name is required"},
				},
			},
			{
				name:         "NameTooShort",
				reqBody:      `{"name": "ad", "display_name": "Administrator"}`,
				expectedCode: float64(apperror.ErrValidationFailed),
				expectedMsg:  "Validation failed",
				expectedFields: []apperror.FieldError{
					{Field: "name", Message: "name must be at least 3 characters long or numeric"},
				},
			},
			{
				name:         "NameTooLong",
				reqBody:      `{"name": "` + strings.Repeat("a", 256) + `", "display_name": "Administrator"}`,
				expectedCode: float64(apperror.ErrValidationFailed),
				expectedMsg:  "Validation failed",
				expectedFields: []apperror.FieldError{
					{Field: "name", Message: "name must be at most 255 characters long or numeric"},
				},
			},
			{
				name:         "DisplayNameRequired",
				reqBody:      `{"name": "admin"}`,
				expectedCode: float64(apperror.ErrValidationFailed),
				expectedMsg:  "Validation failed",
				expectedFields: []apperror.FieldError{
					{Field: "display_name", Message: "display_name is required"},
				},
			},
			{
				name:         "DisplayNameTooShort",
				reqBody:      `{"name": "admin", "display_name": "Ad"}`,
				expectedCode: float64(apperror.ErrValidationFailed),
				expectedMsg:  "Validation failed",
				expectedFields: []apperror.FieldError{
					{Field: "display_name", Message: "display_name must be at least 3 characters long or numeric"},
				},
			},
			{
				name:         "DisplayNameTooLong",
				reqBody:      `{"name": "admin", "display_name": "` + strings.Repeat("a", 256) + `"}`,
				expectedCode: float64(apperror.ErrValidationFailed),
				expectedMsg:  "Validation failed",
				expectedFields: []apperror.FieldError{
					{Field: "display_name", Message: "display_name must be at most 255 characters long or numeric"},
				},
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				mockService := new(mocks.MockRoleService)
				handler := handlers.NewRoleHandler(mockService)

				// Create a test context
				w := httptest.NewRecorder()
				c, _ := gin.CreateTestContext(w)
				req, _ := http.NewRequest("POST", "/api/v1/roles", bytes.NewBufferString(tt.reqBody))
				req.Header.Set("Content-Type", "application/json")
				c.Request = req

				// Call the handler function
				handler.CreateRole(c)

				// Assert the response status code and body
				var response map[string]any
				json.Unmarshal(w.Body.Bytes(), &response)
				assert.Equal(t, http.StatusBadRequest, w.Code)
				assert.Equal(t, tt.expectedCode, response["code"])
				assert.Equal(t, tt.expectedMsg, response["message"])
				assert.Equal(t, tt.expectedFields, utils.MapJsonToFieldErrors(response["fields"]))

				// Assert mocks
				mockService.AssertExpectations(t)
			})
		}
	})

}

func TestUpdateRole(t *testing.T) {
	gin.SetMode(gin.TestMode)

	var requestBody = map[string]any{
		"display_name": "Super Administrator",
	}

	t.Run("UpdateRole Success", func(t *testing.T) {
		mockService := new(mocks.MockRoleService)
		handler := handlers.NewRoleHandler(mockService)

		role := models.Role{
			ID:          1,
			Name:        "admin",
			DisplayName: "Administrator",
		}

		// Mock service method
		mockService.On("GetByID", int64(1)).Return(&role, nil)
		mockService.On("Update", &role).Return(nil)

		jsonValue, _ := json.Marshal(requestBody)

		// Create a test context
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		req, _ := http.NewRequest("PUT", "/api/v1/roles/1", bytes.NewBuffer(jsonValue))
		req.Header.Set("Content-Type", "application/json")
		c.Params = gin.Params{{Key: "id", Value: "1"}}
		c.Request = req

		// Call the handler function
		handler.UpdateRole(c)

		// Assert the response status code and body
		assert.Equal(t, http.StatusOK, w.Code)
		var response map[string]any
		json.Unmarshal(w.Body.Bytes(), &response)
		assert.Equal(t, "Update role successfully", response["message"])

		// Assert mocks
		mockService.AssertExpectations(t)
	})

	t.Run("UpdateRole Error - GetByID Error", func(t *testing.T) {
		mockService := new(mocks.MockRoleService)
		handler := handlers.NewRoleHandler(mockService)

		// Mock service method
		mockService.On("GetByID", int64(2)).Return(nil, apperror.NewNotFoundError("Role not found"))

		// Create a test context
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		jsonValue, _ := json.Marshal(requestBody)
		c.Request, _ = http.NewRequest("PUT", "/api/v1/roles/:id", bytes.NewBuffer(jsonValue))
		c.Params = gin.Params{{Key: "id", Value: "2"}}

		// Call the handler function
		handler.UpdateRole(c)

		// Assert the response status code and body
		var response map[string]any
		json.Unmarshal(w.Body.Bytes(), &response)
		assert.Equal(t, http.StatusNotFound, w.Code)
		assert.Equal(t, float64(apperror.ErrNotFound), response["code"])
		assert.Equal(t, "Role not found", response["message"])

		// Assert mocks
		mockService.AssertExpectations(t)
	})

	t.Run("Parse params error", func(t *testing.T) {
		mockService := new(mocks.MockRoleService)
		handler := handlers.NewRoleHandler(mockService)

		jsonValue, _ := json.Marshal(requestBody)

		// Create a test context
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest("PUT", "/api/v1/roles/:id", bytes.NewBuffer(jsonValue))
		c.Params = gin.Params{{Key: "id", Value: "abc"}}

		// Call the handler function
		handler.UpdateRole(c)

		// Assert the response status code and body
		var response map[string]any
		json.Unmarshal(w.Body.Bytes(), &response)
		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Equal(t, float64(apperror.ErrParseError), response["code"])
		assert.Equal(t, "Invalid RoleID", response["message"])

		// Assert mocks
		mockService.AssertExpectations(t)
	})

	t.Run("UpdateRole - Update Error", func(t *testing.T) {
		mockService := new(mocks.MockRoleService)
		handler := handlers.NewRoleHandler(mockService)

		role := models.Role{
			ID:          1,
			Name:        "admin",
			DisplayName: "Administrator",
		}
		// Mock service method
		mockService.On("GetByID", int64(1)).Return(&role, nil)
		mockService.On("Update", &role).Return(apperror.NewDBUpdateError("Database error"))

		// Create a test context
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		jsonValue, _ := json.Marshal(requestBody)

		c.Request, _ = http.NewRequest("PUT", "/api/v1/roles/:id", bytes.NewBuffer(jsonValue))
		c.Params = gin.Params{{Key: "id", Value: "1"}}

		// Call the handler function
		handler.UpdateRole(c)

		// Assert the response status code and body
		var response map[string]any
		json.Unmarshal(w.Body.Bytes(), &response)
		assert.Equal(t, float64(apperror.ErrDBUpdate), response["code"])
		assert.Equal(t, "Database error", response["message"])
		assert.Equal(t, http.StatusInternalServerError, w.Code)

		// Assert mocks
		mockService.AssertExpectations(t)
	})

	t.Run("UpdateRole - Validation Error", func(t *testing.T) {
		tests := []struct {
			name           string
			reqBody        string
			expectedCode   float64
			expectedMsg    string
			expectedFields []apperror.FieldError
		}{
			{
				name:         "DisplayNameRequired",
				reqBody:      `{}`,
				expectedCode: float64(apperror.ErrValidationFailed),
				expectedMsg:  "Validation failed",
				expectedFields: []apperror.FieldError{
					{Field: "display_name", Message: "display_name is required"},
				},
			},
			{
				name:         "DisplayNameTooShort",
				reqBody:      `{"display_name": "Ad"}`,
				expectedCode: float64(apperror.ErrValidationFailed),
				expectedMsg:  "Validation failed",
				expectedFields: []apperror.FieldError{
					{Field: "display_name", Message: "display_name must be at least 3 characters long or numeric"},
				},
			},
			{
				name:         "DisplayNameTooLong",
				reqBody:      `{"display_name": "` + strings.Repeat("a", 256) + `"}`,
				expectedCode: float64(apperror.ErrValidationFailed),
				expectedMsg:  "Validation failed",
				expectedFields: []apperror.FieldError{
					{Field: "display_name", Message: "display_name must be at most 255 characters long or numeric"},
				},
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				mockService := new(mocks.MockRoleService)
				handler := handlers.NewRoleHandler(mockService)

				// Create a test context
				w := httptest.NewRecorder()
				c, _ := gin.CreateTestContext(w)
				c.Request, _ = http.NewRequest("PUT", "/api/v1/roles/:id", bytes.NewBufferString(tt.reqBody))
				c.Request.Header.Set("Content-Type", "application/json")
				c.Params = gin.Params{{Key: "id", Value: "1"}}

				// Call the handler function
				handler.UpdateRole(c)

				// Assert the response status code and body
				var response map[string]any
				json.Unmarshal(w.Body.Bytes(), &response)

				assert.Equal(t, http.StatusBadRequest, w.Code)
				assert.Equal(t, tt.expectedCode, response["code"])
				assert.Equal(t, tt.expectedMsg, response["message"])
				assert.Equal(t, tt.expectedFields, utils.MapJsonToFieldErrors(response["fields"]))

				// Assert mocks
				mockService.AssertExpectations(t)
			})
		}
	})

}

func TestGetRole(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("GetRole - Success", func(t *testing.T) {
		mockService := new(mocks.MockRoleService)
		handler := handlers.NewRoleHandler(mockService)

		role := &models.Role{
			ID:          1,
			Name:        "admin",
			DisplayName: "Administrator",
			CreatedAt:   time.Date(2023, 10, 1, 0, 0, 0, 0, time.UTC),
			UpdatedAt:   time.Date(2023, 10, 2, 0, 0, 0, 0, time.UTC),
		}
		// Mock the service method
		mockService.On("GetByID", int64(1)).Return(role, nil)

		// Create a test context
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = gin.Params{{Key: "id", Value: "1"}}
		c.Request, _ = http.NewRequest("GET", "/api/v1/roles/:id", nil)

		// Call the handler function
		handler.GetRole(c)

		// Assert the response status code and body
		assert.JSONEq(t, `{
		"id": 1,
		"name": "admin",
		"displayName": "Administrator",
		"createdAt": "2023-10-01T00:00:00Z",
		"updatedAt": "2023-10-02T00:00:00Z",
		"deletedAt": null
	}`, w.Body.String())

		assert.Equal(t, http.StatusOK, w.Code)

		// Assert mocks
		mockService.AssertExpectations(t)
	})

	t.Run("GetRole - Error Invalid ID", func(t *testing.T) {
		// Mock service method
		mockService := new(mocks.MockRoleService)
		handler := handlers.NewRoleHandler(mockService)

		// Create a test context
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = gin.Params{{Key: "id", Value: "abc"}}
		c.Request, _ = http.NewRequest("GET", "/api/v1/roles/:id", nil)

		// Call the handler function
		handler.GetRole(c)

		// Assert the response status code and body
		var response map[string]any
		json.Unmarshal(w.Body.Bytes(), &response)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Equal(t, float64(apperror.ErrParseError), response["code"])
		assert.Equal(t, "Invalid RoleID", response["message"])

		// Assert mocks
		mockService.AssertExpectations(t)
	})

	t.Run("GetRole - Error GetByID NotFound", func(t *testing.T) {
		mockService := new(mocks.MockRoleService)
		handler := handlers.NewRoleHandler(mockService)

		// Mock service methods
		mockService.On("GetByID", int64(2)).Return(nil, apperror.NewNotFoundError("Role not found"))

		// Create a test context
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = gin.Params{{Key: "id", Value: "2"}}
		c.Request, _ = http.NewRequest("GET", "/api/v1/roles/:id", nil)

		// Call the handler function
		handler.GetRole(c)

		// Assert the response status code and body
		assert.Equal(t, http.StatusNotFound, w.Code)
		var response map[string]any
		json.Unmarshal(w.Body.Bytes(), &response)
		assert.Equal(t, float64(apperror.ErrNotFound), response["code"])

		// Assert mocks
		mockService.AssertExpectations(t)
	})
}

func TestDeleteRole(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("DeleteRole - Success", func(t *testing.T) {
		mockService := new(mocks.MockRoleService)
		handler := handlers.NewRoleHandler(mockService)

		// Mock service method
		mockService.On("Delete", int64(1)).Return(nil)

		// Create a test context
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest("DELETE", "/api/v1/roles/:id", nil)
		c.Params = gin.Params{{Key: "id", Value: "1"}}

		// Call the handler function
		handler.DeleteRole(c)

		// Assert the response
		var response map[string]any
		json.Unmarshal(w.Body.Bytes(), &response)
		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, "Delete role successfully", response["message"])

		// Assert mocks
		mockService.AssertExpectations(t)
	})

	t.Run("DeleteRole - Error Invalid ID", func(t *testing.T) {
		mockService := new(mocks.MockRoleService)
		handler := handlers.NewRoleHandler(mockService)

		// Create a test context
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest("DELETE", "/api/v1/roles/:id", nil)
		c.Params = gin.Params{{Key: "id", Value: "abc"}}

		// Call the handler function
		handler.DeleteRole(c)

		// Assert the response
		var response map[string]any
		json.Unmarshal(w.Body.Bytes(), &response)
		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Equal(t, float64(apperror.ErrParseError), response["code"])
		assert.Equal(t, "Invalid RoleID", response["message"])

		// Assert mocks
		mockService.AssertExpectations(t)
	})

	t.Run("DeleteRole - Error Delete", func(t *testing.T) {
		mockService := new(mocks.MockRoleService)
		handler := handlers.NewRoleHandler(mockService)

		err := apperror.NewDBDeleteError("Database error")
		// Mock service method
		mockService.On("Delete", int64(2)).Return(err)

		// Create a test context
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest("DELETE", "/api/v1/roles/:id", nil)
		c.Params = gin.Params{{Key: "id", Value: "2"}}

		// Call the handler function
		handler.DeleteRole(c)

		// Assert the response
		var response map[string]any
		json.Unmarshal(w.Body.Bytes(), &response)
		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.Equal(t, float64(apperror.ErrDBDelete), response["code"])
		assert.Equal(t, "Database error", response["message"])

		// Assert mocks
		mockService.AssertExpectations(t)
	})
}
