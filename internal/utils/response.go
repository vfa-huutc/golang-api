package utils

import (
	"github.com/gin-gonic/gin"
	"github.com/vfa-khuongdv/golang-cms/constants"
)

// RespondWithError sends a standardized error response with a default status code
func RespondWithError(ctx *gin.Context, err constants.ErrorResponse) {
	ctx.JSON(err.StatusCode, gin.H{
		"code":    err.Code,
		"message": err.Message,
	})
}
