package handlers_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/vfa-khuongdv/golang-cms/internal/handlers"
	"github.com/vfa-khuongdv/golang-cms/internal/models"
	"github.com/vfa-khuongdv/golang-cms/internal/utils"
	"github.com/vfa-khuongdv/golang-cms/pkg/apperror"
	"github.com/vfa-khuongdv/golang-cms/tests/mocks"
)

func TestGetSettings(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("GetSetting - Success", func(t *testing.T) {
		mockService := new(mocks.MockSettingService)
		handler := handlers.NewSettingHandler(mockService)

		expected := []models.Setting{
			{
				ID:         1,
				SettingKey: "site_name",
				Value:      "My Site",
				CreatedAt:  time.Date(2023, 10, 1, 0, 0, 0, 0, time.UTC),
				UpdatedAt:  time.Date(2023, 10, 1, 0, 0, 0, 0, time.UTC),
			},
		}
		// Mock the service method
		mockService.On("GetSetting").Return(expected, nil).Once()

		// Create a test context
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest("GET", "/api/v1/settings", nil)

		// Call the handler
		handler.GetSettings(c)

		// Assert the response
		var actual []models.Setting
		json.Unmarshal(w.Body.Bytes(), &actual)
		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, expected, actual)

		// Assert mocks
		mockService.AssertExpectations(t)
	})

	t.Run("GetSetting - Error", func(t *testing.T) {
		mockService := new(mocks.MockSettingService)
		handler := handlers.NewSettingHandler(mockService)

		// Mock the service methods
		mockService.On("GetSetting").Return(([]models.Setting)(nil), apperror.NewNotFoundError("Not found record")).Once()

		// Create a test context
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest("GET", "/api/v1/settings", nil)

		// Call the handler
		handler.GetSettings(c)

		// Assert the response

		var expected = map[string]any{
			"code":    float64(apperror.ErrNotFound),
			"message": "Not found record",
		}
		var actual map[string]any
		json.Unmarshal(w.Body.Bytes(), &actual)

		assert.Equal(t, expected["code"], actual["code"])
		assert.Equal(t, expected["message"], actual["message"])
		assert.Equal(t, http.StatusNotFound, w.Code)

		// Assert mocks
		mockService.AssertExpectations(t)
	})

}

func TestUpdateSettings(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("UpdateSetting - Success", func(t *testing.T) {
		mockService := new(mocks.MockSettingService)
		handler := handlers.NewSettingHandler(mockService)

		// Mock the service methods
		requestBody := map[string]interface{}{
			"settings": []map[string]string{
				{"key": "site_name", "value": "New Site"},
				{"key": "site_url", "value": "https://example.com"},
			},
		}
		body, _ := json.Marshal(requestBody)
		mockService.On("GetSettingByKey", "site_name").Return(&models.Setting{SettingKey: "site_name", Value: "Old"}, nil).Once()
		mockService.On("Update", mock.MatchedBy(func(s *models.Setting) bool {
			return s.SettingKey == "site_name" && s.Value == "New Site"
		})).Return(nil).Once()
		mockService.On("GetSettingByKey", "site_url").Return((&models.Setting{}), apperror.NewNotFoundError("not found")).Once()
		mockService.On("Create", mock.MatchedBy(func(s *models.Setting) bool {
			return s.SettingKey == "site_url" && s.Value == "https://example.com"
		})).Return(nil).Once()

		// Create a test context and request
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest("PUT", "/api/v1/settings", bytes.NewBuffer(body))
		c.Request.Header.Set("Content-Type", "application/json")

		// Call the handler method
		handler.UpdateSettings(c)

		// Assert the response
		assert.Equal(t, http.StatusOK, w.Code)
		assert.JSONEq(t, `{"message":"Update setting successfully"}`, w.Body.String())

		// Assert mocks
		mockService.AssertExpectations(t)
	})

	t.Run("UpdateSetting - Validation error", func(t *testing.T) {
		mockService := new(mocks.MockSettingService)
		handler := handlers.NewSettingHandler(mockService)

		// Missing 'value'
		requestBody := map[string]interface{}{
			"settings": []map[string]string{
				{"key": "site_name"},
			},
		}
		body, _ := json.Marshal(requestBody)

		// Create a test context
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest("PUT", "/api/v1/settings", bytes.NewBuffer(body))

		// Call the handler
		handler.UpdateSettings(c)

		// Assert the response
		expectedBody := map[string]any{
			"code":    float64(apperror.ErrValidationFailed),
			"message": "Validation failed",
			"fields": []apperror.FieldError{
				{
					Field:   "settings.Settings[0].Value",
					Message: "settings.Settings[0].Value is required",
				},
			},
		}
		var actualBody map[string]any
		json.Unmarshal(w.Body.Bytes(), &actualBody)
		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Equal(t, expectedBody["code"], actualBody["code"])
		assert.Equal(t, expectedBody["message"], actualBody["message"])
		assert.Equal(t, expectedBody["fields"], utils.ToFieldErrors(actualBody["fields"]))

		// Assert mocks
		mockService.AssertExpectations(t)
	})

	t.Run("UpdateSetting - Failed to update setting", func(t *testing.T) {
		mockService := new(mocks.MockSettingService)
		handler := handlers.NewSettingHandler(mockService)

		// Mock the service methods
		requestBody := map[string]interface{}{
			"settings": []map[string]string{
				{"key": "site_name", "value": "New Site"},
			},
		}
		body, _ := json.Marshal(requestBody)
		mockService.On("GetSettingByKey", "site_name").Return(&models.Setting{SettingKey: "site_name", Value: "Old"}, nil).Once()
		mockService.On("Update", mock.Anything).Return(apperror.NewDBUpdateError("update failed")).Once()

		// Create a test context and request
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest("PUT", "/api/v1/settings", bytes.NewBuffer(body))
		c.Request.Header.Set("Content-Type", "application/json")

		// Call the handler method
		handler.UpdateSettings(c)

		// Assert the response
		assert.Equal(t, http.StatusOK, w.Code)
		assert.JSONEq(t, `{"message":"Update setting successfully"}`, w.Body.String()) // Still returns success even if one fails internally

		// Assert mocks
		mockService.AssertExpectations(t)
	})

	t.Run("UpdateSetting - Fails but continues", func(t *testing.T) {
		// Mock the service and handler
		mockService := new(mocks.MockSettingService)
		handler := handlers.NewSettingHandler(mockService)

		// Mock the service methods
		requestBody := map[string]interface{}{
			"settings": []map[string]string{
				{"key": "existing_setting", "value": "new_value"},
			},
		}
		body, _ := json.Marshal(requestBody)
		mockService.On("GetSettingByKey", "existing_setting").
			Return(&models.Setting{SettingKey: "existing_setting", Value: "old_value"}, nil).Once()
		mockService.On("Update", mock.MatchedBy(func(s *models.Setting) bool {
			return s.SettingKey == "existing_setting" && s.Value == "new_value"
		})).Return(apperror.NewDBUpdateError("update failed")).Once()

		// Create a test context
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest("PUT", "/api/v1/settings", bytes.NewBuffer(body))
		c.Request.Header.Set("Content-Type", "application/json")

		// Call the handler method
		handler.UpdateSettings(c)

		// Assert the response
		assert.Equal(t, http.StatusOK, w.Code)
		assert.JSONEq(t, `{"message":"Update setting successfully"}`, w.Body.String())

		// Assert mocks
		mockService.AssertExpectations(t)
	})
}
