package handlers_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/vfa-khuongdv/golang-cms/internal/handlers"
	"github.com/vfa-khuongdv/golang-cms/internal/models"
	"github.com/vfa-khuongdv/golang-cms/pkg/apperror"
	"github.com/vfa-khuongdv/golang-cms/tests/mocks"
)

func TestGetAllPermissions(t *testing.T) {
	gin.SetMode(gin.TestMode)
	t.Run("GetAll - Success", func(t *testing.T) {
		mockService := new(mocks.MockPermissionService)
		handler := handlers.NewPermissionHandler(mockService)

		expectedPermissions := []models.Permission{
			{
				ID:        1,
				Resource:  "read",
				Action:    "Read Permission",
				CreatedAt: time.Date(2023, 10, 1, 0, 0, 0, 0, time.UTC),
				UpdatedAt: time.Date(2023, 10, 1, 0, 0, 0, 0, time.UTC),
			},
			{
				ID:        2,
				Resource:  "write",
				Action:    "Write Permission",
				CreatedAt: time.Date(2023, 10, 1, 0, 0, 0, 0, time.UTC),
				UpdatedAt: time.Date(2023, 10, 1, 0, 0, 0, 0, time.UTC),
			},
		}

		// Mock the service method
		mockService.On("GetAll").Return(expectedPermissions, nil).Once()

		// Create a test context
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest("GET", "/api/v1/permissions", nil)

		// Call the handler
		handler.GetAll(c)

		// Assert the response
		assert.Equal(t, http.StatusOK, w.Code)

		assert.JSONEq(t, `[
		{
			"id": 1,
			"resource": "read",
			"action": "Read Permission",
			"created_at": "2023-10-01T00:00:00Z",
			"updated_at": "2023-10-01T00:00:00Z",
			"deleted_at": null
		},
		{
			"id": 2,
			"resource": "write",
			"action": "Write Permission",
			"created_at": "2023-10-01T00:00:00Z",
			"updated_at": "2023-10-01T00:00:00Z",
			"deleted_at": null
		}
	]`, w.Body.String())

		// Assert mocks
		mockService.AssertExpectations(t)
	})

	t.Run("GetAll - Error", func(t *testing.T) {
		mockService := new(mocks.MockPermissionService)
		handler := handlers.NewPermissionHandler(mockService)

		// Mock the service method
		mockService.On("GetAll").Return(nil, apperror.NewDBQueryError("Query error")).Once()

		// Create a test context
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest("GET", "/api/v1/permissions", nil)

		// Call the handler
		handler.GetAll(c)

		// Assert the response
		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.JSONEq(t, fmt.Sprintf(`{"code":%d,"message":"%s"}`, apperror.ErrDBQuery, "Query error"), w.Body.String())

		// Assert mocks
		mockService.AssertExpectations(t)
	})
}
