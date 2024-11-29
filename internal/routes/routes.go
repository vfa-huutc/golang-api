package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/vfa-khuongdv/golang-cms/internal/handlers"
	"github.com/vfa-khuongdv/golang-cms/internal/middlewares"
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
	permissionRepo := repositories.NewPermissionRepository(db)

	// Initialize services
	tokenService := services.NewRefreshTokenService(refreshRepo)
	userService := services.NewUserService(userRepo)
	authService := services.NewAuthService(userRepo, tokenService)
	roleService := services.NewRoleService(roleRepo)
	settingService := services.NewSettingService(settingRepo)
	permissionService := services.NewPermissionService(permissionRepo)

	// Initialize handlers
	authHandler := handlers.NewAuthHandler(authService)
	userHandler := handlers.NewUserHandler(userService)
	roleHandler := handlers.NewRoleHandler(roleService)
	settingHandler := handlers.NewSettingHandler(settingService)
	permissionHandler := handlers.NewPermissionHandler(permissionService)

	// Set Gin mode from environment variable
	ginMode := utils.GetEnv("GIN_MODE", "debug")
	gin.SetMode(ginMode)

	// Add middleware
	router.Use(gin.Recovery())

	router.GET("/healthz", handlers.HealthCheck)

	// Setup API routes
	api := router.Group("/api/v1")
	{
		// Public routes
		api.POST("/login", authHandler.Login)
		api.POST("/refresh-token", authHandler.RefreshToken)
		api.POST("/forgot-password", userHandler.ForgotPassword)
		api.POST("/reset-password", userHandler.ResetPassword)

		// Protected routes (require authentication)
		api.Use(middlewares.AuthMiddleware())
		{
			// Profile management
			api.POST("/change-password", userHandler.ChangePassword)
			api.GET("/profile", userHandler.GetProfile)
			api.PATCH("/profile", userHandler.UpdateProfile)

			// User management
			api.POST("/users", userHandler.CreateUser)
			api.PATCH("/users/:id", userHandler.UpdateUser)
			api.DELETE("/users/:id", userHandler.DeleteUser)

			// Role management
			api.POST("/roles", roleHandler.CreateRole)
			api.GET("/roles/:id", roleHandler.GetRole)
			api.PATCH("/roles/:id", roleHandler.UpdateRole)
			api.DELETE("/roles/:id", roleHandler.DeleteRole)

			// Settings
			api.GET("/settings", settingHandler.GetSettings)
			api.PUT("/settings", settingHandler.UpdateSettings)

			// Permissions
			api.GET("/permissions", permissionHandler.GetAll)
		}
	}

	return router
}
