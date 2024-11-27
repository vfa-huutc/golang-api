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
	// Initialize the default Gin router
	router := gin.Default()

	// Initialize repositories
	userRepo := repositories.NewUserRepsitory(db)
	refreshRepo := repositories.NewRefreshTokenRepository(db)

	// Initialize services
	tokenService := services.NewRefreshTokenService(refreshRepo)
	userService := services.NewUserService(userRepo)
	authService := services.NewAuthService(userRepo, tokenService)

	// Initialize handlers
	authHandler := handlers.NewAuthHandler(authService)
	userHandler := handlers.NewUserHandler(userService)

	// Set Gin mode from environment variable
	ginMode := utils.GetEnv("GIN_MODE", "debug")
	gin.SetMode(ginMode)

	// Add middleware
	router.Use(gin.Recovery())

	// Setup API routes
	api := router.Group("/api/v1")
	{
		api.GET("/healthz", handlers.HealthCheck)
		api.POST("/auth/login", authHandler.LoginHandler)
		api.POST("/auth/refresh-token", authHandler.RefreshToken)
		api.POST("/users", userHandler.RegisterHandler)
	}

	return router
}
