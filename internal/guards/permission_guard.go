package guards

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/vfa-khuongdv/golang-cms/internal/repositories"
	"github.com/vfa-khuongdv/golang-cms/internal/utils"
	"github.com/vfa-khuongdv/golang-cms/pkg/errors"
	"gorm.io/gorm"
)

// RoleGuard provides methods to check permissions for roles
type RoleGuard struct {
	roleRepo repositories.IRoleRepository
	userRepo repositories.IUserRepository
}

// NewRoleGuard creates a new RoleGuard instance with the provided database
func NewRoleGuard(db *gorm.DB) *RoleGuard {
	return &RoleGuard{
		roleRepo: repositories.NewRoleRepository(db),
		userRepo: repositories.NewUserRepository(db),
	}
}

// RequirePermissions returns a middleware that checks if the user has all required permissions
// It requires the auth middleware to be run first to set the user ID in the context
func RequirePermissions(guard *RoleGuard, requiredPerms ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get user ID
		userIdAny, exists := c.Get("UserID")
		if !exists {
			utils.RespondWithError(
				c,
				http.StatusForbidden,
				errors.New(errors.ErrForbidden, "User ID not found"),
			)
		}

		userId := userIdAny.(uint)
		userPermSet := map[string]bool{}

		// Get user permissions
		permissions, err := guard.userRepo.GetUserPermissions(userId)
		if err != nil {
			utils.RespondWithError(
				c,
				http.StatusForbidden,
				errors.New(errors.ErrForbidden, "Failed to retrieve user permissions"),
			)
		}

		// Aggregate permissions from all roles
		for _, p := range permissions {
			permKey := fmt.Sprintf("%s:%s", p.Resource, p.Action)
			userPermSet[permKey] = true
		}

		// Check ALL required permissions are granted
		for _, required := range requiredPerms {
			if !userPermSet[required] {
				utils.RespondWithError(
					c,
					http.StatusForbidden,
					errors.New(errors.ErrForbidden, fmt.Sprintf("Missing permission: %s", required)),
				)
			}
		}

		c.Next()
	}
}
