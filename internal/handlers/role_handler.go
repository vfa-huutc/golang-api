package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/vfa-khuongdv/golang-cms/internal/models"
	"github.com/vfa-khuongdv/golang-cms/internal/services"
	"github.com/vfa-khuongdv/golang-cms/internal/utils"
	"github.com/vfa-khuongdv/golang-cms/pkg/errors"
)

type IRoleHandler interface {
	CreateRole(c *gin.Context)
	UpdateRole(c *gin.Context)
	GetRole(c *gin.Context)
	GetRoles(c *gin.Context)
	DeleteRole(c *gin.Context)
	AssignPermissions(c *gin.Context)
	GetRolePermissions(c *gin.Context)
}

type RoleHandler struct {
	service services.IRoleService
}

func NewRoleHandler(service services.IRoleService) *RoleHandler {
	return &RoleHandler{
		service: service,
	}
}

func (handler *RoleHandler) CreateRole(ctx *gin.Context) {
	var input struct {
		Name        string `json:"name" binding:"required,min=3,max=255"`
		DisplayName string `json:"display_name" binding:"required,min=3,max=255"`
	}

	if err := ctx.ShouldBindJSON(&input); err != nil {
		utils.RespondWithError(
			ctx,
			http.StatusBadRequest,
			errors.New(errors.ErrInvalidData, err.Error()),
		)
		return
	}

	role := models.Role{
		Name:        input.Name,
		DisplayName: input.DisplayName,
	}

	if err := handler.service.Create(&role); err != nil {
		utils.RespondWithError(ctx, http.StatusBadRequest, err)
		return
	}

	utils.RespondWithOK(ctx, http.StatusCreated, gin.H{"message": "Create new role successfully"})
}

func (handler *RoleHandler) UpdateRole(ctx *gin.Context) {
	var input struct {
		DisplayName string `json:"display_name" binding:"required,min=3,max=255"`
	}

	// Bind JSON request body to input struct
	if err := ctx.ShouldBindJSON(&input); err != nil {
		utils.RespondWithError(
			ctx,
			http.StatusBadRequest,
			errors.New(errors.ErrInvalidData, err.Error()),
		)
		return
	}

	// Get role ID from URL parameter
	roleId := ctx.Param("id")
	// Convert role ID string to integer
	id, err := strconv.Atoi(roleId)
	if err != nil {
		utils.RespondWithError(
			ctx,
			http.StatusBadRequest,
			errors.New(errors.ErrInvalidParse, err.Error()),
		)
		return
	}

	// Get role from database by ID
	role, err := handler.service.GetByID(int64(id))
	if err != nil {
		utils.RespondWithError(ctx, http.StatusBadRequest, err)
		return
	}

	// Update role display name
	role.DisplayName = input.DisplayName

	// Save updated role to database
	if err := handler.service.Update(role); err != nil {
		utils.RespondWithError(ctx, http.StatusBadRequest, err)
		return
	}

	utils.RespondWithOK(ctx, http.StatusOK, gin.H{"message": "Update role successfully"})
}

func (handler *RoleHandler) GetRole(ctx *gin.Context) {
	// Get role ID from URL parameter
	roleId := ctx.Param("id")
	// Convert role ID string to integer
	id, err := strconv.Atoi(roleId)
	if err != nil {
		utils.RespondWithError(
			ctx,
			http.StatusBadRequest,
			errors.New(errors.ErrInvalidParse, err.Error()),
		)
		return
	}

	// Get role from database by ID
	role, err := handler.service.GetByID(int64(id))
	if err != nil {
		utils.RespondWithError(ctx, http.StatusBadRequest, err)
		return
	}

	utils.RespondWithOK(ctx, http.StatusOK, role)
}

func (handler *RoleHandler) DeleteRole(ctx *gin.Context) {
	// Get role ID from URL parameter
	roleId := ctx.Param("id")
	// Convert role ID string to integer
	id, err := strconv.Atoi(roleId)
	if err != nil {
		utils.RespondWithError(
			ctx,
			http.StatusBadRequest,
			errors.New(errors.ErrInvalidParse, err.Error()),
		)
		return
	}
	// Delete role from database
	if err := handler.service.Delete(int64(id)); err != nil {
		utils.RespondWithError(ctx, http.StatusBadRequest, err)
		return
	}
	utils.RespondWithOK(ctx, http.StatusOK, gin.H{"message": "Delete role successfully"})
}

// AssignPermissions handles assigning permissions to a role
// It accepts a list of permission IDs to assign to the specified role
func (handler *RoleHandler) AssignPermissions(ctx *gin.Context) {
	var input struct {
		PermissionIDs []uint `json:"permission_ids" binding:"required"`
	}

	// Bind JSON request body to input struct
	if err := ctx.ShouldBindJSON(&input); err != nil {
		utils.RespondWithError(
			ctx,
			http.StatusBadRequest,
			errors.New(errors.ErrInvalidData, err.Error()),
		)
		return
	}

	// Get role ID from URL parameter
	roleId := ctx.Param("id")
	// Convert role ID string to integer
	id, err := strconv.Atoi(roleId)
	if err != nil {
		utils.RespondWithError(
			ctx,
			http.StatusBadRequest,
			errors.New(errors.ErrInvalidParse, err.Error()),
		)
		return
	}

	// Verify the role exists
	_, err = handler.service.GetByID(int64(id))
	if err != nil {
		utils.RespondWithError(
			ctx,
			http.StatusNotFound,
			errors.New(errors.ErrResourceNotFound, "Role not found"),
		)
		return
	}

	// Assign permissions to the role
	if err := handler.service.AssignPermissions(uint(id), input.PermissionIDs); err != nil {
		utils.RespondWithError(ctx, http.StatusBadRequest, err)
		return
	}

	utils.RespondWithOK(ctx, http.StatusOK, gin.H{"message": "Permissions assigned to role successfully"})
}

// GetRolePermissions retrieves all permissions assigned to a role
func (handler *RoleHandler) GetRolePermissions(ctx *gin.Context) {
	// Get role ID from URL parameter
	roleId := ctx.Param("id")
	// Convert role ID string to integer
	id, err := strconv.Atoi(roleId)
	if err != nil {
		utils.RespondWithError(
			ctx,
			http.StatusBadRequest,
			errors.New(errors.ErrInvalidParse, err.Error()),
		)
		return
	}

	// Verify the role exists
	_, err = handler.service.GetByID(int64(id))
	if err != nil {
		utils.RespondWithError(
			ctx,
			http.StatusNotFound,
			errors.New(errors.ErrResourceNotFound, "Role not found"),
		)
		return
	}

	// Get permissions assigned to the role
	permissions, err := handler.service.GetRolePermissions(uint(id))
	if err != nil {
		utils.RespondWithError(ctx, http.StatusBadRequest, err)
		return
	}

	utils.RespondWithOK(ctx, http.StatusOK, permissions)
}
