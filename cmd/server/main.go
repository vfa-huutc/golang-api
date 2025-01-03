package main

import (
	"fmt"

	"github.com/vfa-khuongdv/golang-cms/internal/configs"
	"github.com/vfa-khuongdv/golang-cms/internal/models"
	"github.com/vfa-khuongdv/golang-cms/internal/routes"
	"github.com/vfa-khuongdv/golang-cms/internal/utils"
	"github.com/vfa-khuongdv/golang-cms/pkg/logger"
)

func main() {
	// Load env package
	configs.LoadEnv()

	// Init logger
	logger.Init()

	// MySQL database configuration
	config := configs.DatabaseConfig{
		Host:     utils.GetEnv("DB_HOST", "127.0.0.1"),
		Port:     utils.GetEnvAsInt("DB_PORT", 3306),
		User:     utils.GetEnv("DB_USERNAME", ""),
		Password: utils.GetEnv("DB_PASSWORD", ""),
		DBName:   utils.GetEnv("DB_DATABASE", ""),
		Charset:  "utf8mb4",
	}

	// Initialize database connection
	db := configs.InitDB(config)

	// Run migrations (for GORM)
	err := db.AutoMigrate(
		&models.User{},
		&models.RefreshToken{},
		&models.Role{},
		&models.Setting{},
		&models.Permission{},
	)
	if err != nil {
		logger.Fatal(fmt.Sprintf("Failed to run migrations: %+v", err))
	}
	logger.Info("Migrations completed")

	// Routes
	routes := routes.SetupRouter(db)
	// Port run server
	port := fmt.Sprintf(":%s", utils.GetEnv("PORT", "3000"))
	routes.Run(port)

}
