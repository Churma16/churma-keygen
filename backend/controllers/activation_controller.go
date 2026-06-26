package controllers

import (
	"net/http"

	"churma-keygen/backend/dtos"
	"churma-keygen/backend/usecase"

	"github.com/gin-gonic/gin"
)

type ActivationController struct {
	activationUsecase usecase.ActivationUsecase
}

func NewActivationController(activationUsecase usecase.ActivationUsecase) *ActivationController {
	return &ActivationController{activationUsecase: activationUsecase}
}

func (ctrl *ActivationController) ActivateLicense(c *gin.Context) {
	var req dtos.ActivateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(
			http.StatusBadRequest,
			dtos.NewErrorResponse(http.StatusBadRequest, "License code and hardware ID are required"),
		)
		return
	}

	ipAddress := c.ClientIP()

	resp, err := ctrl.activationUsecase.Activate(req, ipAddress)
	if err != nil {
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
	pubPEM, err := ctrl.activationUsecase.GetPublicKey()
	if err != nil {
		c.JSON(http.StatusInternalServerError, dtos.ErrorResponse{Error: err.Error()})
		return
	}

	c.String(http.StatusOK, pubPEM)
}

func (ctrl *ActivationController) GetContact(c *gin.Context) {
	resp := ctrl.activationUsecase.GetContact()
	c.JSON(http.StatusOK, dtos.NewSuccessResponse(http.StatusOK, "Contact info retrieved successfully", resp))
}
