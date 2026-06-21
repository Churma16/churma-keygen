package controllers

import (
	"net/http"

	"churma-keygen/backend/dtos"
	"churma-keygen/backend/services"

	"github.com/gin-gonic/gin"
)

type AuthController struct {
	authService services.AuthService
}

func NewAuthController(authService services.AuthService) *AuthController {
	return &AuthController{authService: authService}
}

func (ctrl *AuthController) Login(c *gin.Context) {
	var req dtos.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dtos.NewErrorResponse(http.StatusBadRequest, "Username and password are required"))
		return
	}

	resp, err := ctrl.authService.Login(req)
	if err != nil {
		c.JSON(http.StatusUnauthorized, dtos.NewErrorResponse(http.StatusUnauthorized, err.Error()))
		return
	}

	c.JSON(http.StatusOK, dtos.NewSuccessResponse(http.StatusOK, "Login berhasil", resp))
}

func (ctrl *AuthController) GetMe(c *gin.Context) {
	username, _ := c.Get("username")
	role, _ := c.Get("role")
	userID, _ := c.Get("userID")

	resp, err := ctrl.authService.GetMe(userID.(string), username.(string), role.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, dtos.NewErrorResponse(http.StatusInternalServerError, err.Error()))
		return
	}

	c.JSON(http.StatusOK, dtos.NewSuccessResponse(http.StatusOK, "Profil admin berhasil ditemukan", resp))
}

func (ctrl *AuthController) UpdateProfile(c *gin.Context) {
	var req dtos.UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dtos.NewErrorResponse(http.StatusBadRequest, "Username dan password saat ini wajib diisi"))
		return
	}

	userID, _ := c.Get("userID")
	err := ctrl.authService.UpdateProfile(userID.(string), req)
	if err != nil {
		c.JSON(http.StatusBadRequest, dtos.NewErrorResponse(http.StatusBadRequest, err.Error()))
		return
	}

	c.JSON(http.StatusOK, dtos.NewSuccessResponse(http.StatusOK, "Profil admin berhasil diperbarui", nil))
}
