package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/vfa-khuongdv/golang-cms/internal/handlers"
	"github.com/vfa-khuongdv/golang-cms/internal/repositories"
	"github.com/vfa-khuongdv/golang-cms/internal/services"
	"gorm.io/gorm"
)

func SetupRouter(db *gorm.DB) *gin.Engine {
	router := gin.Default()
	userRepo := repositories.NewUserRepsitory(db)
	authService := services.NewAuthService(userRepo)
	userService := services.NewUserService(userRepo)

	authHandler := handlers.NewAuthHandler(authService)
	userHandler := handlers.NewUserHandler(userService)

	// Middleware
	router.Use(gin.Recovery())

	// Routes
	api := router.Group("/api/v1")
	{
		api.GET("/healthz", handlers.HealthCheck)
		api.POST("/auth/login", authHandler.LoginHandler)
		api.POST("/users", userHandler.RegisterHandler)
	}

	return router
}
