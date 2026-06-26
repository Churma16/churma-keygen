package controllers

import (
	"net/http"

	"churma-keygen/backend/dtos"
	"churma-keygen/backend/usecase"

	"github.com/gin-gonic/gin"
)

type SettingController struct {
	settingUsecase usecase.SettingUsecase
}

func NewSettingController(settingUsecase usecase.SettingUsecase) *SettingController {
	return &SettingController{settingUsecase: settingUsecase}
}

func (ctrl *SettingController) GetSetting(c *gin.Context) {
	key := c.Param("key")
	if key == "" {
		c.JSON(http.StatusBadRequest, dtos.NewErrorResponse(http.StatusBadRequest, "Setting key is required"))
		return
	}

	setting, err := ctrl.settingUsecase.GetSetting(key)
	if err != nil {
		c.JSON(http.StatusNotFound, dtos.NewErrorResponse(http.StatusNotFound, "Setting not found"))
		return
	}

	c.JSON(http.StatusOK, dtos.NewSuccessResponse(http.StatusOK, "Setting retrieved successfully", setting))
}

func (ctrl *SettingController) UpdateSetting(c *gin.Context) {
	key := c.Param("key")
	if key == "" {
		c.JSON(http.StatusBadRequest, dtos.NewErrorResponse(http.StatusBadRequest, "Setting key is required"))
		return
	}

	var req dtos.UpdateSettingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dtos.NewErrorResponse(http.StatusBadRequest, err.Error()))
		return
	}

	setting, err := ctrl.settingUsecase.UpdateSetting(key, req.Value)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dtos.NewErrorResponse(http.StatusInternalServerError, err.Error()))
		return
	}

	c.JSON(http.StatusOK, dtos.NewSuccessResponse(http.StatusOK, "Setting updated successfully", setting))
}
