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
