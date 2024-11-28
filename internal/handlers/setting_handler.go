package handlers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/vfa-khuongdv/golang-cms/internal/models"
	"github.com/vfa-khuongdv/golang-cms/internal/services"
)

type ISettingHandler interface {
	GetSettings(c *gin.Context)
	UpdateSettings(c *gin.Context)
}

type SettingHandler struct {
	service *services.SettingService
}

func NewSettingHandler(service *services.SettingService) *SettingHandler {
	return &SettingHandler{service: service}
}

func (handler *SettingHandler) GetSettings(c *gin.Context) {
	// Get settings from service
	settings, err := handler.service.GetSetting()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, settings)
}

func (handler *SettingHandler) UpdateSettings(c *gin.Context) {

	type KeyValue struct {
		Key   string `json:"key" binding:"required"`
		Value string `json:"value" binding:"required"`
	}

	type Settings struct {
		Settings []KeyValue `json:"settings" binding:"required,dive"`
	}

	var input Settings

	// Bind JSON request body to input struct
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Iterate through settings array from request
	for _, v := range input.Settings {
		// Get existing setting by key
		value, err := handler.service.GetSettingByKey(v.Key)
		if err != nil {
			newSetting := models.Setting{
				SettingKey: v.Key,
				Value:      v.Value,
			}

			if err := handler.service.Create(&newSetting); err != nil {
				fmt.Printf("Create new setting error for key:%s value:%s\n", v.Key, v.Value)
				continue
			}

		}
		// Update setting value
		value.Value = v.Value

		// Save updated setting
		if err := handler.service.Update(value); err != nil {
			fmt.Printf("Update setting error for key:%s value:%s\n", v.Key, v.Value)
		}
	}

	c.JSON(http.StatusOK, gin.H{"message": "Update setting successfully"})
}
