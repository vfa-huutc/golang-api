package repository

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type DatabaseConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	DBName   string
	SSLMode  string
}

// InitDB initializes the database connection using GORM
var DB *gorm.DB

func ConnectionDB() {
	database, err := gorm.Open(sqlite.Open("test.db"))
	if err != nil {
		panic("Failed to connect to database!")
	}
	DB = database
}
