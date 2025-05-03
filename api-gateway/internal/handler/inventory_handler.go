package handler

import (
	"net/http"
	"strconv"

	grpc "api-gateway/internal/client/grpc"

	"github.com/gin-gonic/gin"
)

type InventoryHandler struct{}

func (h InventoryHandler) GetSneakers(c *gin.Context) {
	role := c.Query("role")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	resp, err := grpc.GRPCGetSneakers(role, int32(page), int32(limit))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}

func (h InventoryHandler) CreateSneaker(c *gin.Context) {
	var body struct {
		Role  string `json:"role"`
		Brand string `json:"brand"`
		Model string `json:"model"`
		Price int32  `json:"price"`
		Color string `json:"color"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp, err := grpc.GRPCCreateSneaker(body.Role, body.Brand, body.Model, body.Price, body.Color)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, resp)
}

func (h InventoryHandler) EditSneaker(c *gin.Context) {
	var body struct {
		Role  string `json:"role"`
		Id    string `json:"id"`
		Brand string `json:"brand"`
		Model string `json:"model"`
		Price int32  `json:"price"`
		Color string `json:"color"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp, err := grpc.GRPCEditSneaker(body.Role, body.Id, body.Brand, body.Model, body.Price, body.Color)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}

func (h InventoryHandler) RemoveSneaker(c *gin.Context) {
	var body struct {
		Role string `json:"role"`
		Id   string `json:"id"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp, err := grpc.GRPCRemoveSneaker(body.Role, body.Id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}

func (h InventoryHandler) GetPublicSneakers(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	resp, err := grpc.GRPCGetPublicSneakers(int32(page), int32(limit))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}
