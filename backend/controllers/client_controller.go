package controllers

import (
	"net/http"

	"churma-keygen/backend/dtos"
	"churma-keygen/backend/usecase"

	"github.com/gin-gonic/gin"
)

type ClientController struct {
	clientUsecase usecase.ClientUsecase
}

func NewClientController(clientUsecase usecase.ClientUsecase) *ClientController {
	return &ClientController{clientUsecase: clientUsecase}
}

func (ctrl *ClientController) GetClients(c *gin.Context) {
	clients, err := ctrl.clientUsecase.GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, dtos.NewErrorResponse(http.StatusInternalServerError, "Failed to fetch clients"))
		return
	}

	c.JSON(http.StatusOK, dtos.NewSuccessResponse(http.StatusOK, "Daftar klien berhasil ditemukan", clients))
}

func (ctrl *ClientController) GetClientStats(c *gin.Context) {
	stats, err := ctrl.clientUsecase.GetStats()
	if err != nil {
		c.JSON(http.StatusInternalServerError, dtos.NewErrorResponse(http.StatusInternalServerError, "Failed to fetch client statistics"))
		return
	}

	c.JSON(http.StatusOK, dtos.NewSuccessResponse(http.StatusOK, "Statistik klien berhasil ditemukan", stats))
}

func (ctrl *ClientController) CreateClient(c *gin.Context) {
	var req dtos.CreateClientRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dtos.NewErrorResponse(http.StatusBadRequest, err.Error()))
		return
	}

	client, err := ctrl.clientUsecase.Create(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dtos.NewErrorResponse(http.StatusInternalServerError, "Failed to create client"))
		return
	}

	c.JSON(http.StatusCreated, dtos.NewSuccessResponse(http.StatusCreated, "Klien berhasil dibuat", client))
}

func (ctrl *ClientController) UpdateClient(c *gin.Context) {
	id := c.Param("id")
	var req dtos.UpdateClientRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dtos.NewErrorResponse(http.StatusBadRequest, err.Error()))
		return
	}

	client, err := ctrl.clientUsecase.Update(id, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dtos.NewErrorResponse(http.StatusInternalServerError, err.Error()))
		return
	}

	c.JSON(http.StatusOK, dtos.NewSuccessResponse(http.StatusOK, "Klien berhasil diperbarui", client))
}

func (ctrl *ClientController) DeleteClient(c *gin.Context) {
	id := c.Param("id")
	err := ctrl.clientUsecase.Delete(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dtos.NewErrorResponse(http.StatusInternalServerError, err.Error()))
		return
	}

	c.JSON(http.StatusOK, dtos.NewSuccessResponse(http.StatusOK, "Client deleted successfully", nil))
}
