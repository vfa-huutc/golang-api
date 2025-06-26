package apperror_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vfa-khuongdv/golang-cms/pkg/apperror"
)

func TestValidateError(t *testing.T) {
	// Example test case for the Error method
	err := apperror.NewValidationError("Validation failed", []apperror.FieldError{
		{Field: "username", Message: "Username is required"},
		{Field: "email", Message: "Email is invalid"},
	})

	expected := "code: 4001, message: Validation failed, fields: [{username Username is required} {email Email is invalid}]"

	assert.Equal(t, expected, err.Error())
}

func TestValidateWrap(t *testing.T) {
	// Example test case for the Wrap method
	err := apperror.NewValidationError("Validation failed", []apperror.FieldError{
		{Field: "password", Message: "Password is too short"},
	})

	wrappedErr := err.Wrap(400, 4001, "Wrapped validation error")
	expected := "code: 4001, message: Wrapped validation error, fields: [{password Password is too short}]"

	assert.NotNil(t, wrappedErr)
	assert.Equal(t, 4001, wrappedErr.Code)
	assert.Equal(t, "Wrapped validation error", wrappedErr.Message)
	assert.Equal(t, 1, len(wrappedErr.Fields))
	assert.Equal(t, expected, wrappedErr.Error())
}

func TestValidateNewValidationError(t *testing.T) {
	// Example test case for the NewValidationError function
	fieldErrors := []apperror.FieldError{
		{Field: "email", Message: "Email is required"},
		{Field: "password", Message: "Password must be at least 6 characters"},
	}

	err := apperror.NewValidationError("Validation failed", fieldErrors)

	assert.Equal(t, 4001, err.Code)
	assert.Equal(t, "Validation failed", err.Message)
	assert.Equal(t, 2, len(err.Fields))

	for _, fieldError := range err.Fields {
		switch fieldError.Field {
		case "email":
			assert.Equal(t, "Email is required", fieldError.Message)
		case "password":
			assert.Equal(t, "Password must be at least 6 characters", fieldError.Message)
		}
	}

}
