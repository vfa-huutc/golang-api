package seeders

import (
	"github.com/vfa-khuongdv/golang-cms/pkg/logger"
	"gorm.io/gorm"
)

// Run executes all seed functions to populate the database with initial data
// It takes a GORM database connection as input and panics if any seeding operation fails
func Run(db *gorm.DB) {

	// SeendPermissions seeds the permissions table
	if err := SeedPermissions(db); err != nil {
		logger.Infof("Something else error when run seeding permission: %+v", err)
	}

	// SeedRoles seeds the roles table
	if err := SeedRoles(db); err != nil {
		logger.Infof("Something else error when run seeding user: %+v", err)
	}

	// SeedUsers seeds the users table
	if err := SeedUsers(db); err != nil {
		logger.Infof("Something else error when run seeding user: %+v", err)
	}

}
