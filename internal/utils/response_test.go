package utils_test

import (
	originError "errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/vfa-khuongdv/golang-cms/internal/utils"
	"github.com/vfa-khuongdv/golang-cms/pkg/errors"
)

func TestRespondWithError_AppError(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)

	appErr := &errors.AppError{
		Code:    1001,
		Message: "App error occurred",
	}

	utils.RespondWithError(ctx, http.StatusBadRequest, appErr)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	expectedJSON := `{"code":1001,"message":"App error occurred"}`
	assert.JSONEq(t, expectedJSON, w.Body.String())
}

func TestRespondWithError_InternalServerError(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)

	internalErr := originError.New("Internal server error occurred")

	utils.RespondWithError(ctx, http.StatusInternalServerError, internalErr)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	expectedJSON := `{"code":1000,"message":"Internal server error occurred"}`
	assert.JSONEq(t, expectedJSON, w.Body.String())
}

func TestRespondWithError_GenericError(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)

	genericErr := errors.New(errors.ErrInternal, "generic error message")

	utils.RespondWithError(ctx, http.StatusInternalServerError, genericErr)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	expectedJSON := `{"code":1000,"message":"generic error message"}`
	assert.JSONEq(t, expectedJSON, w.Body.String())
}

func TestRespondWithOK(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)

	body := gin.H{"success": true, "data": "some data"}

	utils.RespondWithOK(ctx, http.StatusOK, body)

	assert.Equal(t, http.StatusOK, w.Code)
	expectedJSON := `{"success":true,"data":"some data"}`
	assert.JSONEq(t, expectedJSON, w.Body.String())
}
