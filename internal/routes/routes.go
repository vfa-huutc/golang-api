package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/vfa-khuongdv/golang-cms/internal/guards"
	"github.com/vfa-khuongdv/golang-cms/internal/handlers"
	"github.com/vfa-khuongdv/golang-cms/internal/middlewares"
	"github.com/vfa-khuongdv/golang-cms/internal/repositories"
	"github.com/vfa-khuongdv/golang-cms/internal/services"
	"github.com/vfa-khuongdv/golang-cms/internal/utils"
	"gorm.io/gorm"
)

func SetupRouter(db *gorm.DB) *gin.Engine {
	// Set Gin mode from environment variable
	ginMode := utils.GetEnv("GIN_MODE", "debug")
	gin.SetMode(ginMode)

	// Initialize the default Gin router
	router := gin.Default()

	stage := utils.GetEnv("STAGE", "dev")

	// Set up Swagger documentation only in non-production environments
	if stage != "prod" {
		router.StaticFile("/docs/swagger.json", "./docs/swagger.json")
		router.StaticFile("/swagger", "./docs/swagger.html")
		router.StaticFile("/api-docs", "./docs/swagger.html")
	}

	// Initialize repositories
	userRepo := repositories.NewUserRepsitory(db)
	refreshRepo := repositories.NewRefreshTokenRepository(db)
	roleRepo := repositories.NewRoleRepository(db)
	settingRepo := repositories.NewSettingRepository(db)
	permissionRepo := repositories.NewPermissionRepository(db)

	// Initialize services
	REDIS_HOST := utils.GetEnv("REDIS_HOST", "localhost:6379")
	REDIS_PASS := utils.GetEnv("REDIS_PASS", "")
	REDIS_DB := utils.GetEnvAsInt("REDIS_DB", 0)

	redisService := services.NewRedisService(REDIS_HOST, REDIS_PASS, REDIS_DB)
	tokenService := services.NewRefreshTokenService(refreshRepo)
	userService := services.NewUserService(userRepo)
	authService := services.NewAuthService(userRepo, tokenService)
	roleService := services.NewRoleService(roleRepo)
	settingService := services.NewSettingService(settingRepo)
	permissionService := services.NewPermissionService(permissionRepo)

	// Initialize role guard for permission checks
	roleGuard := guards.NewRoleGuard(db)

	// Initialize handlers
	authHandler := handlers.NewAuthHandler(authService)
	userHandler := handlers.NewUserHandler(userService, redisService)
	roleHandler := handlers.NewRoleHandler(roleService)
	settingHandler := handlers.NewSettingHandler(settingService)
	permissionHandler := handlers.NewPermissionHandler(permissionService)

	// Add middleware for CORS and logging
	router.Use(
		middlewares.CORSMiddleware(),
		middlewares.LogMiddleware(),
		gin.Recovery(),
		middlewares.EmptyBodyMiddleware(),
	)

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
		authenticated := api.Group("/")
		authenticated.Use(middlewares.AuthMiddleware())
		{
			// Profile management (available to all authenticated users)
			authenticated.POST("/change-password", userHandler.ChangePassword)
			authenticated.GET("/profile", userHandler.GetProfile)
			authenticated.PATCH("/profile", userHandler.UpdateProfile)

			// User management routes with permission checks
			authenticated.POST("/users", guards.RequirePermissions(roleGuard, "users:create"), userHandler.CreateUser)
			authenticated.GET("/users/:id", guards.RequirePermissions(roleGuard, "users:view"), userHandler.GetUser)
			authenticated.PATCH("/users/:id", guards.RequirePermissions(roleGuard, "users:update"), userHandler.UpdateUser)
			authenticated.DELETE("/users/:id", guards.RequirePermissions(roleGuard, "users:delete"), userHandler.DeleteUser)

			// Role management routes with permission checks
			authenticated.POST("/roles", guards.RequirePermissions(roleGuard, "roles:create"), roleHandler.CreateRole)
			authenticated.GET("/roles/:id", guards.RequirePermissions(roleGuard, "roles:view"), roleHandler.GetRole)
			authenticated.PATCH("/roles/:id", guards.RequirePermissions(roleGuard, "roles:update"), roleHandler.UpdateRole)
			authenticated.DELETE("/roles/:id", guards.RequirePermissions(roleGuard, "roles:delete"), roleHandler.DeleteRole)

			// Settings with permission checks
			authenticated.GET("/settings", guards.RequirePermissions(roleGuard, "settings:view"), settingHandler.GetSettings)
			authenticated.PATCH("/settings", guards.RequirePermissions(roleGuard, "settings:update"), settingHandler.UpdateSettings)

			// Permissions with permission checks
			authenticated.GET("/permissions", guards.RequirePermissions(roleGuard, "roles:view"), permissionHandler.GetAll)
		}
	}

	return router
}
