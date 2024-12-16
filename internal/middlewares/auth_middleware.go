package middlewares

import (
	"net/http"
	"strings"
	"github.com/gin-gonic/gin"
	"github.com/vfa-khuongdv/golang-cms/configs"
)

// AuthMiddleware is a Gin middleware function that handles JWT authentication
// It validates the Authorization header and extracts the JWT token
// The middleware checks if:
// - Authorization header exists and has "Bearer " prefix
// - Token is valid and can be parsed
// If validation succeeds, it sets the user ID from token claims in context
// If validation fails, it returns 401 Unauthorized
func AuthMiddleware() gin.HandlerFunc {

	return func(ctx *gin.Context) {

		authHeader := ctx.GetHeader("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
			ctx.Abort()
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		claims, err := configs.ValidateToken(tokenString)
		if err != nil {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			ctx.Abort()
			return
		}

		ctx.Set("UserID", claims.ID)
		ctx.Next()
	}
}
