package errors

import "fmt"

const (
	// General errors
	ErrCodeInternal   = 1000 // Internal server error
	ErrCodeNotFound   = 1001 // Not found
	ErrCodeBadRequest = 1002 // Bad request

	// Database errors
	ErrCodeDBConnection = 2000 //	Database connection error
	ErrCodeDBQuery      = 2001 // Database query error
	ErrCodeDBInsert     = 2002 //	Database insert error
	ErrCodeDBUpdate     = 2003 // Database update error
	ErrCodeDBDelete     = 2004 // Database delete error

	// Authentication errors
	ErrCodeUnauthorized        = 3000 // Unauthorized access
	ErrCodeForbidden           = 3001 // Forbidden access
	ErrCodeTokenExpired        = 3002 // Token has expired
	ErrCodeInvalidPassword     = 3003 // Invalid password
	ErrCodeFailedToHashed      = 3004 // Failed to hash password
	ErrCodeNotMatchedPassword  = 3005 // Password not matched
	ErrCodeOldAndNewShouldDiff = 3006 // Old and new password should be different

	// Common errors
	ErrCodeParseError = 4000 // Parse error, missing fields, etc.
	ErrCodeValidation = 4001 // Validation error
	ErrCodeSetCache   = 4002 // Set cache error
	ErrorGetCache     = 4003 // Get cache error
	ErrorDeleteCache  = 4004 // Delete cache error
	ErrorListCache    = 4005 // List cache error
	ErrorExistsCache  = 4006 // Exists cache error

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
