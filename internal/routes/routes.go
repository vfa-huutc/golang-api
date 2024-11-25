package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/vfa-khuongdv/golang-cms/internal/handlers"
)

func SetupRouter() *gin.Engine {
	router := gin.Default()

	// Middleware
	router.Use(gin.Logger(), gin.Recovery())

	// Routes
	api := router.Group("/api")
	{
		api.GET("/healthz", handlers.HealthCheck)
	}

	return router
}
