package controllers

import (
	"net/http"

	"churma-keygen/backend/dtos"
	"churma-keygen/backend/services"

	"github.com/gin-gonic/gin"
)

type ActivationController struct {
	activationService services.ActivationService
}

func NewActivationController(activationService services.ActivationService) *ActivationController {
	return &ActivationController{activationService: activationService}
}

func (ctrl *ActivationController) ActivateLicense(c *gin.Context) {
	var req dtos.ActivateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dtos.NewErrorResponse(http.StatusBadRequest, "License code and hardware ID are required"))
		return
	}

	ipAddress := c.ClientIP()

	resp, err := ctrl.activationService.Activate(req, ipAddress)
	if err != nil {
		// Distinguish between Conflict (HWID mismatch), Forbidden (suspended/expired), and NotFound (invalid key)
		switch err.Error() {
		case "invalid license code":
			c.JSON(http.StatusNotFound, dtos.NewErrorResponse(http.StatusNotFound, err.Error()))
		case "this license key is suspended or revoked", "this license key has expired":
			c.JSON(http.StatusForbidden, dtos.NewErrorResponse(http.StatusForbidden, err.Error()))
		case "hardware ID mismatch. This license is bound to another machine":
			c.JSON(http.StatusConflict, dtos.NewErrorResponse(http.StatusConflict, err.Error()))
		default:
			c.JSON(http.StatusInternalServerError, dtos.NewErrorResponse(http.StatusInternalServerError, err.Error()))
		}
		return
	}

	c.JSON(http.StatusOK, dtos.NewSuccessResponse(http.StatusOK, "Lisensi berhasil diaktivasi", resp))
}

func (ctrl *ActivationController) GetPublicKey(c *gin.Context) {
	pubPEM, err := ctrl.activationService.GetPublicKey()
	if err != nil {
		c.JSON(http.StatusInternalServerError, dtos.ErrorResponse{Error: err.Error()})
		return
	}

	c.String(http.StatusOK, pubPEM)
}
