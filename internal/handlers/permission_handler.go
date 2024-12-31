package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/vfa-khuongdv/golang-cms/internal/services"
	"github.com/vfa-khuongdv/golang-cms/internal/utils"
)

type IPermissionHandler interface {
	GetAll(c *gin.Context)
}

type PermissionHandler struct {
	service *services.PermissionService
}

func NewPermissionHandler(service *services.PermissionService) *PermissionHandler {
	return &PermissionHandler{service: service}
}

func (handlder *PermissionHandler) GetAll(ctx *gin.Context) {
	permissions, err := handlder.service.GetAll()

	if err != nil {
		utils.RespondWithError(ctx, http.StatusBadRequest, err)
		return
	}

	utils.RespondWithOK(ctx, http.StatusOK, permissions)
}
