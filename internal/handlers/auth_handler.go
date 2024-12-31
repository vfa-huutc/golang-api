package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/vfa-khuongdv/golang-cms/internal/services"
	"github.com/vfa-khuongdv/golang-cms/internal/utils"
	"github.com/vfa-khuongdv/golang-cms/pkg/errors"
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

func (handler *AuthHandler) Login(ctx *gin.Context) {
	var credentials struct {
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required,min=6,max=255"`
	}

	if err := ctx.ShouldBindJSON(&credentials); err != nil {
		utils.RespondWithError(
			ctx,
			http.StatusBadRequest,
			errors.New(errors.ErrCodeValidation, err.Error()),
		)
		return
	}

	// login handler
	res, err := handler.authService.Login(credentials.Email, credentials.Password, ctx)
	if err != nil {
		utils.RespondWithError(ctx, http.StatusUnauthorized, err)
		return
	}

	utils.RespondWithOK(ctx, http.StatusOK, res)
}

func (handler *AuthHandler) RefreshToken(ctx *gin.Context) {
	var token struct {
		RefreshToken string `json:"refresh_token" binding:"required"`
	}

	// Bind JSON request body to token struct
	if err := ctx.ShouldBindJSON(&token); err != nil {
		utils.RespondWithError(
			ctx,
			http.StatusBadRequest,
			errors.New(errors.ErrCodeValidation, err.Error()),
		)
		return
	}

	// Call auth service to refresh the token
	res, err := handler.authService.RefreshToken(token.RefreshToken, ctx)
	if err != nil {
		utils.RespondWithError(ctx, http.StatusUnauthorized, err)
		return
	}

	utils.RespondWithOK(ctx, http.StatusOK, res)
}
