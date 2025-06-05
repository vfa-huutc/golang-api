package errors_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	app_errors "github.com/vfa-khuongdv/golang-cms/pkg/errors"
)

func TestAppError_Error(t *testing.T) {
	t.Run("without underlying error", func(t *testing.T) {
		appErr := app_errors.New(app_errors.ErrInternal, "internal error")
		expected := "code: 1000, message: internal error"
		assert.Equal(t, expected, appErr.Error())
	})

	t.Run("with underlying error", func(t *testing.T) {
		underlying := errors.New("database down")
		appErr := app_errors.Wrap(app_errors.ErrDBConnection, "db connection failed", underlying)
		expected := "code: 2000, message: db connection failed, error: database down"
		assert.Equal(t, expected, appErr.Error())
	})
}

func TestWrap(t *testing.T) {
	underlying := errors.New("underlying error")
	appErr := app_errors.Wrap(app_errors.ErrInvalidData, "invalid request", underlying)

	assert.NotNil(t, appErr)
	assert.Equal(t, app_errors.ErrInvalidData, appErr.Code)
	assert.Equal(t, "invalid request2", appErr.Message)
	assert.Equal(t, underlying, appErr.Err)
}

func TestNew(t *testing.T) {
	appErr := app_errors.New(app_errors.ErrUnauthorized, "unauthorized")

	assert.NotNil(t, appErr)
	assert.Equal(t, app_errors.ErrUnauthorized, appErr.Code)
	assert.Equal(t, "unauthorized", appErr.Message)
	assert.Nil(t, appErr.Err)
}
