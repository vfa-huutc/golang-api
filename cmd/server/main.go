package main

import (
	"log"

	"github.com/vfa-khuongdv/golang-cms/configs"
	"github.com/vfa-khuongdv/golang-cms/internal/models"
	"github.com/vfa-khuongdv/golang-cms/internal/routes"
)

func main() {
	// MySQL database configuration
	config := configs.DatabaseConfig{
		Host:     "127.0.0.1",
		Port:     3306,
		User:     "root",
		Password: "password",
		DBName:   "myapp",
		Charset:  "utf8mb4",
	}

	// Initialize database connection
	db := configs.ConnectDB(config)

	// Run migrations (for GORM)
	err := db.AutoMigrate(&models.User{})
	if err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	log.Println("Migrations completed")

	routes := routes.SetupRouter()
	routes.Run(":8080")

}
