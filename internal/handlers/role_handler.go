package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/vfa-khuongdv/golang-cms/internal/models"
	"github.com/vfa-khuongdv/golang-cms/internal/services"
)

type IRoleHandler interface {
	CreateRole(c *gin.Context)
	UpdateRole(c *gin.Context)
	GetRole(c *gin.Context)
	GetRoles(c *gin.Context)
	DeleteRole(c *gin.Context)
}

type RoleHandler struct {
	service *services.RoleService
}

func NewRoleHandler(service *services.RoleService) *RoleHandler {
	return &RoleHandler{
		service: service,
	}
}

func (handler *RoleHandler) CreateRole(c *gin.Context) {
	var input struct {
		Name        string `json:"name" binding:"required,min=3,max=255"`
		DisplayName string `json:"display_name" binding:"required,min=3,max=255"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	role := models.Role{
		Name:        input.Name,
		DisplayName: input.DisplayName,
	}

	if err := handler.service.Create(&role); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Create new role successfully"})
}

func (handler *RoleHandler) UpdateRole(c *gin.Context) {
	var input struct {
		DisplayName string `json:"display_name" binding:"required,min=3,max=255"`
	}

	// Bind JSON request body to input struct
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get role ID from URL parameter
	roleId := c.Param("id")
	// Convert role ID string to integer
	id, err := strconv.Atoi(roleId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get role from database by ID
	role, err := handler.service.GetByID(int64(id))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Update role display name
	role.DisplayName = input.DisplayName

	// Save updated role to database
	if err := handler.service.Update(role); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Update role successfully"})

}

func (handler *RoleHandler) GetRole(c *gin.Context) {
	// Get role ID from URL parameter
	roleId := c.Param("id")
	// Convert role ID string to integer
	id, err := strconv.Atoi(roleId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get role from database by ID
	role, err := handler.service.GetByID(int64(id))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, role)
}

func (handler *RoleHandler) DeleteRole(c *gin.Context) {
	// Get role ID from URL parameter
	roleId := c.Param("id")
	// Convert role ID string to integer
	id, err := strconv.Atoi(roleId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// Delete role from database
	if err := handler.service.Delete(int64(id)); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Delete role successfully"})
}
