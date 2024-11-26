package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/vfa-khuongdv/golang-cms/internal/handlers"
	"github.com/vfa-khuongdv/golang-cms/internal/repositories"
	"github.com/vfa-khuongdv/golang-cms/internal/services"
	"github.com/vfa-khuongdv/golang-cms/internal/utils"
	"gorm.io/gorm"
)

func SetupRouter(db *gorm.DB) *gin.Engine {
	router := gin.Default()
	userRepo := repositories.NewUserRepsitory(db)
	refreshRepo := repositories.NewRefreshTokenRepository(db)

	tokenService := services.NewRefreshTokenService(refreshRepo)
	userService := services.NewUserService(userRepo)
	authService := services.NewAuthService(userRepo, tokenService)

	authHandler := handlers.NewAuthHandler(authService)
	userHandler := handlers.NewUserHandler(userService)

	// Mode
	ginMode := utils.GetEnv("GIN_MODE", "debug")
	gin.SetMode(ginMode)

	// Middleware
	router.Use(gin.Recovery())

	// Router
	api := router.Group("/api/v1")
	{
		api.GET("/healthz", handlers.HealthCheck)
		api.POST("/auth/login", authHandler.LoginHandler)
		api.POST("/users", userHandler.RegisterHandler)
	}

	return router
}
