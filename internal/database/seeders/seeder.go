package seeders

import (
	"github.com/vfa-khuongdv/golang-cms/pkg/logger"
	"gorm.io/gorm"
)

// Run executes all seed functions to populate the database with initial data
// It takes a GORM database connection as input and panics if any seeding operation fails
func Run(db *gorm.DB) {

	// Execute the SeedUsers function to populate initial user data
	// If an error occurs during seeding, panic and terminate execution
	if err := SeedUsers(db); err != nil {
		logger.Infof("Something else error when run seeding user: %+v", err)
	}

	if err := SeedPermissions(db); err != nil {
		logger.Infof("Something else error when run seeding permission: %+v", err)
	}

}
