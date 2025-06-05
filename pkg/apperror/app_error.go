package apperror

import "fmt"

// AppError represents a custom error with a code and message.
type AppError struct {
	HttpStatusCode int    `json:"-"`       // HTTP status code (optional)
	Code           int    `json:"code"`    // Error code
	Message        string `json:"message"` // Error message
}

// Error implements the error interface.
func (e *AppError) Error() string {
	return fmt.Sprintf("code: %d, message: %s", e.Code, e.Message)
}

// Wrap creates a new AppError with an underlying error.
func Wrap(httpStatusCode, code int, message string, err error) *AppError {
	return &AppError{
		HttpStatusCode: httpStatusCode,
		Code:           code,
		Message:        message,
	}
}

// New creates a new AppError without an underlying error.
func New(httpStatusCode, code int, message string) *AppError {
	return &AppError{
		HttpStatusCode: httpStatusCode,
		Code:           code,
		Message:        message,
	}
}
