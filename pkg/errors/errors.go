package errors

import "fmt"

const (
	// General errors
	ErrInternal   = 1000 // Internal server error
	ErrNotFound   = 1001 // Resource not found
	ErrBadRequest = 1002 // Invalid or bad request

	// Database errors
	ErrDBConnection = 2000 // Failed to connect to DB
	ErrDBQuery      = 2001 // DB query error
	ErrDBInsert     = 2002 // DB insert error
	ErrDBUpdate     = 2003 // DB update error
	ErrDBDelete     = 2004 // DB delete error

	// Authentication errors
	ErrUnauthorized       = 3000 // Unauthorized access
	ErrForbidden          = 3001 // Forbidden access
	ErrTokenExpired       = 3002 // Token has expired
	ErrInvalidPassword    = 3003 // Invalid password
	ErrPasswordHashFailed = 3004 // Failed to hash password
	ErrPasswordMismatch   = 3005 // Password mismatch
	ErrPasswordUnchanged  = 3006 // Old and new password are the same

	// Common / misc errors
	ErrParseError  = 4000 // Parsing or field error
	ErrInvalidData = 4001 // Validation failed
	ErrCacheSet    = 4002 // Set cache error
	ErrCacheGet    = 4003 // Get cache error
	ErrCacheDelete = 4004 // Delete cache error
	ErrCacheList   = 4005 // List cache error
	ErrCacheExists = 4006 // Cache key exists check error
)

// AppError represents a custom error with a code and message.
type AppError struct {
	Code    int    `json:"code"`    // Error code
	Message string `json:"message"` // Error message
	Err     error  `json:"-"`       // Underlying error (optional)
}

// Error implements the error interface.
func (e *AppError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("code: %d, message: %s, error: %v", e.Code, e.Message, e.Err)
	}
	return fmt.Sprintf("code: %d, message: %s", e.Code, e.Message)
}

// Wrap creates a new AppError with an underlying error.
func Wrap(code int, message string, err error) *AppError {
	return &AppError{
		Code:    code,
		Message: message,
		Err:     err,
	}
}

// New creates a new AppError without an underlying error.
func New(code int, message string) *AppError {
	return &AppError{
		Code:    code,
		Message: message,
	}
}
