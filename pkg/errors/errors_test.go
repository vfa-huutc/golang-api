package errors_test

import (
	originError "errors"
	"testing"

	"github.com/vfa-khuongdv/golang-cms/pkg/errors"

	"github.com/stretchr/testify/assert"
)

func TestAppError_Error(t *testing.T) {
	t.Run("without underlying error", func(t *testing.T) {
		appErr := errors.New(errors.ErrServerInternal, "internal error")
		expected := "code: 1000, message: internal error"
		assert.Equal(t, expected, appErr.Error())
	})

	t.Run("with underlying error", func(t *testing.T) {
		underlying := originError.New("database down")
		appErr := errors.Wrap(errors.ErrDatabaseConnection, "db connection failed", underlying)
		expected := "code: 2000, message: db connection failed, error: database down"
		assert.Equal(t, expected, appErr.Error())
	})
}

func TestWrap(t *testing.T) {
	underlying := originError.New("underlying error")
	appErr := errors.Wrap(errors.ErrInvalidRequest, "invalid request", underlying)

	assert.NotNil(t, appErr)
	assert.Equal(t, errors.ErrInvalidRequest, appErr.Code)
	assert.Equal(t, "invalid request", appErr.Message)
	assert.Equal(t, underlying, appErr.Err)
}

func TestNew(t *testing.T) {
	appErr := errors.New(errors.ErrAuthUnauthorized, "unauthorized")

	assert.NotNil(t, appErr)
	assert.Equal(t, errors.ErrAuthUnauthorized, appErr.Code)
	assert.Equal(t, "unauthorized", appErr.Message)
	assert.Nil(t, appErr.Err)
}
