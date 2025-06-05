package apperror

import (
	"fmt"
)

type FieldError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

type ValidationError struct {
	Code    int          `json:"code"`    // Error code
	Message string       `json:"message"` // Error message
	Fields  []FieldError `json:"fields"`  // List of validation errors
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("code: %d, message: %s, fields: %v", e.Code, e.Message, e.Fields)
}

func (e *ValidationError) Wrap(httpStatusCode int, code int, message string) *ValidationError {
	return &ValidationError{
		Code:    code,
		Message: message,
		Fields:  e.Fields,
	}
}

func NewValidationError(message string, fieldErrors []FieldError) *ValidationError {
	return &ValidationError{
		Code:    ErrValidationFailed,
		Message: message,
		Fields:  fieldErrors,
	}
}
