package utils

import (
	"github.com/gin-gonic/gin"
	"github.com/vfa-khuongdv/golang-cms/constants"
)

// RespondWithError sends a standardized error response to the client
// Parameters:
//   - ctx: Gin context for handling the HTTP response
//   - err: ErrorResponse struct containing status code, error code and message
//
// Returns: None
func RespondWithError(ctx *gin.Context, err constants.ErrorResponse) {
	ctx.JSON(err.StatusCode, gin.H{
		"code":    err.Code,
		"message": err.Message,
	})
}
