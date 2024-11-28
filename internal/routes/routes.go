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
	roleRepo := repositories.NewRoleRepository(db)
	settingRepo := repositories.NewSettingRepository(db)

	// Initialize services
	tokenService := services.NewRefreshTokenService(refreshRepo)
	userService := services.NewUserService(userRepo)
	authService := services.NewAuthService(userRepo, tokenService)
	roleService := services.NewRoleService(roleRepo)
	settingService := services.NewSettingService(settingRepo)

	// Initialize handlers
	authHandler := handlers.NewAuthHandler(authService)
	userHandler := handlers.NewUserHandler(userService)
	roleHandler := handlers.NewRoleHandler(roleService)
	settingHandler := handlers.NewSettingHandler(settingService)

	// Set Gin mode from environment variable
	ginMode := utils.GetEnv("GIN_MODE", "debug")
	gin.SetMode(ginMode)

	// Add middleware
	router.Use(gin.Recovery())

	// Setup API routes
	api := router.Group("/api/v1")
	{
		api.GET("/healthz", handlers.HealthCheck)
		api.POST("/auth/login", authHandler.Login)
		api.POST("/auth/refresh-token", authHandler.RefreshToken)

		api.POST("/users", userHandler.CreateUser)

		api.POST("/roles", roleHandler.CreateRole)
		api.GET("/roles/:id", roleHandler.GetRole)
		api.PATCH("/roles/:id", roleHandler.UpdateRole)
		api.DELETE("/roles/:id", roleHandler.DeleteRole)

		api.GET("/settings", settingHandler.GetSettings)
	}

	return router
}
