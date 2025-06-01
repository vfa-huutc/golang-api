package handlers_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/vfa-khuongdv/golang-cms/internal/handlers"
	"github.com/vfa-khuongdv/golang-cms/internal/models"
	appError "github.com/vfa-khuongdv/golang-cms/pkg/errors"
	"github.com/vfa-khuongdv/golang-cms/tests/mocks"
)

func TestGetSettings(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("Success", func(t *testing.T) {
		mockService := new(mocks.MockSettingService)
		expected := []models.Setting{
			{
				ID:         1,
				SettingKey: "site_name",
				Value:      "My Site",
				CreatedAt:  time.Date(2023, 10, 1, 0, 0, 0, 0, time.UTC),
				UpdatedAt:  time.Date(2023, 10, 1, 0, 0, 0, 0, time.UTC),
			},
		}
		mockService.On("GetSetting").Return(expected, nil).Once()
		handler := handlers.NewSettingHandler(mockService)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest("GET", "/api/v1/settings", nil)
		handler.GetSettings(c)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.JSONEq(t, `[
			{ "id": 1, "settingKey": "site_name", "value": "My Site", "createdAt": "2023-10-01T00:00:00Z","updatedAt": "2023-10-01T00:00:00Z","deletedAt": null }
		]`, w.Body.String())
		mockService.AssertExpectations(t)
	})

	t.Run("Error", func(t *testing.T) {
		mockService := new(mocks.MockSettingService)
		mockService.On("GetSetting").Return(([]models.Setting)(nil), appError.New(appError.ErrDBQuery, "Query error")).Once()
		handler := handlers.NewSettingHandler(mockService)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest("GET", "/api/v1/settings", nil)
		handler.GetSettings(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.JSONEq(t, fmt.Sprintf(`{"code":%d,"message":"%s"}`, appError.ErrDBQuery, "Query error"), w.Body.String())
		mockService.AssertExpectations(t)
	})

	t.Run("Create setting fails but continues", func(t *testing.T) {
		mockService := new(mocks.MockSettingService)
		handler := handlers.NewSettingHandler(mockService)

		requestBody := map[string]interface{}{
			"settings": []map[string]string{
				{"key": "new_setting", "value": "value1"},
			},
		}
		body, _ := json.Marshal(requestBody)

		// Simulate "not found" on GetSettingByKey to trigger Create
		mockService.On("GetSettingByKey", "new_setting").
			Return(&models.Setting{}, appError.New(appError.ErrDBQuery, "not found")).Once()

		// Simulate Create failure
		mockService.On("Create", mock.MatchedBy(func(s *models.Setting) bool {
			return s.SettingKey == "new_setting" && s.Value == "value1"
		})).Return(appError.New(appError.ErrDBInsert, "insert failed")).Once()

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest("PUT", "/api/v1/settings", bytes.NewBuffer(body))
		c.Request.Header.Set("Content-Type", "application/json")

		handler.UpdateSettings(c)

		// Despite the create failure, response is still 200 OK with success message
		assert.Equal(t, http.StatusOK, w.Code)
		assert.JSONEq(t, `{"message":"Update setting successfully"}`, w.Body.String())

		mockService.AssertExpectations(t)
	})

}

func TestUpdateSettings(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("Success - update and create settings", func(t *testing.T) {
		mockService := new(mocks.MockSettingService)
		handler := handlers.NewSettingHandler(mockService)

		requestBody := map[string]interface{}{
			"settings": []map[string]string{
				{"key": "site_name", "value": "New Site"},
				{"key": "site_url", "value": "https://example.com"},
			},
		}
		body, _ := json.Marshal(requestBody)

		// Mock: site_name exists → Update
		mockService.On("GetSettingByKey", "site_name").Return(&models.Setting{SettingKey: "site_name", Value: "Old"}, nil).Once()
		mockService.On("Update", mock.MatchedBy(func(s *models.Setting) bool {
			return s.SettingKey == "site_name" && s.Value == "New Site"
		})).Return(nil).Once()

		// Mock: site_url not exist → Create
		mockService.On("GetSettingByKey", "site_url").Return((&models.Setting{}), appError.New(appError.ErrDBQuery, "not found")).Once()
		mockService.On("Create", mock.MatchedBy(func(s *models.Setting) bool {
			return s.SettingKey == "site_url" && s.Value == "https://example.com"
		})).Return(nil).Once()

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest("PUT", "/api/v1/settings", bytes.NewBuffer(body))
		c.Request.Header.Set("Content-Type", "application/json")

		handler.UpdateSettings(c)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.JSONEq(t, `{"message":"Update setting successfully"}`, w.Body.String())

		mockService.AssertExpectations(t)
	})

	t.Run("Fail - validation error", func(t *testing.T) {
		mockService := new(mocks.MockSettingService)
		handler := handlers.NewSettingHandler(mockService)

		// Missing 'value'
		body := `{"settings":[{"key":"site_name"}]}`

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest("PUT", "/api/v1/settings", bytes.NewBufferString(body))
		c.Request.Header.Set("Content-Type", "application/json")

		handler.UpdateSettings(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		var res map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &res)
		assert.Equal(t, float64(appError.ErrInvalidData), res["code"])
	})

	t.Run("Partial Failure - one update fails", func(t *testing.T) {
		mockService := new(mocks.MockSettingService)
		handler := handlers.NewSettingHandler(mockService)

		requestBody := map[string]interface{}{
			"settings": []map[string]string{
				{"key": "site_name", "value": "New Site"},
			},
		}
		body, _ := json.Marshal(requestBody)

		mockService.On("GetSettingByKey", "site_name").Return(&models.Setting{SettingKey: "site_name", Value: "Old"}, nil).Once()
		mockService.On("Update", mock.Anything).Return(appError.New(appError.ErrDBUpdate, "update failed")).Once()

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest("PUT", "/api/v1/settings", bytes.NewBuffer(body))
		c.Request.Header.Set("Content-Type", "application/json")

		handler.UpdateSettings(c)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.JSONEq(t, `{"message":"Update setting successfully"}`, w.Body.String()) // Still returns success even if one fails internally

		mockService.AssertExpectations(t)
	})

	t.Run("Update setting fails but continues", func(t *testing.T) {
		mockService := new(mocks.MockSettingService)
		handler := handlers.NewSettingHandler(mockService)

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
		})).Return(appError.New(appError.ErrDBUpdate, "update failed")).Once()

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest("PUT", "/api/v1/settings", bytes.NewBuffer(body))
		c.Request.Header.Set("Content-Type", "application/json")

		handler.UpdateSettings(c)

		// Response is still success despite update error
		assert.Equal(t, http.StatusOK, w.Code)
		assert.JSONEq(t, `{"message":"Update setting successfully"}`, w.Body.String())

		mockService.AssertExpectations(t)
	})

}
