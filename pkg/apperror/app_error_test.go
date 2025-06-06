package apperror_test

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vfa-khuongdv/golang-cms/pkg/apperror"
)

func TestAppError_Error(t *testing.T) {
	t.Run("without underlying error", func(t *testing.T) {
		appErr := apperror.New(
			http.StatusInternalServerError,
			apperror.ErrInternal,
			"internal error",
		)
		expected := "code: 1000, message: internal error"
		assert.Equal(t, expected, appErr.Error())
	})

}

func TestWrap(t *testing.T) {
	underlying := apperror.New(
		http.StatusBadRequest,
		apperror.ErrBadRequest,
		"invalid input",
	)
	appErr := apperror.Wrap(
		http.StatusBadRequest,
		apperror.ErrValidationFailed,
		"invalid request",
		underlying,
	)

	assert.NotNil(t, appErr)
	assert.Equal(t, apperror.ErrValidationFailed, appErr.Code)
	assert.Equal(t, "invalid request", appErr.Message)
}

func TestNew(t *testing.T) {
	appErr := apperror.New(
		http.StatusUnauthorized,
		apperror.ErrUnauthorized,
		"unauthorized",
	)

	assert.NotNil(t, appErr)
	assert.Equal(t, apperror.ErrUnauthorized, appErr.Code)
	assert.Equal(t, "unauthorized", appErr.Message)
}

func TestIsAppError(t *testing.T) {
	t.Run("is AppError", func(t *testing.T) {
		appErr := apperror.New(
			http.StatusForbidden,
			apperror.ErrForbidden,
			"forbidden",
		)
		assert.True(t, apperror.IsAppError(appErr))
	})

	t.Run("is not AppError", func(t *testing.T) {
		err := assert.AnError
		assert.False(t, apperror.IsAppError(err))
	})
}

func TestToAppError(t *testing.T) {
	t.Run("is AppError", func(t *testing.T) {
		appErr := apperror.New(
			http.StatusNotFound,
			apperror.ErrNotFound,
			"not found",
		)
		result, ok := apperror.ToAppError(appErr)
		assert.True(t, ok)
		assert.Equal(t, appErr, result)
	})

	t.Run("is not AppError", func(t *testing.T) {
		err := assert.AnError
		result, ok := apperror.ToAppError(err)
		assert.False(t, ok)
		assert.Nil(t, result)
	})
}

func TestAppErrorWithUnderlyingError(t *testing.T) {
	underlying := assert.AnError
	appErr := apperror.Wrap(
		http.StatusInternalServerError,
		apperror.ErrInternal,
		"internal error",
		underlying,
	)

	expected := "code: 1000, message: internal error, error: " + underlying.Error()
	assert.Equal(t, expected, appErr.Error())
	assert.Equal(t, http.StatusInternalServerError, appErr.HttpStatusCode)
}
