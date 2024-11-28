package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/vfa-khuongdv/golang-cms/internal/services"
)

type IAuthHandler interface {
	Login(c *gin.Context)
	RefreshToken(c *gin.Context)
}

type AuthHandler struct {
	authService *services.AuthService
}

func NewAuthHandler(authService *services.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

func (handler *AuthHandler) Login(c *gin.Context) {
	var credentials struct {
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required,min=6,max=255"`
	}

	if err := c.ShouldBindJSON(&credentials); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// login handler
	res, err := handler.authService.Login(credentials.Email, credentials.Password, c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, res)

}

func (handler *AuthHandler) RefreshToken(c *gin.Context) {
	var token struct {
		RefreshToken string `json:"refresh_token"`
	}

	// Bind JSON request body to token struct
	if err := c.ShouldBindJSON(&token); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
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
