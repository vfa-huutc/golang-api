package utils

import (
	"github.com/gin-gonic/gin"
	"github.com/vfa-khuongdv/golang-cms/pkg/errors"
)

func RespondWithError(ctx *gin.Context, statusCode int, err error) {
	if appErr, ok := err.(*errors.AppError); ok {
		ctx.JSON(
			statusCode,
			gin.H{"code": appErr.Code, "message": appErr.Message},
		)
		return
	} else {
		ctx.JSON(
			statusCode,
			gin.H{"code": errors.ErrCodeInternal, "message": err.Error()},
		)
		return
	}
}

func RespondWithOK(ctx *gin.Context, statusCode int, body interface{}) {
	ctx.JSON(statusCode, body)
}
