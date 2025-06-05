package middlewares_test

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/vfa-khuongdv/golang-cms/internal/middlewares"
	"github.com/vfa-khuongdv/golang-cms/pkg/apperror"
)

func TestEmptyBodyMiddleware_RejectsEmptyBody(t *testing.T) {
	router := gin.New()
	router.Use(middlewares.EmptyBodyMiddleware())
	router.POST("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "OK"})
	})

	req := httptest.NewRequest(http.MethodPost, "/test", nil)
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusBadRequest, resp.Code)

	expectedJSON := fmt.Sprintf(`{
		"code": %d,
		"message": "Request body cannot be empty"
	}`, apperror.ErrEmptyData)

	assert.JSONEq(t, expectedJSON, resp.Body.String())
}

func TestEmptyBodyMiddleware_AllowsNonEmptyBody(t *testing.T) {
	router := gin.New()
	router.Use(middlewares.EmptyBodyMiddleware())
	router.POST("/test", func(c *gin.Context) {
		body, _ := c.GetRawData()
		c.JSON(http.StatusOK, gin.H{"received": string(body)})
	})

	body := bytes.NewBufferString(`{"key":"value"}`)
	req := httptest.NewRequest(http.MethodPost, "/test", body)
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)
	assert.JSONEq(t, `{"received": "{\"key\":\"value\"}"}`, resp.Body.String())
}

func TestEmptyBodyMiddleware_IgnoreGET(t *testing.T) {
	router := gin.New()
	router.Use(middlewares.EmptyBodyMiddleware())
	router.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "OK"})
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)
	assert.JSONEq(t, `{"message": "OK"}`, resp.Body.String())
}
