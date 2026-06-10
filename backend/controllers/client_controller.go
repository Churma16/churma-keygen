package controllers

import (
	"net/http"

	"churma-keygen/backend/dtos"
	"churma-keygen/backend/services"

	"github.com/gin-gonic/gin"
)

type ClientController struct {
	clientService services.ClientService
}

func NewClientController(clientService services.ClientService) *ClientController {
	return &ClientController{clientService: clientService}
}

func (ctrl *ClientController) GetClients(c *gin.Context) {
	clients, err := ctrl.clientService.GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, dtos.ErrorResponse{Error: "Failed to fetch clients"})
		return
	}

	c.JSON(http.StatusOK, clients)
}

func (ctrl *ClientController) GetClientStats(c *gin.Context) {
	stats, err := ctrl.clientService.GetStats()
	if err != nil {
		c.JSON(http.StatusInternalServerError, dtos.ErrorResponse{Error: "Failed to fetch client statistics"})
		return
	}

	c.JSON(http.StatusOK, stats)
}

func (ctrl *ClientController) CreateClient(c *gin.Context) {
	var req dtos.CreateClientRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dtos.ErrorResponse{Error: err.Error()})
		return
	}

	client, err := ctrl.clientService.Create(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dtos.ErrorResponse{Error: "Failed to create client"})
		return
	}

	c.JSON(http.StatusCreated, client)
}

func (ctrl *ClientController) UpdateClient(c *gin.Context) {
	id := c.Param("id")
	var req dtos.UpdateClientRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dtos.ErrorResponse{Error: err.Error()})
		return
	}

	client, err := ctrl.clientService.Update(id, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dtos.ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, client)
}

func (ctrl *ClientController) DeleteClient(c *gin.Context) {
	id := c.Param("id")
	err := ctrl.clientService.Delete(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dtos.ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, dtos.MessageResponse{Message: "Client deleted successfully"})
}
