package configs

import (
	"fmt"

	log "github.com/sirupsen/logrus"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type DatabaseConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	DBName   string
	Charset  string
}

var DB *gorm.DB

// InitDB initializes a MySQL database connection using GORM
// Parameters:
//   - config: DatabaseConfig struct containing database connection parameters
//
// Returns:
//   - *gorm.DB: Database connection instance
//
// Note: Also sets the global DB variable with the connection instance
func InitDB(config DatabaseConfig) *gorm.DB {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=True&loc=Local",
		config.User,
		config.Password,
		config.Host,
		config.Port,
		config.DBName,
		config.Charset,
	)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to the MySQL database: %v", err)
	}

	log.Println("MySQL database connection established successfully")
	DB = db
	return db
}
