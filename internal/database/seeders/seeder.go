package seeders

import (
	"fmt"

	"gorm.io/gorm"
)

// Run executes all seed functions to populate the database with initial data
// It takes a GORM database connection as input and panics if any seeding operation fails
func Run(db *gorm.DB) {

	// Execute the SeedUsers function to populate initial user data
	// If an error occurs during seeding, panic and terminate execution
	if err := SeedUsers(db); err != nil {
		fmt.Println("Something else error when run seeding user", err)
	}

	if err := SeedPermissions(db); err != nil {
		fmt.Println("Something else error when run seeding permission", err)
	}

}
