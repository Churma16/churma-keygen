package controllers

import (
	"net/http"

	"churma-keygen/backend/dtos"
	"churma-keygen/backend/usecase"

	"github.com/gin-gonic/gin"
)

type LicenseController struct {
	licenseUsecase usecase.LicenseUsecase
}

func NewLicenseController(licenseUsecase usecase.LicenseUsecase) *LicenseController {
	return &LicenseController{licenseUsecase: licenseUsecase}
}

func (ctrl *LicenseController) GetLicenses(c *gin.Context) {
	licenses, err := ctrl.licenseUsecase.GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, dtos.NewErrorResponse(http.StatusInternalServerError, "Failed to fetch licenses"))
		return
	}

	c.JSON(http.StatusOK, dtos.NewSuccessResponse(http.StatusOK, "Daftar lisensi berhasil ditemukan", licenses))
}

func (ctrl *LicenseController) GenerateLicense(c *gin.Context) {
	var req dtos.GenerateLicenseRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dtos.NewErrorResponse(http.StatusBadRequest, err.Error()))
		return
	}

	license, err := ctrl.licenseUsecase.Generate(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dtos.NewErrorResponse(http.StatusInternalServerError, err.Error()))
		return
	}

	c.JSON(http.StatusCreated, dtos.NewSuccessResponse(http.StatusCreated, "Lisensi berhasil dibuat", license))
}

func (ctrl *LicenseController) UpdateLicenseStatus(c *gin.Context) {
	id := c.Param("id")
	var req dtos.UpdateLicenseStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dtos.NewErrorResponse(http.StatusBadRequest, err.Error()))
		return
	}

	license, err := ctrl.licenseUsecase.UpdateStatus(id, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dtos.NewErrorResponse(http.StatusInternalServerError, err.Error()))
		return
	}

	c.JSON(http.StatusOK, dtos.NewSuccessResponse(http.StatusOK, "Status lisensi berhasil diperbarui", license))
}

func (ctrl *LicenseController) DeleteLicense(c *gin.Context) {
	id := c.Param("id")
	err := ctrl.licenseUsecase.Delete(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dtos.NewErrorResponse(http.StatusInternalServerError, err.Error()))
		return
	}

	c.JSON(http.StatusOK, dtos.NewSuccessResponse(http.StatusOK, "License deleted successfully", nil))
}

func (ctrl *LicenseController) GetActivationLogs(c *gin.Context) {
	logs, err := ctrl.licenseUsecase.GetActivationLogs()
	if err != nil {
		c.JSON(http.StatusInternalServerError, dtos.NewErrorResponse(http.StatusInternalServerError, "Failed to fetch activation logs"))
		return
	}

	c.JSON(http.StatusOK, dtos.NewSuccessResponse(http.StatusOK, "Log aktivasi berhasil ditemukan", logs))
}
