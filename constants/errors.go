package constants

import "net/http"

// ErrorResponse defines the structure for API error responses
type ErrorResponse struct {
	Code       string `json:"code"`
	Message    string `json:"message"`
	StatusCode int    `json:"status_code"`
}

var ErrorMessages = map[string]ErrorResponse{
	"InvalidCredentials": {
		Code:       "B0001",
		Message:    "Username and password are required",
		StatusCode: http.StatusBadRequest,
	},
	"Unauthorized": {
		Code:       "C0001",
		Message:    "Unauthorized",
		StatusCode: http.StatusUnauthorized,
	},
	"Forbidden": {
		Code:       "C0002",
		Message:    "Forbidden",
		StatusCode: http.StatusForbidden,
	},
	"NotFound": {
		Code:       "C0003",
		Message:    "Not Found",
		StatusCode: http.StatusNotFound,
	},
	"InvalidAuthorization": {
		Code:       "C0004",
		Message:    "Authorization header required",
		StatusCode: http.StatusUnauthorized,
	},
	"InvalidToken": {
		Code:       "C0005",
		Message:    "Authorization header required",
		StatusCode: http.StatusUnauthorized,
	},
	"Internal": {
		Code:       "A0001",
		Message:    "Internal Server Error",
		StatusCode: http.StatusInternalServerError,
	},
}
