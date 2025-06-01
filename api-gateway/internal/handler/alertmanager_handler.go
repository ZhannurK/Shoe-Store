package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func AlertWebhookHandler(c *gin.Context) {
	var alertData interface{}
	if err := c.ShouldBindJSON(&alertData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "received", "data": alertData})
}
