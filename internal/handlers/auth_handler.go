package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/vfa-khuongdv/golang-cms/internal/services"
)

type AuthHandler struct {
	authService *services.AuthService
}

func NewAuthHandler(authService *services.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

// LoginHandler process login request
func (handler *AuthHandler) LoginHandler(c *gin.Context) {
	var credentials struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if err := c.ShouldBindJSON(&credentials); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	// login handler
	res, err := handler.authService.Login(credentials.Username, credentials.Password, c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, res)

}

// RefreshToken handles the token refresh request
// It validates the refresh token and returns a new access token
// Parameters:
//   - c: Gin context containing the HTTP request/response
//
// Returns:
//   - 200 OK with new tokens on success
//   - 400 Bad Request if invalid JSON
//   - 401 Unauthorized if refresh token is invalid
func (handler *AuthHandler) RefreshToken(c *gin.Context) {
	var token struct {
		RefreshToken string `json:"refresh_token"`
	}

	// Bind JSON request body to token struct
	if err := c.ShouldBindJSON(&token); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	// Call auth service to refresh the token
	res, err := handler.authService.RefreshToken(token.RefreshToken, c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, res)

}
